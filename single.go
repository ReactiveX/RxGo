package rxgo

// Single is a observable with a single element.
type Single interface {
	Iterable
	Filter(apply Predicate, opts ...Option) OptionalSingle
	Map(apply Func, opts ...Option) Single
}

type single struct {
	iterable Iterable
}

func newSingleFromOperator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) Single {
	return &single{
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

func (s *single) Observe(opts ...Option) <-chan Item {
	return s.iterable.Observe()
}

func (s *single) Filter(apply Predicate, opts ...Option) OptionalSingle {
	return newOptionalSingleFromOperator(s, func(item Item, dst chan<- Item, stop func()) {
		if apply(item.Value) {
			dst <- item
		}
		stop()
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

func (s *single) Map(apply Func, opts ...Option) Single {
	return newSingleFromOperator(s, func(item Item, dst chan<- Item, stop func()) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			stop()
		} else {
			dst <- FromValue(res)
			stop()
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}
