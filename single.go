package rxgo

import "context"

// Single is a observable with a single element.
type Single interface {
	Iterable
	Filter(ctx context.Context, apply Predicate) OptionalSingle
	Map(ctx context.Context, apply Function) Single
}

type single struct {
	iterable Iterable
}

func newSingleFromOperator(ctx context.Context, iterable Iterable, nextFunc, errFunc, endFunc Operator) Single {
	next := operator(ctx, iterable, nextFunc, errFunc, endFunc)

	return &single{
		iterable: newChannelIterable(next),
	}
}

func (s *single) Observe(opts ...Option) <-chan Item {
	return s.iterable.Observe()
}

func (s *single) Filter(ctx context.Context, apply Predicate) OptionalSingle {
	return newOptionalSingleFromOperator(ctx, s, func(item Item, dst chan<- Item, stop func()) {
		if apply(item.Value) {
			dst <- item
		}
		stop()
	}, defaultEndFuncOperator, defaultEndFuncOperator)
}

func (s *single) Map(ctx context.Context, apply Function) Single {
	return newSingleFromOperator(ctx, s, func(item Item, dst chan<- Item, stop func()) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			stop()
		} else {
			dst <- FromValue(res)
			stop()
		}
	}, defaultEndFuncOperator, defaultEndFuncOperator)
}
