package rxgo

type rangeIterable struct {
	start, count int
}

func newRangeIterable(start, count int) Iterable {
	return &rangeIterable{
		start: start,
		count: count,
	}
}

func (i *rangeIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(opts...)
	next := option.buildChannel()

	go func() {
		for idx := i.start; idx <= i.start+i.count; idx++ {
			next <- FromValue(idx)
		}
		close(next)
	}()
	return next
}
