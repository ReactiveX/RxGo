package rxgo

import "context"

type funcIterable struct {
	f func(ctx context.Context, next chan<- Item)
}

func newFuncIterable(f func(ctx context.Context, next chan<- Item)) Iterable {
	return &funcIterable{f: f}
}

func (i *funcIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(opts...)
	var next chan Item
	if toBeBuffered, cap := option.withBuffer(); toBeBuffered {
		next = make(chan Item, cap)
	} else {
		next = make(chan Item)
	}
	var ctx context.Context
	withContext, c := option.withContext()
	if withContext {
		ctx = c
	} else {
		ctx = context.Background()
	}

	go i.f(ctx, next)

	return next
}
