package rxgo

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// Returns an Observable that will resubscribe to the source stream when the source stream
// completes.
func Repeat[T any, N constraints.Unsigned](count N) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			var (
				repeated = N(0)
				upStream Subscriber[T]
				forEach  <-chan Notification[T]
			)

			setupStream := func() {
				wg.Add(1)
				upStream = source.SubscribeOn(wg.Done)
				forEach = upStream.ForEach()
			}

			setupStream()

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-forEach:
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						Error[T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						repeated++
						if repeated < count {
							setupStream()
							continue
						}

						Complete[T]().Send(subscriber)
						break observe
					}

					item.Send(subscriber)
				}
			}

			wg.Wait()
		})
	}
}

// Used to perform side-effects for notifications from the source observable
func Tap[T any](cb Observer[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		if cb == nil {
			cb = NewObserver[T](nil, nil, nil)
		}
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				obs.Next(v)
				cb.Next(v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
				cb.Error(err)
			},
			func(obs Observer[T]) {
				obs.Complete()
				cb.Complete()
			},
		)
	}
}

// Returns an observable that asserts that only one value is emitted from the observable
// that matches the predicate. If no predicate is provided, then it will assert that the
// observable only emits one value.
// FIXME: should rename `Single2` to `Single`
func Single2[T any](predicate func(v T, index uint) bool) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			// var (
			// 	index    uint
			// 	found    bool
			// 	matches  uint
			// 	hasValue bool
			// )
			// source.SubscribeSync(
			// 	func(v T) {
			// 		result := predicate(v, index)
			// 		if result {
			// 			found = result
			// 			matches++
			// 		}
			// 		hasValue = true
			// 		index++
			// 	},
			// 	subscriber.Error,
			// 	func() {
			// 		if !hasValue {
			// 			subscriber.Error(ErrEmpty)
			// 		} else if !found {
			// 			subscriber.Error(ErrNotFound)
			// 		} else if matches > 1 {
			// 			subscriber.Error(ErrSequence)
			// 		}
			// 		subscriber.Complete()
			// 	},
			// )
		})
	}
}

// Delays the emission of items from the source Observable by a given timeout.
func Delay[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				time.Sleep(duration)
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

// Delays the emission of items from the source Observable by a given time span
// determined by the emissions of another Observable.
func DelayWhen[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				time.Sleep(duration)
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

// Combines the source Observable with other Observables to create an Observable
// whose values are calculated from the latest values of each, only when the source emits.
func WithLatestFrom[A any, B any](input Observable[B]) OperatorFunc[A, Tuple[A, B]] {
	return func(source Observable[A]) Observable[Tuple[A, B]] {
		return newObservable(func(subscriber Subscriber[Tuple[A, B]]) {
			var (
				allOk              [2]bool
				activeSubscription = 2
				wg                 = new(sync.WaitGroup)
				upStream           = source.SubscribeOn(wg.Done)
				notifySteam        = input.SubscribeOn(wg.Done)
				latestA            A
				latestB            B
			)

			wg.Add(activeSubscription)

			stopAll := func() {
				upStream.Stop()
				notifySteam.Stop()
				activeSubscription = 0
			}

			onNext := func() {
				if allOk[0] && allOk[1] {
					subscriber.Send() <- Next(NewTuple(latestA, latestB))
				}
			}

			// All input Observables must emit at least one value before
			// the output Observable will emit a value.
			for activeSubscription > 0 {
				select {
				case <-subscriber.Closed():
					stopAll()

				case item := <-notifySteam.ForEach():
					if item == nil {
						continue
					}

					allOk[1] = true
					if item.Done() {
						activeSubscription--
						continue
					}

					if err := item.Err(); err != nil {
						stopAll()
						subscriber.Send() <- Error[Tuple[A, B]](err)
						continue
					}

					latestB = item.Value()
					onNext()

				case item := <-upStream.ForEach():
					if item == nil {
						continue
					}

					allOk[0] = true
					if item.Done() {
						activeSubscription--
						continue
					}

					if err := item.Err(); err != nil {
						stopAll()
						subscriber.Send() <- Error[Tuple[A, B]](err)
						continue
					}

					latestA = item.Value()
					onNext()
				}
			}

			wg.Wait()
		})
	}
}

type TimeoutConfig[T any] struct {
	With func() Observable[T]
	Each time.Duration
}

type timeoutConfig[T any] interface {
	time.Duration | TimeoutConfig[T]
}

// Errors if Observable does not emit a value in given time span.
// FIXME:  DATA RACE and send on closed channel
func Timeout[T any, C timeoutConfig[T]](config C) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
				timeout  <-chan time.Time
			)

			switch v := any(config).(type) {
			case time.Duration:
				timeout = time.After(v)
			case TimeoutConfig[T]:
				panic("unimplemented")
			}

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

					// Reset timeout
					timeout = make(<-chan time.Time)
					item.Send(subscriber)
					if item.Err() != nil || item.Done() {
						break loop
					}

				case <-timeout:
					upStream.Stop()
					Error[T](ErrTimeout).Send(subscriber)
					break loop
				}
			}

			wg.Wait()
		})
	}
}

// Collects all source emissions and emits them as an array when the source completes.
func ToArray[T any]() OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		var (
			result = make([]T, 0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[[]T], v T) {
				result = append(result, v)
			},
			func(obs Observer[[]T], err error) {
				// When the source Observable errors no array will be emitted.
				obs.Error(err)
			},
			func(obs Observer[[]T]) {
				obs.Next(result)
				obs.Complete()
			},
		)
	}
}
