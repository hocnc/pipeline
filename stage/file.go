package stage

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
)

func GetFilePaths(ctx context.Context, folders <-chan string, errChannel chan<- error) <-chan string {
	outChannel := make(chan string)

	go func() {
		defer close(outChannel)

		for {

			select {
			case <-ctx.Done():
				return
			case folder, ok := <-folders:
				if ok {
					errChannel <- filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {

						if !info.Mode().IsRegular() {
							return nil
						}

						select {
						case <-ctx.Done():
							return fmt.Errorf("Walk Canceled")
						case outChannel <- path:
						}
						return nil
					})
				} else {
					return
				}
			}
		}
	}()

	return outChannel
}
