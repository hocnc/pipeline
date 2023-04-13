package producer

import (
	"context"
	"testing"

	"github.com/hocnc/pipeline/end"
)

func TestGetStdin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)

	defer cancel()
	defer close(errChannel)

	inputs := GetStdin(ctx, errChannel)
	end.End(ctx, cancel, inputs, errChannel)
}
