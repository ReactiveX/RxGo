package rxgo

import "context"

type Iterable interface {
	Done() <-chan struct{}
	Error() <-chan error
	Next() <-chan interface{}
}

type Source struct {
	ctx  context.Context
	next <-chan interface{}
	errs <-chan error
}

func newIterable(ctx context.Context, next <-chan interface{}, errs <-chan error) Iterable {
	return &Source{ctx: ctx, next: next, errs: errs}
}

func (s *Source) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s *Source) Error() <-chan error {
	return s.errs
}

func (s *Source) Next() <-chan interface{} {
	return s.next
}
