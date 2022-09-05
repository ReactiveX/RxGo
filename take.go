package rxgo

// Emits only the first count values emitted by the source Observable.
func Take[T any](count uint) OperatorFunc[T, T] {
	if count == 0 {
		return func(source IObservable[T]) IObservable[T] {
			return EMPTY[T]()
		}
	}

	return func(source IObservable[T]) IObservable[T] {
		var (
			seen = uint(0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				seen++
				if seen <= count {
					obs.Next(v)
					if count <= seen {
						obs.Complete()
					}
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

// Waits for the source to complete, then emits the last N values from the source,
// as specified by the count argument.
func TakeLast[T any](count uint) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			values = make([]T, count)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if uint(len(values)) >= count {
					// shift the item from queue
					values = values[1:]
				}
				values = append(values, v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				for _, v := range values {
					obs.Next(v)
				}
				obs.Complete()
			},
		)
	}
}

// Emits the values emitted by the source Observable until a notifier Observable emits a value.
func TakeUntil[T any, R any](notifier IObservable[R]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				upStream     = source.subscribeOn()
				recv         = upStream.ForEach()
				notifyStream = notifier.subscribeOn()
				// stop         bool
			)

			for {
				select {
				case item, ok := <-recv:
					if !ok {
						return
					}
					select {
					case <-subscriber.Closed():
						return
					case subscriber.Send() <- item:
					}
				case <-notifyStream.ForEach():
					notifyStream.Stop()
					upStream.Stop()
					// select {
					// case <-subscriber.Closed():
					// 	return
					// case subscriber.Send() <- newComplete[T]():
					// 	log.Println("SSS")
					// 	// stop = true
					// }
				}
			}

		})
	}
}

// Emits values emitted by the source Observable so long as each value satisfies the given predicate,
// and then completes as soon as this predicate is not satisfied.
func TakeWhile[T any](predicate func(value T, index uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if !predicate(v, index) {
					obs.Complete()
					return
				}
				obs.Next(v)
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
