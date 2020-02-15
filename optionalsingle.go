package rxgo

import "context"

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
}

func newOptionalSingleFromOperator(ctx context.Context, iterable Iterable, nextFunc, errFunc, doneFunc Operator) OptionalSingle {
	next := operator(ctx, iterable, nextFunc, errFunc, doneFunc)

	return &optionalSingle{
		iterable: newChannelIterable(next),
	}
}

type optionalSingle struct {
	iterable Iterable
}

func (o *optionalSingle) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe()
}
