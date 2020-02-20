package rxgo

import "context"

// Single is a observable with a single element.
type Single interface {
	Iterable
	Filter(apply Predicate, opts ...Option) OptionalSingle
	Map(apply Func, opts ...Option) Single
}

// SingleImpl implements Single.
type SingleImpl struct {
	iterable Iterable
}

func newSingleFromOperator(iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd, opts ...Option) Single {
	return &SingleImpl{
		iterable: operator(iterable, nextFunc, errFunc, endFunc, opts...),
	}
}

// Observe observes a Single by returning its channel.
func (s *SingleImpl) Observe(opts ...Option) <-chan Item {
	return s.iterable.Observe(opts...)
}

// Filter emits only those items from a Single that pass a predicate test.
func (s *SingleImpl) Filter(apply Predicate, opts ...Option) OptionalSingle {
	return newOptionalSingleFromOperator(s, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if apply(item.V) {
			dst <- item
		}
		operator.stop()
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// Map transforms the items emitted by a Single by applying a function to each item.
func (s *SingleImpl) Map(apply Func, opts ...Option) Single {
	return newSingleFromOperator(s, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		res, err := apply(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
		} else {
			dst <- Of(res)
			operator.stop()
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}
