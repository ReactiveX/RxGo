package rxgo

import (
	"errors"
	"sync"
	"sync/atomic"
)

func CombineLatestAll() {
	// TODO: implement `CombineLatestAll`
}

// Create an observable that combines the latest values from all passed observables
// and the source into arrays and emits them.
func CombineLatestWith[T any](sources ...Observable[T]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		sources = append([]Observable[T]{source}, sources...)
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				noOfSource   = len(sources)
				emitCount    = new(atomic.Uint32)
				errOnce      = new(atomic.Pointer[error])
				errCh        = make(chan error, 1)
				stopCh       = make(chan struct{})
				wg           = new(sync.WaitGroup)
				latestValues = make([]T, noOfSource)
			)

			wg.Add(noOfSource)

			// To ensure the output array always has the same length,
			// combineLatest will actually wait for all input Observables
			// to emit at least once, before it starts emitting results.
			onNext := func() {
				if emitCount.Load() == uint32(noOfSource) {
					Next(latestValues).Send(subscriber)
				}
			}

			observeStream := func(index int, upStream Subscriber[T]) {
				var (
					emitted bool
				)

			loop:
				for errOnce.Load() != nil {
					select {
					case <-subscriber.Closed():
						upStream.Stop()
						break loop

					case <-stopCh:
						upStream.Stop()
						break loop

					case item, ok := <-upStream.ForEach():
						if !ok {
							break loop
						}

						if err := item.Err(); err != nil {
							errOnce.CompareAndSwap(nil, &err)
							break loop
						}

						if item.Done() {
							break loop
						}

						//  Passing an empty array will result in an Observable that completes immediately.
						if !emitted {
							emitCount.Add(1)
							emitted = true
						}
						latestValues[index] = item.Value()
						onNext()
					}
				}
			}

			go func() {
				select {
				case <-subscriber.Closed():
					return
				case _, ok := <-errCh:
					if !ok {
						return
					}
					close(stopCh)
				}
			}()

			for i, source := range sources {
				subscriber := source.SubscribeOn(wg.Done)
				go observeStream(i, subscriber)
			}

			wg.Wait()

			select {
			case <-errCh:
			default:
				// Close error channel gracefully
				close(errCh)
			}

			Complete[[]T]().Send(subscriber)
		})
	}
}

// Converts a higher-order Observable into a first-order Observable by  concatenating the inner Observables in order.
func ConcatAll[T any]() OperatorFunc[Observable[T], T] {
	return func(source Observable[Observable[T]]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index      uint
				upStream   = source.SubscribeOn(wg.Done)
				downStream Subscriber[T]
			)

			unsubscribeAll := func() {
				upStream.Stop()
				if downStream != nil {
					downStream.Stop()
				}
			}

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					// If the upstream closed, we break
					if !ok {
						break outerLoop
					}

					if err := item.Err(); err != nil {
						Error[T](err).Send(subscriber)
						break outerLoop
					}

					if item.Done() {
						Complete[T]().Send(subscriber)
						break outerLoop
					}

					wg.Add(1)
					// we should wait the projection to complete
					downStream = item.Value().SubscribeOn(wg.Done)

				innerLoop:
					for {
						select {
						case <-subscriber.Closed():
							unsubscribeAll()
							break outerLoop

						case item, ok := <-downStream.ForEach():
							if !ok {
								break innerLoop
							}

							if err := item.Err(); err != nil {
								unsubscribeAll()
								item.Send(subscriber)
								break outerLoop
							}

							if item.Done() {
								// downStream.Stop()
								break innerLoop
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

// Emits all of the values from the source observable, then, once it completes, subscribes to each observable source provided, one at a time, emitting all of their values, and not subscribing to the next one until it completes.
func ConcatWith[T any](sources ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		sources = append([]Observable[T]{source}, sources...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg  = new(sync.WaitGroup)
				err error
			)

		outerLoop:
			for len(sources) > 0 {
				wg.Add(1)
				firstSource := sources[0]
				upStream := firstSource.SubscribeOn(wg.Done)

			innerLoop:
				for {
					select {
					case <-subscriber.Closed():
						upStream.Stop()
						return

					case item, ok := <-upStream.ForEach():
						if !ok {
							break innerLoop
						}

						if item.Done() {
							// start another loop
							break innerLoop
						}

						if err = item.Err(); err != nil {
							break outerLoop
						}

						item.Send(subscriber)
					}
				}

				sources = sources[1:]
			}

			if err != nil {
				Error[T](err).Send(subscriber)
			} else {
				Complete[T]().Send(subscriber)
			}

			wg.Wait()
		})
	}
}

// FIXME: Accepts an Array of ObservableInput or a dictionary Object of ObservableInput
// and returns an Observable that emits either an array of values in the exact same
// order as the passed array, or a dictionary of values in the same shape as the
// passed dictionary.
func ForkJoin[T any](sources ...Observable[T]) Observable[[]T] {
	return newObservable(func(subscriber Subscriber[[]T]) {
		var (
			noOfSource = len(sources)
		)
		// forkJoin is an operator that takes any number of input observables which can be
		// passed either as an array or a dictionary of input observables. If no input
		// observables are provided (e.g. an empty array is passed), then the resulting
		// stream will complete immediately.
		if noOfSource < 1 {
			Complete[[]T]().Send(subscriber)
			return
		}

		var (
			wg  = new(sync.WaitGroup)
			err = new(atomic.Pointer[error])
			// Single buffered channel will not accept more than one signal
			errCh         = make(chan error, 1)
			stopCh        = make(chan struct{})
			latestValues  = make([]T, noOfSource)
			subscriptions = make([]Subscriber[T], noOfSource)
		)

		wg.Add(noOfSource)

		go func() {
			select {
			case <-subscriber.Closed():
				return
			case v, ok := <-errCh:
				if !ok {
					return
				}
				err.Swap(&v)
				close(stopCh)
			}
		}()

		// In order for the resulting array to have the same length as the number of
		// input observables, whenever any of the given observables completes without
		// emitting any value, forkJoin will complete at that moment as well and it
		// will not emit anything either, even if it already has some last values
		// from other observables.
		onNext := func(index int, v T) {
			// mu.Lock()
			// defer mu.Unlock()
			latestValues[index] = v
		}

		observeStream := func(index int, upStream Subscriber[T]) {
		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case <-stopCh:
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					// if one error, everything error
					if err := item.Err(); err != nil {
						errCh <- err // FIXME: data race
						break observe
					}

					if item.Done() {
						break observe
					}

					// forkJoin will wait for all passed observables to emit and complete
					// and then it will emit an array or an object with last values from
					// corresponding observables.
					onNext(index, item.Value())
				}
			}
		}

		for i, source := range sources {
			subscriber := source.SubscribeOn(wg.Done)
			subscriptions[i] = subscriber
			go observeStream(i, subscriber)
		}

		wg.Wait()

		// Remove dangling go-routine
		select {
		case <-errCh:
		default:
			// Close error channel gracefully
			close(errCh)
		}

		for _, sub := range subscriptions {
			sub.Stop()
		}

		if exception := err.Load(); exception != nil {
			if errors.Is(*exception, ErrEmpty) {
				Complete[[]T]().Send(subscriber)
				return
			}

			Error[[]T](*exception).Send(subscriber)
			return
		}

		Next(latestValues).Send(subscriber)
		Complete[[]T]().Send(subscriber)
	})
}

// FIXME: Merge the values from all observables to a single observable result.
func MergeWith[T any](input Observable[T], inputs ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		inputs = append([]Observable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg                  = new(sync.WaitGroup)
				mu                  = new(sync.RWMutex)
				activeSubCount      = new(atomic.Int32)
				noOfInputs          = len(inputs)
				activeSubscriptions = make([]Subscriber[T], noOfInputs)
				err                 = new(atomic.Pointer[error])
				stopCh              = make(chan struct{})
				errCh               = make(chan error, 1)
			)

			onError := func(err error) {
				mu.Lock()
				defer mu.Unlock()
				select {
				case errCh <- err:
				default:
				}
			}

			go func() {
				select {
				case <-subscriber.Closed():
					return
				case v, ok := <-errCh:
					if !ok {
						return
					}
					err.Swap(&v)
					close(stopCh)
				}
			}()

			observeStream := func(index int, stream Subscriber[T]) {
				defer activeSubCount.Add(-1)

			observe:
				for {
					select {
					case <-subscriber.Closed():
						stream.Stop()
						break observe

					case <-stopCh:
						stream.Stop()
						break observe

					case item, ok := <-stream.ForEach():
						if !ok {
							break observe
						}

						if err := item.Err(); err != nil {
							onError(err)
							break observe
						}

						if item.Done() {
							break observe
						}

						item.Send(subscriber)
					}
				}
			}

			wg.Add(noOfInputs)
			// activeSubCount.Store(int32(noOfInputs))

			for i, input := range inputs {
				activeSubscriptions[i] = input.SubscribeOn(wg.Done)
				go observeStream(i, activeSubscriptions[i])
			}

			wg.Wait()

			// Remove dangling go-routine
			select {
			case <-errCh:
			default:
				mu.Lock()
				// Close error channel gracefully
				close(errCh)
				mu.Unlock()
			}

			// Stop all stream
			for _, sub := range activeSubscriptions {
				sub.Stop()
			}

			if exception := err.Load(); exception != nil {
				Error[T](*exception).Send(subscriber)
				return
			}

			Complete[T]().Send(subscriber)
		})
	}
}

// FIXME: Creates an Observable that mirrors the first source Observable to emit a
// next, error or complete notification from the combination of the Observable
// to which the operator is applied and supplied Observables.
func RaceWith[T any](input Observable[T], inputs ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		inputs = append([]Observable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg         = new(sync.WaitGroup)
				noOfInputs = len(inputs)
				// fastest             = int(-1)
				stopCh              = make(chan int)
				activeSubscriptions = make([]Subscriber[T], noOfInputs)
				// mu                  = new(sync.RWMutex)
				// unsubscribed        bool
			)

			wg.Add(noOfInputs)

			// unsubscribeAll := func() {
			// 	mu.Lock()
			// 	for _, v := range activeSubscriptions {
			// 		v.Stop()
			// 		log.Println(v)
			// 	}
			// 	mu.Unlock()
			// }

			benchmarkStream := func(index int, stream Subscriber[T]) {
				defer stream.Stop()

			observe:
				for {
					select {
					case <-subscriber.Closed():
						break observe

					case v := <-stopCh:
						if v != index {
							break observe
						}

						// case item, ok := <-stream.ForEach():
						// 	if !ok {
						// 		break observe
						// 	}

						// mu.Lock()
						// select {
						// case fastestCh <- index:
						// 	// Inform I'm the winner
						// default:
						// 	stream.Stop()
						// 	break observe
						// }

						// mu.Lock()
						// defer mu.Unlock()
						// for _, sub := range activeSubscriptions {
						// 	sub.Stop()
						// }
						// activeSubscriptions = []Subscriber[T]{}
						// log.Println("ForEach ah", index, item)
						// fastestCh <- idx
						// obs.Stop()
					}
				}
			}

			for i, v := range inputs {
				activeSubscriptions[i] = v.SubscribeOn(wg.Done)
				go benchmarkStream(i, activeSubscriptions[i])
			}

			// unsubscribeAll()

			// stopCh <- <-fastestCh
			// for {
			// select {
			// case stopCh <- <-fastestCh:

			// 	// log.Println(item, ok)
			// }
			// }

			wg.Wait()

			Complete[T]().Send(subscriber)
		})
	}
}

// Combines multiple Observables to create an Observable whose values are calculated
// from the values, in order, of each of its input Observables.
func ZipWith[T any](input Observable[T], inputs ...Observable[T]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		inputs = append([]Observable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg         = new(sync.WaitGroup)
				noOfSource = uint(len(inputs))
				observers  = make([]Subscriber[T], 0, noOfSource)
			)

			wg.Add(int(noOfSource))

			for _, input := range inputs {
				observers = append(observers, input.SubscribeOn(wg.Done))
			}

			unsubscribeAll := func() {
				for _, obs := range observers {
					obs.Stop()
				}
			}

			var (
				result    []T
				completed uint
			)

			setupValues := func() {
				result = make([]T, noOfSource)
				completed = 0
			}

			setupValues()

		outerLoop:
			for {

			innerLoop:
				for i, obs := range observers {
					select {
					case <-subscriber.Closed():
						unsubscribeAll()
						break outerLoop

					case item, ok := <-obs.ForEach():
						if !ok || item.Done() {
							completed++
							unsubscribeAll()
							break innerLoop
						}

						if err := item.Err(); err != nil {
							unsubscribeAll()
							Error[[]T](err).Send(subscriber)
							break outerLoop
						}

						if item != nil {
							result[i] = item.Value()
						}
					}
				}

				// Any of the stream completed, we will escape
				if completed > 0 {
					Complete[[]T]().Send(subscriber)
					break outerLoop
				}

				Next(result).Send(subscriber)

				// Reset the values for next loop
				setupValues()
			}

			wg.Wait()
		})
	}
}

// Collects all observable inner sources from the source, once the source
// completes, it will subscribe to all inner sources, combining their
// values by index and emitting them.
func ZipAll[T any]() OperatorFunc[Observable[T], []T] {
	return func(source Observable[Observable[T]]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				runNext     bool
				upStream    = source.SubscribeOn(wg.Done)
				observables = make([]Observable[T], 0)
			)

			// Collects all observable inner sources from the source, once the
			// source completes, it will subscribe to all inner sources,
			// combining their values by index and emitting them.
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

					if err := item.Err(); err != nil {
						Error[[]T](err).Send(subscriber)
						break loop
					}

					if item.Done() {
						runNext = true
						break loop
					}

					observables = append(observables, item.Value())
				}
			}

			var noOfObservables = len(observables)
			if runNext && noOfObservables > 0 {
				var (
					observers = make([]Subscriber[T], 0, noOfObservables)
					result    []T
					completed uint
				)

				setupValues := func() {
					result = make([]T, noOfObservables)
					completed = 0
				}

				unsubscribeAll := func() {
					for _, obs := range observers {
						obs.Stop()
					}
				}

				setupValues()
				wg.Add(noOfObservables)
				for _, obs := range observables {
					observers = append(observers, obs.SubscribeOn(wg.Done))
				}

			outerLoop:
				for {
				innerLoop:
					for i, obs := range observers {
						select {
						case <-subscriber.Closed():
							unsubscribeAll()
							break outerLoop

						case item, ok := <-obs.ForEach():
							if !ok || item.Done() {
								completed++
								unsubscribeAll()
								break innerLoop
							}

							if err := item.Err(); err != nil {
								unsubscribeAll()
								Error[[]T](err).Send(subscriber)
								break outerLoop
							}

							if item != nil {
								result[i] = item.Value()
							}
						}
					}

					// Any of the stream completed, we will escape
					if completed > 0 {
						Complete[[]T]().Send(subscriber)
						break outerLoop
					}

					Next(result).Send(subscriber)

					// Reset the values for next loop
					setupValues()
				}
			}

			wg.Wait()
		})
	}
}
