package rxgo

type justIterable struct {
	items []interface{}
	opts  []Option
}

func newJustIterable(items ...interface{}) func(opts ...Option) Iterable {
	return func(opts ...Option) Iterable {
		return &justIterable{
			items: items,
			opts:  opts,
		}
	}
}

func (i *justIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()

	go SendItems(option.buildContext(), next, CloseChannel, i.items)
	return next
}
