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
		iterable: newIterable(next),
		handler:  handler,
	}
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
	handler := func(ctx context.Context, src <-chan Item, dst chan<- Item) {
		for {
			select {
			case <-ctx.Done():
				close(dst)
				return
			case i, ok := <-src:
				if !ok {
					close(dst)
					return
				}
				if i.IsError() {
					dst <- i
					close(dst)
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
	return newObservable(ctx, o, handler)
}
