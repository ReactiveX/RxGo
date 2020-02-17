package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
	// TODO Map
}

func newOptionalSingleFromOperator(iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd, opts ...Option) OptionalSingle {
	return &optionalSingle{
		iterable: operator(iterable, nextFunc, errFunc, endFunc, opts...),
	}
}

type optionalSingle struct {
	iterable Iterable
}

func (o *optionalSingle) Observe() <-chan Item {
	return o.iterable.Observe()
}
