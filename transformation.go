package rxgo

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Buffers the source Observable values until closingNotifier emits.
func Buffer[T any, R any](closingNotifier Observable[R]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg     = new(sync.WaitGroup)
				buffer = make([]T, 0)
			)

			wg.Add(2)

			var (
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = closingNotifier.SubscribeOn(wg.Done)
			)

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					notifyStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						notifyStream.Stop()
						break observe
					}

					if err := item.Err(); err != nil {
						notifyStream.Stop()
						Error[[]T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						notifyStream.Stop()
						// Flush out all remaining buffer
						if len(buffer) >= 0 {
							Next(buffer).Send(subscriber)
						}
						Complete[[]T]().Send(subscriber)
						break observe
					}

					buffer = append(buffer, item.Value())

				case _, ok := <-notifyStream.ForEach():
					if !ok {
						upStream.Stop()
						break observe
					}

					Next(buffer).Send(subscriber)

					// Reset buffer
					buffer = make([]T, 0)
				}
			}

			wg.Wait()
		})
	}
}

// Buffers the source Observable values until the size hits the maximum bufferSize given.
func BufferCount[T any](bufferSize uint, startBufferEvery ...uint) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		var (
			buffers = [][]T{}
			// startFrom = bufferSize
		)
		// if len(startBufferEvery) > 0 {
		// 	startFrom = startBufferEvery[0]
		// }
		return createOperatorFunc(
			source,
			func(obs Observer[[]T], v T) {
				for idx := range buffers {
					buffers[idx] = append(buffers[idx], v)

					// if uint(len(buffers[idx])) >= bufferSize {

					// }
				}
				// count++
				// buffer = append(buffer, v)
				// // if len(buffer) >= bufCap {

				// // 	// Reset buffer
				// // 	buffer = nextBuffer()
				// // 	log.Println("Cap ->", cap(buffer))
				// // }
				// log.Println(count, startFrom)
				// if count >= startFrom {
				// 	obs.Next(buffer)
				// 	buffer = make([]T, 0, bufferSize)
				// }
			},
			func(obs Observer[[]T], err error) {
				obs.Error(err)
			},
			func(obs Observer[[]T]) {
				for _, b := range buffers {
					obs.Next(b)
				}
				obs.Complete()
			},
		)
	}
}

// Buffers the source Observable values for a specific time period.
func BufferTime[T any](bufferTimeSpan time.Duration) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				buffer   = make([]T, 0)
				upStream = source.SubscribeOn(wg.Done)
				timer    *time.Timer
			)

			setTimer := func() {
				if timer != nil {
					timer.Stop()
				}
				timer = time.NewTimer(bufferTimeSpan)
			}

			setTimer()

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						Error[[]T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						Next(buffer).Send(subscriber)
						Complete[[]T]().Send(subscriber)
						break observe
					}

					buffer = append(buffer, item.Value())

				case <-timer.C:
					Next(buffer).Send(subscriber)
					setTimer()
				}
			}

			// FIXME: I don't know how to stop timer
			timer.Stop()

			wg.Wait()
		})
	}
}

// Buffers the source Observable values starting from an emission from openings and ending
// when the output of closingSelector emits.
func BufferToggle[T any, O any](openings Observable[O], closingSelector func(value O) Observable[O]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				allowed     bool
				buffer      []T
				startStream = openings.SubscribeOn(wg.Done)
				upStream    = source.SubscribeOn(wg.Done)
				emitStream  Subscriber[O]
				stopCh      <-chan Notification[O]
			)

			setupValues := func() {
				allowed = false
				buffer = make([]T, 0)
				stopCh = make(<-chan Notification[O])
			}

			unsubscribeAll := func() {
				startStream.Stop()
				if emitStream != nil {
					emitStream.Stop()
				}
			}

			setupValues()

			// Buffers values from the source by opening the buffer via signals from an
			// Observable provided to openings, and closing and sending the buffers when a
			// Subscribable or Promise returned by the closingSelector function emits.
		observe:
			for {
				select {
				case <-subscriber.Closed():
					break observe

				case item, ok := <-startStream.ForEach():
					if !ok {
						break observe
					}

					allowed = true
					if emitStream != nil {
						// Unsubscribe the previous one
						emitStream.Stop()
					}
					wg.Add(1)
					emitStream = closingSelector(item.Value()).SubscribeOn(wg.Done)
					stopCh = emitStream.ForEach()

				case <-stopCh:
					Next(buffer).Send(subscriber)
					if emitStream != nil {
						emitStream.Stop()
					}
					setupValues()

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						unsubscribeAll()
						Error[[]T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						unsubscribeAll()
						Complete[[]T]().Send(subscriber)
						break observe
					}

					if allowed {
						buffer = append(buffer, item.Value())
					}
				}
			}

			wg.Wait()
		})
	}
}

// Buffers the source Observable values, using a factory function of closing Observables to
// determine when to close, emit, and reset the buffer.
func BufferWhen[T any, R any](closingSelector func() Observable[R]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				index         uint
				buffer        = make([]T, 0)
				upStream      = source.SubscribeOn(wg.Done)
				closingStream = closingSelector().SubscribeOn(wg.Done)
			)

			stopStreams := func() {
				upStream.Stop()
				closingStream.Stop()
			}

			onError := func(err error) {
				stopStreams()
				Error[[]T](err).Send(subscriber)
			}

			onComplete := func() {
				stopStreams()
				if len(buffer) > 0 {
					Next(buffer).Send(subscriber)
				}
				Complete[[]T]().Send(subscriber)
			}

		observe:
			for {
				select {
				case <-subscriber.Closed():
					stopStreams()
					break observe

				case item, ok := <-upStream.ForEach():
					// If the upstream closed, we break
					if !ok {
						stopStreams()
						break observe
					}

					if err := item.Err(); err != nil {
						onError(err)
						break observe
					}

					if item.Done() {
						onComplete()
						break observe
					}

					buffer = append(buffer, item.Value())
					index++

				case item, ok := <-closingStream.ForEach():
					if !ok {
						stopStreams()
						break observe
					}

					if err := item.Err(); err != nil {
						onError(err)
						break observe
					}

					if item.Done() {
						onComplete()
						break observe
					}

					Next(buffer).Send(subscriber)
					// Reset buffer values after sent
					buffer = make([]T, 0)
				}
			}

			wg.Wait()
		})
	}
}

// Projects each source value to an Observable which is merged in the output Observable,
// in a serialized fashion waiting for each one to complete before merging the next.
func ConcatMap[T any, R any](project ProjectionFunc[T, R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index    uint
				upStream = source.SubscribeOn(wg.Done)
			)

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					// If the upstream closed, we break
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						Error[R](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						Complete[R]().Send(subscriber)
						break observe
					}

					wg.Add(1)
					// we should wait the projection to complete
					stream := project(item.Value(), index).SubscribeOn(wg.Done)
				observeInner:
					for {
						select {
						case <-subscriber.Closed():
							upStream.Stop()
							stream.Stop()
							break observe

						case item, ok := <-stream.ForEach():
							if !ok {
								upStream.Stop()
								break observeInner
							}

							if err := item.Err(); err != nil {
								upStream.Stop()
								stream.Stop()
								item.Send(subscriber)
								break observe
							}

							if item.Done() {
								stream.Stop()
								break observeInner
							}

							item.Send(subscriber)
						}
					}

					index++
				}
			}
			wg.Wait()
		})
	}
}

// Projects each source value to an Observable which is merged in the output Observable
// only if the previous projected Observable has completed.
func ExhaustMap[T any, R any](project ProjectionFunc[T, R]) OperatorFunc[T, R] {
	if project == nil {
		panic(`rxgo: "ExhaustMap" expected project func`)
	}
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index      uint
				enabled    = new(atomic.Pointer[bool])
				upStream   = source.SubscribeOn(wg.Done)
				downStream Subscriber[R]
			)

			flag := true
			enabled.Store(&flag)

			stopDownStream := func() {
				if downStream != nil {
					downStream.Stop()
				}
			}

			observeStream := func(index uint, forEach <-chan Notification[R]) {
			observe:
				for {
					select {
					case <-subscriber.Closed():
						break observe

					case <-upStream.Closed():
						break observe

					case item, ok := <-forEach:
						if !ok {
							break observe
						}

						// TODO: handle error please

						if item.Done() {
							stopDownStream()
							flag := true
							enabled.Swap(&flag)
							break observe
						}

						item.Send(subscriber)
					}
				}
			}

		loop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					stopDownStream()
					break loop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break loop
					}

					if err := item.Err(); err != nil {
						stopDownStream()
						Error[R](err).Send(subscriber)
						break loop
					}

					if item.Done() {
						stopDownStream()
						Complete[R]().Send(subscriber)
						break loop
					}

					log.Println(item, ok)
					if *enabled.Load() {
						flag := false
						enabled.Swap(&flag)
						wg.Add(1)
						downStream = project(item.Value(), index).SubscribeOn(wg.Done)
						go observeStream(index, downStream.ForEach())
					}
					index++
				}
			}

			wg.Wait()
		})
	}
}

// Converts a higher-order Observable into a first-order Observable by dropping inner
// Observables while the previous inner Observable has not yet completed.
func ExhaustAll[T any]() OperatorFunc[Observable[T], T] {
	return ExhaustMap(func(value Observable[T], _ uint) Observable[T] {
		return value
	})
}

// Groups the items emitted by an Observable according to a specified criterion,
// and emits these grouped items as GroupedObservables, one GroupedObservable per group.
func GroupBy[T any, K comparable](keySelector func(value T) K) OperatorFunc[T, GroupedObservable[K, T]] {
	if keySelector == nil {
		panic(`rxgo: "GroupBy" expected keySelector func`)
	}
	return func(source Observable[T]) Observable[GroupedObservable[K, T]] {
		return newObservable(func(subscriber Subscriber[GroupedObservable[K, T]]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				key      K
				upStream = source.SubscribeOn(wg.Done)
				keySet   = make(map[K]Subject[T])
			)

		loop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break loop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break loop
					}

					log.Println(item)
					key = keySelector(item.Value())
					if _, exists := keySet[key]; !exists {
						keySet[key] = NewSubscriber[T]()
					} else {
						log.Println("HELLO")
						keySet[key].Send() <- item
						log.Println("HELLO END")
					}
				}
			}

			wg.Wait()
		})
	}
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func Map[T any, R any](mapper func(T, uint) (R, error)) OperatorFunc[T, R] {
	if mapper == nil {
		panic(`rxgo: "Map" expected mapper func`)
	}
	return func(source Observable[T]) Observable[R] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[R], v T) {
				output, err := mapper(v, index)
				index++
				if err != nil {
					obs.Error(err)
					return
				}
				obs.Next(output)
			},
			func(obs Observer[R], err error) {
				obs.Error(err)
			},
			func(obs Observer[R]) {
				obs.Complete()
			},
		)
	}
}

// Projects each source value to an Observable which is merged in the output Observable.
func MergeMap[T any, R any](project ProjectionFunc[T, R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index    uint
				upStream = source.SubscribeOn(wg.Done)
			)

			observeStream := func(index uint, v T) {
				var (
					downStream = project(v, index).SubscribeOn(wg.Done)
				)

			loop:
				for {
					select {
					case <-subscriber.Closed():
						downStream.Stop()
						break loop

					case item, ok := <-downStream.ForEach():
						if !ok {
							break loop
						}

						if item.Done() {
							break loop
						}

						if err := item.Err(); err != nil {
							break loop
						}

						item.Send(subscriber)
					}
				}
			}

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					// If the upstream closed, we break
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						Error[R](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						// Even upstream is done, we need to wait
						// downstream done as well
						break observe
					}

					// Everytime
					wg.Add(1)
					go observeStream(index, item.Value())
					index++
				}
			}

			wg.Wait()

			Complete[R]().Send(subscriber)
		})
	}
}

// Applies an accumulator function over the source Observable where the accumulator function
// itself returns an Observable, then each intermediate Observable returned is merged into
// the output Observable.
func MergeScan[V any, A any](accumulator func(acc A, value V, index uint) Observable[A], seed A, concurrent ...uint) OperatorFunc[V, A] {
	return func(source Observable[V]) Observable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index      uint
				finalValue = seed
				upStream   = source.SubscribeOn(wg.Done)
			)

			// MergeScan internally keeps the value of the acc parameter:
			// as long as the source Observable emits without inner Observable emitting,
			// the acc will be set to seed.

			observeStream := func() {
				Next(finalValue).Send(subscriber)
			}

		loop:
			for {
				select {
				case <-subscriber.Closed():
				case item, ok := <-upStream.ForEach():
					if !ok {
						break loop
					}

					wg.Add(1)
					accumulator(finalValue, item.Value(), index).SubscribeOn(wg.Done)
					observeStream()
					index++
				}
			}

			wg.Wait()
		})
	}
}

// Groups pairs of consecutive emissions together and emits them as an array of two values.
func PairWise[T any]() OperatorFunc[T, Tuple[T, T]] {
	return func(source Observable[T]) Observable[Tuple[T, T]] {
		var (
			result     = make([]T, 0, 2)
			noOfRecord int
		)
		return createOperatorFunc(
			source,
			func(obs Observer[Tuple[T, T]], v T) {
				result = append(result, v)
				noOfRecord = len(result)
				if noOfRecord >= 2 {
					obs.Next(NewTuple(result[0], result[1]))
					result = result[1:]
				}
			},
			func(obs Observer[Tuple[T, T]], err error) {
				obs.Error(err)
			},
			func(obs Observer[Tuple[T, T]]) {
				obs.Complete()
			},
		)
	}
}

// Useful for encapsulating and managing state. Applies an accumulator (or "reducer function")
// to each value from the source after an initial state is established --
// either via a seed value (second argument), or from the first value from the source.
func Scan[V any, A any](accumulator AccumulatorFunc[A, V], seed A) OperatorFunc[V, A] {
	if accumulator == nil {
		panic(`rxgo: "Scan" expected accumulator func`)
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
				obs.Next(result)
				index++
			},
			func(obs Observer[A], err error) {
				obs.Error(err)
			},
			func(obs Observer[A]) {
				obs.Complete()
			},
		)
	}
}

// Projects each source value to an Observable which is merged in the output Observable,
// emitting values only from the most recently projected Observable.
func SwitchMap[T any, R any](project func(value T, index uint) Observable[R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				index uint
				wg    = new(sync.WaitGroup)
				// mu    = new(sync.RWMutex)
				stop = make(chan struct{})
				// closing  = make(chan struct{})
				upStream = source.SubscribeOn(wg.Done)
				// stream   Subscriber[R]
			)

			wg.Add(1)

			closeStream := func() {
				close(stop)
				stop = make(chan struct{})
			}

			startStream := func(obs Observable[R]) {
				defer wg.Done()
				stream := obs.SubscribeOn()
				defer stream.Stop()

			loop:
				for {
					select {
					case <-stop:
						break loop

					case <-subscriber.Closed():
						stream.Stop()
						return

					case item, ok := <-stream.ForEach():
						if !ok {
							break loop
						}

						item.Send(subscriber)
					}
				}
			}

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					closeStream()

				case item := <-upStream.ForEach():
					if err := item.Err(); err != nil {
						break observe
					}

					if item.Done() {
						break observe
					}

					closeStream()
					wg.Add(1)
					go startStream(project(item.Value(), index))
					index++
				}
			}

			wg.Wait()
		})
	}
}
