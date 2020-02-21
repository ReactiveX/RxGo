package rxgo

import (
	"sync"
)

type deferIterable struct {
	f    []Producer
	opts []Option
}

func newDeferIterable(f []Producer, opts ...Option) Iterable {
	return &deferIterable{
		f:    f,
		opts: opts,
	}
}

func (i *deferIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()
	ctx := option.buildContext()

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
