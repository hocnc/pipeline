package producer

import (
	"context"
	"testing"

	"github.com/hocnc/pipeline/end"
)

func TestGetPaths(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)
	defer close(errChannel)
	defer cancel()

	filePaths := GetFilePaths(ctx, "samples", errChannel)
	end.EndDebug(ctx, cancel, filePaths, errChannel)
}
