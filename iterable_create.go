package rxgo

import (
	"sync"
)

type createIterable struct {
	opts []Option
	next <-chan Item
}

func newCreateIterable(fs []Producer, opts ...Option) Iterable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	wg := sync.WaitGroup{}
	done := func() {
		wg.Done()
	}
	for _, f := range fs {
		wg.Add(1)
		go f(ctx, next, done)
	}
	go func() {
		wg.Wait()
		close(next)
	}()

	return &createIterable{
		opts: opts,
		next: next,
	}
}

func (i *createIterable) Observe() <-chan Item {
	return i.next
}
