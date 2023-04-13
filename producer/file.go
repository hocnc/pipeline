package producer

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
)

func GetFilePaths(ctx context.Context, folder string, errChannel chan<- error) <-chan string {
	outChannel := make(chan string)

	go func() {
		defer close(outChannel)

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

	}()

	return outChannel
}
