package rxgo

type justIterable struct {
	items interface{}
	opts  []Option
}

func newJustIterable(items interface{}, opts ...Option) Iterable {
	return &justIterable{
		items: items,
		opts:  opts,
	}
}

func (i *justIterable) Observe() <-chan Item {
	option := parseOptions(i.opts...)
	next := option.buildChannel()

	go SendItems(next, CloseChannel, i.items)
	return next
}
