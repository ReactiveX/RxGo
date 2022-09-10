package rxgo

import (
	"sync"
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
						break observe
					}

					buffer = append(buffer, item.Value())

				case _, ok := <-notifyStream.ForEach():
					if !ok {
						upStream.Stop()
						break observe
					}

					if !Next(buffer).Send(subscriber) {
						break observe
					}

					// Reset buffer
					buffer = make([]T, 0)
				}
			}

			wg.Wait()

			Complete[[]T]().Send(subscriber)
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

// Projects each source value to an Observable which is merged in the output Observable,
// in a serialized fashion waiting for each one to complete before merging the next.
func ConcatMap[T any, R any](project func(value T, index uint) Observable[R]) OperatorFunc[T, R] {
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
				keySet   = make(map[K]*groupedObservable[K, T])
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

					key = keySelector(item.Value())
					if _, exists := keySet[key]; !exists {
						keySet[key] = newGroupedObservable[K, T]()
					} else {
						keySet[key].connector.Send() <- item
					}
				}
			}

			wg.Wait()
		})
	}
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func Map[T any, R any](mapper func(T, uint) (R, error)) OperatorFunc[T, R] {
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
func MergeMap[T any, R any](project func(value T, index uint) Observable[R]) OperatorFunc[T, R] {
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
						return

					case item, ok := <-downStream.ForEach():
						if !ok {
							return
						}

						if item.Done() {
							break loop
						}

						if err := item.Err(); err != nil {
							break loop
						}

						item.Send(subscriber)
						// log.Println("Item ->", item)
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

			// if errCount < 1 {
			// 	Complete[R]().Send(subscriber)
			// }
		})
	}
}

// Applies an accumulator function over the source Observable where the accumulator
// function itself returns an Observable, then each intermediate Observable returned
// is merged into the output Observable.
func MergeScan[T any, R any](accumulator func(acc R, value T, index uint) Observable[R], seed R) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return nil
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
