package rxgo

type sliceIterable struct {
	items []Item
	opts  []Option
}

func newSliceIterable(items []Item, opts ...Option) Iterable {
	return &sliceIterable{
		items: items,
		opts:  opts,
	}
}

func (i *sliceIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		for _, item := range i.items {
			select {
			case <-ctx.Done():
				return
			case next <- item:
			}
		}
		close(next)
	}()
	return next
}
