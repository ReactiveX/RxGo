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
	next := make(chan interface{})
	errs := make(chan error)

	go handler(ctx, source.Next(), source.Error(), next, errs)

	return &observable{
		iterable: newIterable(ctx, next, errs),
		handler:  handler,
	}
}

func (o *observable) Done() <-chan struct{} {
	return o.iterable.Done()
}

func (o *observable) Error() <-chan error {
	return o.iterable.Error()
}

func (o *observable) Next() <-chan interface{} {
	return o.iterable.Next()
}

func (o *observable) ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc) {
	handler := func(ctx context.Context, nextSrc <-chan interface{}, errsSrc <-chan error, nextDst chan<- interface{}, errsDst chan<- error) {
		for {
			select {
			case <-ctx.Done():
				doneFunc()
				return
			case i := <-nextSrc:
				nextFunc(i)
			case err := <-errsSrc:
				errFunc(err)
				return
			}
		}
	}
	newObservable(ctx, o, handler)
}

func (o *observable) Map(ctx context.Context, apply Func) Observable {
	cancel, cancelFunc := context.WithCancel(ctx)
	handler := func(ctx context.Context, nextSrc <-chan interface{}, errsSrc <-chan error, nextDst chan<- interface{}, errsDst chan<- error) {
		for {
			select {
			case <-ctx.Done():
				cancelFunc()
			case i := <-nextSrc:
				res, err := apply(i)
				if err != nil {
					errsDst <- err
					return
				}
				nextDst <- res
			case err := <-errsSrc:
				errsDst <- err
			}
		}
	}
	return newObservable(cancel, o, handler)
}
