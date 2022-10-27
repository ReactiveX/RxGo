package rxgo

// Counts the number of emissions on the source and emits that number when the source completes.
func Count[T any](predicate ...PredicateFunc[T]) OperatorFunc[T, uint] {
	cb := skipPredicate[T]
	if len(predicate) > 0 {
		cb = predicate[0]
	}
	return func(source Observable[T]) Observable[uint] {
		var (
			count uint
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[uint], v T) {
				if cb(v, index) {
					count++
				}
				index++
			},
			func(obs Observer[uint], err error) {
				obs.Error(err)
			},
			func(obs Observer[uint]) {
				obs.Next(count)
				obs.Complete()
			},
		)
	}
}

// The Max operator operates on an Observable that emits numbers
// (or items that can be compared with a provided function),
// and when source Observable completes it emits a single item: the item with the largest value.
func Max[T any](comparer ...ComparerFunc[T, T]) OperatorFunc[T, T] {
	cb := maximum[T]
	if len(comparer) > 0 {
		cb = comparer[0]
	}
	return func(source Observable[T]) Observable[T] {
		var (
			lastValue T
			first     = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if first {
					lastValue = v
					first = false
					return
				}

				if cb(lastValue, v) < 0 {
					lastValue = v
				}
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Next(lastValue)
				obs.Complete()
			},
		)
	}
}

// The Min operator operates on an Observable that emits numbers
// (or items that can be compared with a provided function),
// and when source Observable completes it emits a single item: the item with the smallest value.
func Min[T any](comparer ...ComparerFunc[T, T]) OperatorFunc[T, T] {
	cb := minimum[T]
	if len(comparer) > 0 {
		cb = comparer[0]
	}
	return func(source Observable[T]) Observable[T] {
		var (
			lastValue T
			first     = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if first {
					lastValue = v
					first = false
					return
				}

				if cb(lastValue, v) >= 0 {
					lastValue = v
				}
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Next(lastValue)
				obs.Complete()
			},
		)
	}
}

// Applies an accumulator function over the source Observable, and returns
// the accumulated result when the source completes, given an optional seed value.
func Reduce[V any, A any](accumulator AccumulatorFunc[A, V], seed A) OperatorFunc[V, A] {
	if accumulator == nil {
		panic(`rxgo: "Reduce" expected accumulator func`)
	}
	return func(source Observable[V]) Observable[A] {
		var (
			index  uint
			result = seed
			err    error
		)
		return createOperatorFunc(
			source,
			func(obs Observer[A], v V) {
				result, err = accumulator(result, v, index)
				if err != nil {
					obs.Error(err)
					return
				}
				index++
			},
			func(obs Observer[A], err error) {
				obs.Error(err)
			},
			func(obs Observer[A]) {
				obs.Next(result)
				obs.Complete()
			},
		)
	}
}
