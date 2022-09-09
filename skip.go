package rxgo

// Returns an Observable that skips the first count items emitted by the source Observable.
func Skip[T any](count uint) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				index++
				if count >= index {
					return
				}
				obs.Next(v)
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

// Skip a specified number of values before the completion of an observable.
func SkipLast[T any](skipCount uint) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			values = make([]T, 0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				values = append(values, v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				values = values[:uint(len(values))-skipCount]
				for _, v := range values {
					obs.Next(v)
				}
				obs.Complete()
			},
		)
	}
}

// Returns an Observable that skips items emitted by the source Observable until a
// second Observable emits an item.
func SkipUntil[A any, B any](notifier Observable[B]) OperatorFunc[A, A] {
	return func(source Observable[A]) Observable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			// var (
			// 	mu           = new(sync.RWMutex)
			// 	ok           bool
			// 	subscription Subscription
			// )
			// subscription = notifier.Subscribe(func(b B) {
			// 	mu.Lock()
			// 	ok = true
			// 	mu.Unlock()
			// }, func(err error) {}, func() {})

			// source.SubscribeSync(
			// 	func(v A) {
			// 		mu.RLock()
			// 		defer mu.RUnlock()
			// 		if ok {
			// 			subscriber.Next(v)
			// 		}
			// 	},
			// 	func(err error) {
			// 		unsubscribeStream(subscription)
			// 		subscriber.Error(err)
			// 	},
			// 	func() {
			// 		unsubscribeStream(subscription)
			// 		subscriber.Complete()
			// 	},
			// )
		})
	}
}

// Returns an Observable that skips all items emitted by the source Observable
// as long as a specified condition holds true, but emits all further source items
// as soon as the condition becomes false.
func SkipWhile[T any](predicate func(v T, index uint) bool) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
			pass  bool
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if pass {
					obs.Next(v)
					return
				}
				if !predicate(v, index) {
					pass = true
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
