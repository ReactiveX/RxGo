package rxgo

// Emits a given value if the source Observable completes without emitting any
// next value, otherwise mirrors the source Observable.
func DefaultIfEmpty[T any](defaultValue T) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			hasValue bool
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				hasValue = true
				obs.Next(v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				if !hasValue {
					obs.Next(defaultValue)
				}
				obs.Complete()
			},
		)
	}
}

// Returns an Observable that emits whether or not every item of the
// source satisfies the condition specified.
func Every[T any](predicate PredicateFunc[T]) OperatorFunc[T, bool] {
	return func(source Observable[T]) Observable[bool] {
		var (
			allOk = true
			index uint
		)
		cb := skipPredicate[T]
		if predicate != nil {
			cb = predicate
		}
		return createOperatorFunc(
			source,
			func(obs Observer[bool], v T) {
				allOk = allOk && cb(v, index)
			},
			func(obs Observer[bool], err error) {
				obs.Error(err)
			},
			func(obs Observer[bool]) {
				obs.Next(allOk)
				obs.Complete()
			},
		)
	}
}

// Emits only the first value emitted by the source Observable that meets some condition.
func Find[T any](predicate PredicateFunc[T]) OperatorFunc[T, Optional[T]] {
	return func(source Observable[T]) Observable[Optional[T]] {
		var (
			found bool
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[Optional[T]], v T) {
				if predicate(v, index) {
					found = true
					obs.Next(Some(v))
					obs.Complete()
					return
				}
				index++
			},
			func(obs Observer[Optional[T]], err error) {
				obs.Error(err)
			},
			func(obs Observer[Optional[T]]) {
				if !found {
					obs.Next(None[T]())
				}
				obs.Complete()
			},
		)
	}
}

// Emits only the index of the first value emitted by the source Observable that meets some condition.
func FindIndex[T any](predicate PredicateFunc[T]) OperatorFunc[T, int] {
	var (
		index uint
		found bool
	)
	return func(source Observable[T]) Observable[int] {
		return createOperatorFunc(
			source,
			func(obs Observer[int], v T) {
				if predicate(v, index) {
					found = true
					obs.Next(int(index))
					obs.Complete()
				}
				index++
			},
			func(obs Observer[int], err error) {
				obs.Error(err)
			},
			func(obs Observer[int]) {
				if !found {
					obs.Next(-1)
				}
				obs.Complete()
			},
		)
	}
}

// Emits false if the input Observable emits any values,
// or emits true if the input Observable completes without emitting any values.
func IsEmpty[T any]() OperatorFunc[T, bool] {
	return func(source Observable[T]) Observable[bool] {
		var (
			empty = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[bool], v T) {
				empty = false
			},
			func(obs Observer[bool], err error) {
				obs.Error(err)
			},
			func(obs Observer[bool]) {
				obs.Next(empty)
				obs.Complete()
			},
		)
	}
}

// If the source observable completes without emitting a value, it will emit an error.
// The error will be created at that time by the optional errorFactory argument, otherwise,
// the error will be `ErrEmpty`.
func ThrowIfEmpty[T any](errorFactory ...ErrorFunc) OperatorFunc[T, T] {
	factory := func() error {
		return ErrEmpty
	}
	if len(errorFactory) > 0 {
		factory = errorFactory[0]
	}
	return func(source Observable[T]) Observable[T] {
		var (
			empty = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				empty = false
				obs.Next(v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				if empty {
					obs.Error(factory())
					return
				}
				obs.Complete()
			},
		)
	}
}
