package end

import (
	"context"
	"log"
)

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
