package pipeline

import (
	"context"
	"testing"
)

func TestGetStdin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)

	defer cancel()
	defer close(errChannel)

	inputs := GetStdin(ctx, errChannel)
	End(ctx, cancel, inputs, errChannel)
}
