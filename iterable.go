package rxgo

import "context"

type Iterable interface {
	Done() <-chan struct{}
	Next() <-chan Item
}

type Source struct {
	ctx  context.Context
	next <-chan Item
}

func newIterable(ctx context.Context, next <-chan Item) Iterable {
	return &Source{ctx: ctx, next: next}
}

func (s *Source) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s *Source) Next() <-chan Item {
	return s.next
}
