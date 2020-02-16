package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
}

func newOptionalSingleFromOperator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) OptionalSingle {
	return &optionalSingle{
		iterable: newColdIterable(func() <-chan Item {
			next, ctx, option := buildOptionValues(opts...)
			if withPool, pool := option.withPool(); withPool {
				parallel(ctx, pool, next, iterable, nextFunc, errFunc, endFunc)
			} else {
				seq(ctx, next, iterable, nextFunc, errFunc, endFunc)
			}
			return next
		}),
	}
}

type optionalSingle struct {
	iterable Iterable
}

func (o *optionalSingle) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe()
}
