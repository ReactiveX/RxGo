package rxgo

type sliceIterable struct {
	items []Item
}

func newSliceIterable(items []Item) Iterable {
	return &sliceIterable{items: items}
}

func (i *sliceIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(opts...)
	next := option.buildChannel()

	go func() {
		for _, item := range i.items {
			next <- item
		}
		close(next)
	}()
	return next
}
