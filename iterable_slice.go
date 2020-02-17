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

func (i *sliceIterable) Observe() <-chan Item {
	option := parseOptions(i.opts...)
	next := option.buildChannel()

	go func() {
		for _, item := range i.items {
			next <- item
		}
		close(next)
	}()
	return next
}
