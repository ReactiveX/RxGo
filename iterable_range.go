package rxgo

type rangeIterable struct {
	start, count int
	opts         []Option
}

func newRangeIterable(start, count int, opts ...Option) Iterable {
	return &rangeIterable{
		start: start,
		count: count,
		opts:  opts,
	}
}

func (i *rangeIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(append(i.opts, opts...)...)
	ctx := option.buildContext(emptyContext)
	next := option.buildChannel()

	go func() {
		for idx := i.start; idx <= i.start+i.count-1; idx++ {
			select {
			case <-ctx.Done():
				return
			case next <- Of(idx):
			}
		}
		close(next)
	}()
	return next
}
