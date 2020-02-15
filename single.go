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

func newSingleFromOperator(ctx context.Context, source Single, nextFunc Operator, errFunc Operator) Single {
	next := make(chan Item)

	stop := func() {}
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(next)
				return
			case i, ok := <-source.Observe():
				if !ok {
					close(next)
					return
				}
				if i.IsError() {
					errFunc(i, next, stop)
					close(next)
					return
				}
				nextFunc(i, next, stop)
				close(next)
				return
			}
		}
	}()

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
	}, func(item Item, dst chan<- Item, stop func()) {
		dst <- item
		stop()
	})
}

func (s *single) Map(ctx context.Context, apply Function) Single {
	return newSingleFromOperator(ctx, s, func(item Item, dst chan<- Item, stop func()) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			stop()
		}
		dst <- FromValue(res)
	}, func(item Item, dst chan<- Item, stop func()) {
		dst <- item
		stop()
	})
}
