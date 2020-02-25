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
	for _, f := range i.f {
		f := f
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(ctx, next)
		}()
	}
	go func() {
		wg.Wait()
		close(next)
	}()

	return next
}
