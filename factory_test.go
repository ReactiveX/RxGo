package rxgo

import (
	"context"
	"testing"
)

func Test_FromChannel(t *testing.T) {
	next := make(chan interface{})
	errs := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		next <- 1
		next <- 2
		next <- 3
		cancel()
	}()

	obs := FromChannel(ctx, next, errs)
	assertObservable(t, ctx, obs, hasItems(1, 2, 3))
}
