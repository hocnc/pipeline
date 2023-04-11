package pipeline

import (
	"context"
	"log"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

var Limit = runtime.NumCPU()

func splitWithBandWidth(result string, limit int) []string {
	var newResult []string
	var tmp strings.Builder
	i := 0

	values := strings.Split(result, "\n")
	n := len(values)

	bandwidth := n / limit
	if bandwidth == 0 {
		bandwidth = 1
	}
	for _, v := range values {
		if i == bandwidth {
			// fmt.Println("-----")
			// fmt.Print(tmp.String())
			newResult = append(newResult, tmp.String())
			i = 0
			tmp.Reset()
		}
		tmp.WriteString(v + "\n")
		i++
	}
	if tmp.Len() != 0 {
		newResult = append(newResult, tmp.String())
	}
	return newResult
}

func LoadBalancer(
	ctx context.Context,
	inChannel <-chan string,
	errChannel chan<- error) <-chan string {
	outChannel := make(chan string)

	go func() {
		defer close(outChannel)

		for {

			select {
			case <-ctx.Done():
				return
			case val, ok := <-inChannel:

				if ok {

					for _, v := range splitWithBandWidth(val, Limit) {
						outChannel <- v
					}

				} else {
					return
				}
			}
		}

	}()

	return outChannel
}

func Producer(ctx context.Context, values []string, errChannel chan<- error) <-chan string {
	outChannel := make(chan string)

	go func() {
		defer close(outChannel)

		for _, value := range values {
			if value == "\n" {
				continue
			}

			// log.Println(value)
			select {
			case <-ctx.Done():
				return
			case outChannel <- value:
			}
		}
	}()

	return outChannel
}

func ANew[T any](
	ctx context.Context,
	anew func(context.Context, T) (T, bool),
	insertDB func(context.Context, T) error,
	notify func(context.Context, T) error,
	isFinal bool,
	inChannel <-chan T,
	errChannel chan<- error) <-chan T {

	outChannel := make(chan T)

	go func() {
		defer close(outChannel)

		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-inChannel:

				if ok {

					//ANew and Send to outChannel
					newResult, isEmpty := anew(ctx, result)

					//If newResult is empty, next
					if isEmpty {
						continue
					}

					//Send to the next step
					if !isFinal {
						outChannel <- newResult
					}

					// //Insert into Database
					err := insertDB(ctx, newResult)
					if err != nil {
						errChannel <- err
						break
					}
					//Notify
					err = notify(ctx, newResult)
					if err != nil {
						errChannel <- err
						break
					}
				} else {
					return
				}

			}

		}
	}()

	return outChannel
}

func EndDebug[T any](ctx context.Context, cancelFunc context.CancelFunc, values <-chan T, errors chan error) {
	for {
		select {
		case <-ctx.Done():
			log.Print(ctx.Err().Error())
			return
		case err := <-errors:
			if err != nil {
				log.Println("error: ", err.Error())
				close(errors)
				cancelFunc()
			}
		case value, ok := <-values:
			if ok {
				log.Println(value)
			} else {
				log.Println("Done")
				return
			}
		}
	}
}

func End[T any](ctx context.Context, cancelFunc context.CancelFunc, values <-chan T, errors chan error) {
	for {
		select {
		case <-ctx.Done():
			log.Print(ctx.Err().Error())
			return
		case err := <-errors:
			if err != nil {
				log.Println("error: ", err.Error())
				close(errors)
				cancelFunc()
			}
		case _, ok := <-values:
			if ok {
				log.Println("Scanning")
			} else {
				log.Println("Done")
				return
			}
		}
	}
}

func Run[In any, Out any](
	ctx context.Context,
	inChannel <-chan In,
	errChannel chan<- error,
	fns ...func(context.Context, In) (Out, error)) <-chan Out {
	sem := semaphore.NewWeighted(int64(Limit))

	outChannel := make(chan Out)
	go func() {
		defer close(outChannel)

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-inChannel:

				if ok {
					for _, fn := range fns {

						if err := sem.Acquire(ctx, 1); err != nil {
							log.Printf("Failed to acquire semmaphore: %v", err)
							return
						}

						go func(ctx context.Context, val In, fn func(context.Context, In) (Out, error)) {
							defer sem.Release(1)

							result, err := fn(ctx, val)
							if err != nil {
								errChannel <- err
							} else {
								outChannel <- result
							}

						}(ctx, val, fn)

					}
				} else {
					//Make sure all done
					if err := sem.Acquire(ctx, int64(Limit)); err != nil {
						log.Printf("Failed to acquire semaphore: %v", err)
					}
					return
				}

			}

		}
	}()

	return outChannel
}

func Merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	multiplexedStream := make(chan T, 100)

	multiplex := func(c <-chan T) {
		defer wg.Done()

		for i := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	//Waitting untill all channels are closed
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
