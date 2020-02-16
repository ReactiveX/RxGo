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

func newSingleFromOperator(iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd, opts ...Option) Single {
	return &single{
		iterable: operator(iterable, nextFunc, errFunc, endFunc, opts...),
	}
}

func (s *single) Observe(opts ...Option) <-chan Item {
	return s.iterable.Observe()
}

func (s *single) Filter(apply Predicate, opts ...Option) OptionalSingle {
	return newOptionalSingleFromOperator(s, func(item Item, dst chan<- Item, operator operatorOptions) {
		if apply(item.Value) {
			dst <- item
		}
		operator.stop()
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

func (s *single) Map(apply Func, opts ...Option) Single {
	return newSingleFromOperator(s, func(item Item, dst chan<- Item, operator operatorOptions) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			operator.stop()
		} else {
			dst <- FromValue(res)
			operator.stop()
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}
