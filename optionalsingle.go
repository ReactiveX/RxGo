package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
	// TODO Map
}

func newOptionalSingleFromOperator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) OptionalSingle {
	return &optionalSingle{
		iterable: operator(iterable, nextFunc, errFunc, endFunc, opts...),
	}
}

type optionalSingle struct {
	iterable Iterable
}

func (o *optionalSingle) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe()
}
