package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
	// TODO Map
}

func newOptionalSingleFromOperator(iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd, opts ...Option) OptionalSingle {
	return &OptionalSingleImpl{
		iterable: operator(iterable, nextFunc, errFunc, endFunc, opts...),
	}
}

// OptionalSingleImpl implements OptionalSingle.
type OptionalSingleImpl struct {
	iterable Iterable
}

// Observe observes an OptionalSingle by returning its channel.
func (o *OptionalSingleImpl) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}
