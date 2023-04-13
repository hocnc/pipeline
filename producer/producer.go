package producer

import "context"

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
