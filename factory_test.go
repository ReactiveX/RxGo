package rxgo

import (
	"context"
	"testing"
)

func Test_FromChannel(t *testing.T) {
	next := make(chan Item)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		cancel()
	}()

	obs := FromChannel(ctx, next)
	assertObservable(t, ctx, obs, hasItems(1, 2, 3))
}
