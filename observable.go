package rxgo

import (
	"context"
)

// Observable is the basic observable interface.
type Observable interface {
	Iterable
	Filter(ctx context.Context, apply Predicate) Observable
	ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc)
	Map(ctx context.Context, apply Func) Observable
	SkipWhile(ctx context.Context, apply Predicate) Observable
}

type observable struct {
	iterable Iterable
	handler  Handler
}

func newObservable(ctx context.Context, source Observable, handler Handler) Observable {
	next := make(chan Item)

	go handler(ctx, source.Next(), next)

	return &observable{
		iterable: newChannelIterable(next),
		handler:  handler,
	}
}

func newOperator(ctx context.Context, source Observable, nextFunc Operator, errFunc Operator) Observable {
	next := make(chan Item)

	stopped := false
	stop := func() {
		stopped = true
	}
	go func() {
		for !stopped {
			select {
			case <-ctx.Done():
				close(next)
				return
			case i, ok := <-source.Next():
				if !ok {
					close(next)
					return
				}
				if i.IsError() {
					errFunc(i, next, stop)
				} else {
					nextFunc(i, next, stop)
				}
			}
		}
		close(next)
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

func (o *observable) Next() <-chan Item {
	return o.iterable.Next()
}

func (o *observable) Filter(ctx context.Context, apply Predicate) Observable {
	return newOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if apply(item.Value) {
			dst <- item
		}
	}, func(item Item, dst chan<- Item, stop func()) {
		dst <- item
		stop()
	})
}

func (o *observable) ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc) {
	handler := func(ctx context.Context, src <-chan Item, dst chan<- Item) {
		for {
			select {
			case <-ctx.Done():
				doneFunc()
				return
			case i, ok := <-src:
				if !ok {
					doneFunc()
					return
				}
				if i.IsError() {
					errFunc(i.Err)
					return
				}
				nextFunc(i.Value)
			}
		}
	}
	newObservable(ctx, o, handler)
}

func (o *observable) Map(ctx context.Context, apply Func) Observable {
	return newOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
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

func (o *observable) SkipWhile(ctx context.Context, apply Predicate) Observable {
	skip := true

	return newOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if !skip {
			dst <- item
		} else {
			if !apply(item.Value) {
				skip = false
				dst <- item
			}
		}
	}, func(item Item, dst chan<- Item, stop func()) {
		dst <- item
		stop()
	})
}
