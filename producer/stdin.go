package pipeline

import (
	"bufio"
	"context"
	"os"
)

func GetStdin(ctx context.Context, errChannel chan<- error) <-chan string {
	outChannel := make(chan string)

	go func() {
		defer close(outChannel)

		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			input := sc.Text()

			select {
			case <-ctx.Done():
				return
			case outChannel <- input:
			}
		}

	}()

	return outChannel
}
