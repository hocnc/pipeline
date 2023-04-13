package stage

import (
	"context"
	"testing"

	"github.com/hocnc/pipeline/end"
	"github.com/hocnc/pipeline/producer"
)

func TestGetFilePaths(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)

	defer cancel()
	defer close(errChannel)

	filePaths := producer.GetFilePaths(ctx, "samples", errChannel)
	end.EndDebug(ctx, cancel, filePaths, errChannel)
}
