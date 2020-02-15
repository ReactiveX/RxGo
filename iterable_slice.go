package rxgo

type sliceIterable struct {
	items []Item
}

func newSliceIterable(items []Item) Iterable {
	return &sliceIterable{items: items}
}

func (i *sliceIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(opts...)
	var next chan Item
	if toBeBuffered, cap := option.withBuffer(); toBeBuffered {
		next = make(chan Item, cap)
	} else {
		next = make(chan Item)
	}

	go func() {
		for _, item := range i.items {
			next <- item
		}
		close(next)
	}()
	return next
}
