package rxgo

import (
	"context"
	"sync"
)

type funcsIterable struct {
	f []Scatter
}

func newFuncsIterable(f ...Scatter) Iterable {
	return &funcsIterable{f: f}
}

func (i *funcsIterable) Observe(opts ...Option) <-chan Item {
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

	wg := sync.WaitGroup{}
	done := func() {
		wg.Done()
	}
	for _, f := range i.f {
		wg.Add(1)
		go f(ctx, next, done)
	}
	go func() {
		wg.Wait()
		close(next)
	}()

	return next
}
