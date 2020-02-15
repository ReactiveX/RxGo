package rxgo

import "context"

func FromChannel(ctx context.Context, next <-chan Item) Observable {
	return &observable{
		iterable: newIterable(ctx, next),
	}
}
