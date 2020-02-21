package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
	Run(opts ...Option) Disposed
	// TODO Map
}

// OptionalSingleImpl implements OptionalSingle.
type OptionalSingleImpl struct {
	iterable Iterable
}

// Observe observes an OptionalSingle by returning its channel.
func (o *OptionalSingleImpl) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}

// Run creates an observer without consuming the emitted items.
func (o *OptionalSingleImpl) Run(opts ...Option) Disposed {
	dispose := make(chan struct{})
	option := parseOptions(opts...)
	ctx := option.buildContext()

	go func() {
		defer close(dispose)
		observe := o.Observe()
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-observe:
				if !ok {
					return
				}
			}
		}
	}()

	return dispose
}
