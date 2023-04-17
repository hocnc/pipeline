package main

import (
	"context"

	"github.com/hocnc/pipeline/end"
	"github.com/hocnc/pipeline/producer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	errChannel := make(chan error)

	defer cancel()
	defer close(errChannel)

	inputs := producer.GetStdin(ctx, errChannel)
	end.End(ctx, cancel, inputs, errChannel)
}
