package rxgo

func Take[T any](count uint) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		if count == 0 {
			return EMPTY[T]()
		}

		seen := uint(0)
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(v T) {
					seen++
					if seen <= count {
						subscriber.Next(v)
						if count <= seen {
							subscriber.Complete()
						}
					}
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

func Filter[T any](filter func(T, uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		index := uint(0)
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(v T) {
					if filter(v, index) {
						subscriber.Next(v)
					}
					index++
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

func Map[T any, R any](mapper func(T, uint) R) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		index := uint(0)
		return newObservable(func(subscriber Subscriber[R]) {
			source.SubscribeSync(
				func(v T) {
					output := mapper(v, index)
					subscriber.Next(output)
					index++
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}
