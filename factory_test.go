package rxgo

import (
	"context"
	"testing"
)

func Test_FromChannel(t *testing.T) {
	next := make(chan Item)
	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		close(next)
	}()

	obs := FromChannel(next)
	assertObservable(t, context.Background(), obs, hasItems(1, 2, 3), hasNotRaisedError())
}
