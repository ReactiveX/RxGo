package rxgo

type deferIterable struct {
	fs   []Producer
	opts []Option
}

func newDeferIterable(f []Producer, opts ...Option) Iterable {
	return &deferIterable{
		fs:   f,
		opts: opts,
	}
}

func (i *deferIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		for _, f := range i.fs {
			f(ctx, next)
		}
	}()

	return next
}
