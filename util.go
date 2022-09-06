package rxgo

func createOperatorFunc[T any, R any](
	source IObservable[T],
	onNext func(Observer[R], T),
	onError func(Observer[R], error),
	onComplete func(Observer[R]),
) IObservable[R] {
	return newObservable(func(subscriber Subscriber[R]) {
		var (
			// terminated = subscriber.Closed()
			stop bool

			// input stream
			upStream = source.SubscribeOn()
		)

		obs := &consumerObserver[R]{
			onNext: func(v R) {
				select {
				case <-subscriber.Closed():
					stop = true
					return
				case subscriber.Send() <- NextNotification(v):
				}
			},
			onError: func(err error) {
				upStream.Stop()
				stop = true
				select {
				case <-subscriber.Closed():
					return
				case subscriber.Send() <- ErrorNotification[R](err):
				}
			},
			onComplete: func() {
				// Inform the up stream to stop emit value
				upStream.Stop()
				stop = true
				select {
				case <-subscriber.Closed():
					return
				case subscriber.Send() <- CompleteNotification[R]():
				}
			},
		}

		for !stop {
			select {
			// If only the stream terminated, break it
			case <-subscriber.Closed():
				stop = true
				upStream.Stop()
				return

			case item, ok := <-upStream.ForEach():
				if !ok {
					// If only the data stream closed, break it
					stop = true
					return
				}

				if err := item.Err(); err != nil {
					onError(obs, err)
					return
				}

				if item.Done() {
					onComplete(obs)
					return
				}

				onNext(obs, item.Value())
			}
		}
	})
}
