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
	for _, f := range fs {
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

	return &createIterable{
		opts: opts,
		next: next,
	}
}

func (i *createIterable) Observe(_ ...Option) <-chan Item {
	return i.next
}
