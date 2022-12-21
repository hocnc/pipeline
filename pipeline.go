package pipeline

import "sync"

func Take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func PrimeFinder(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})
	go func() {
		defer close(primeStream)
		for integer := range intStream {
			integer -= 1
			prime := true
			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()
	return primeStream
}

func RepeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}

	}()
	return valueStream
}

// func MyFanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
// 	multiplexedStream := make(chan interface{})
// 	go func() {
// 		defer close(multiplexedStream)
// 		for c := range channels {
// 			//Waitting c channel until the c channel is closed -- Waitting forever
// 			for v := range c {
// 				select {
// 				case <-done:
// 					return
// 				case multiplexedStream <- v:
// 				}
// 			}
// 		}
// 	}()
// 	return multiplexedStream
// }

func FanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{}, 100)

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()

		for i := range c {
			select {
			case <-done:
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
