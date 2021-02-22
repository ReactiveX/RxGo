package rxgo

import (
	"context"
	"errors"
)

type testStruct struct {
	ID int `json:"id"`
}

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")
)

func channelValue(ctx context.Context, items ...interface{}) chan Item {
	next := make(chan Item)
	go func() {
		for _, item := range items {
			switch item := item.(type) {
			default:
				Of(item).SendContext(ctx, next)
			case error:
				Error(item).SendContext(ctx, next)
			}
		}
		close(next)
	}()
	return next
}

func testObservable(ctx context.Context, items ...interface{}) Observable {
	return FromChannel(channelValue(ctx, items...))
}
