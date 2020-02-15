package rxgo

import (
	"context"
)

type Observable interface {
	Iterable
	ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc)
	Map(ctx context.Context, apply Func) Observable
}

type observable struct {
	iterable Iterable
	handler  Handler
}

func newObservable(ctx context.Context, source Observable, handler Handler) Observable {
	next := make(chan Item)

	go handler(ctx, source.Next(), next)

	return &observable{
		iterable: newIterable(ctx, next),
		handler:  handler,
	}
}

func (o *observable) Done() <-chan struct{} {
	return o.iterable.Done()
}

func (o *observable) Next() <-chan Item {
	return o.iterable.Next()
}

func (o *observable) ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc) {
	handler := func(ctx context.Context, src <-chan Item, dst chan<- Item) {
		for {
			select {
			case <-ctx.Done():
				doneFunc()
				return
			case i := <-src:
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
	cancel, cancelFunc := context.WithCancel(ctx)
	handler := func(ctx context.Context, src <-chan Item, dst chan<- Item) {
		for {
			select {
			case <-ctx.Done():
				cancelFunc()
			case i := <-src:
				if i.IsError() {
					dst <- i
					return
				}

				res, err := apply(i.Value)
				if err != nil {
					dst <- FromError(err)
					return
				}
				dst <- FromValue(res)
			}
		}
	}
	return newObservable(cancel, o, handler)
}
