package rxgo

import "context"

func FromChannel(ctx context.Context, next <-chan interface{}, errs <-chan error) Observable {
	return &observable{
		iterable: newIterable(ctx, next, errs),
	}
}
