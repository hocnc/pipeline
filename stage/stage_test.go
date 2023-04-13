package pipeline

import (
	"context"
	"testing"
)

func TestGetFilePaths(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)

	defer cancel()
	defer close(errChannel)

	filePaths := GetFilePaths(ctx, "samples", errChannel)
	EndDebug(ctx, cancel, filePaths, errChannel)
}
