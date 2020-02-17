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

func (i *rangeIterable) Observe() <-chan Item {
	option := parseOptions(i.opts...)
	next := option.buildChannel()

	go func() {
		for idx := i.start; idx <= i.start+i.count; idx++ {
			next <- Of(idx)
		}
		close(next)
	}()
	return next
}
