package rxgo

import "context"

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
}

func newOptionalSingleFromOperator(ctx context.Context, source Single, nextFunc Operator, errFunc Operator) OptionalSingle {
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
