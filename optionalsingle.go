package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
}

func newOptionalSingleFromOperator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) OptionalSingle {
	next := operator(iterable, nextFunc, errFunc, endFunc, opts...)

	return &optionalSingle{
		iterable: newChannelIterable(next),
	}
}

type optionalSingle struct {
	iterable Iterable
}

func (o *optionalSingle) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe()
}
