package end

import (
	"context"
	"log"
)

func End(ctx context.Context, cancelFunc context.CancelFunc, values <-chan string, errors chan error) {
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
				// fmt.Print(value)
			} else {
				return
			}
		}
	}
}
