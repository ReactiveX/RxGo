package rxgo

import "reflect"

// Returns an Observable that emits all items emitted by the source Observable
// that are distinct by comparison from previous items.
func Distinct[T any, K comparable](keySelector func(value T) K) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			keySet = make(map[K]bool)
			exists bool
			key    K
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				key = keySelector(v)
				_, exists = keySet[key]
				if !exists {
					keySet[key] = true
					obs.Next(v)
				}
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Complete()
			},
		)
	}
}

// Returns a result Observable that emits all values pushed by the source observable
// if they are distinct in comparison to the last value the result observable emitted.
func DistinctUntilChanged[T any](comparator ...ComparatorFunc[T, T]) OperatorFunc[T, T] {
	cb := func(prev T, current T) bool {
		return reflect.DeepEqual(prev, current)
	}
	if len(comparator) > 0 {
		cb = comparator[0]
	}
	return func(source Observable[T]) Observable[T] {
		var (
			lastValue T
			first     = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if first || !cb(lastValue, v) {
					obs.Next(v)
					first = false
					lastValue = v
				}
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Complete()
			},
		)
	}
}

// Filter emits only those items from an Observable that pass a predicate test.
func Filter[T any](predicate PredicateFunc[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
		)
		cb := skipPredicate[T]
		if predicate != nil {
			cb = predicate
		}
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if cb(v, index) {
					obs.Next(v)
				}
				index++
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Complete()
			},
		)
	}
}
