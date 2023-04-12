package pipeline

import (
	"context"
	"testing"
)

func TestGetPaths(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)
	defer close(errChannel)
	defer cancel()

	filePaths := GetFilePaths(ctx, "samples", errChannel)
	EndDebug(ctx, cancel, filePaths, errChannel)
}
