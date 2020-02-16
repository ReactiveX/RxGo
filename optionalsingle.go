package rxgo

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
}

func newOptionalSingleFromOperator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) OptionalSingle {
	return &optionalSingle{
		iterable: newColdIterable(func() <-chan Item {
			next, ctx, pool := buildOptionValues(opts...)
			if pool == 0 {
				seq(ctx, next, iterable, nextFunc, errFunc, endFunc)
			} else {
				parallel(ctx, pool, next, iterable, nextFunc, errFunc, endFunc)
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
