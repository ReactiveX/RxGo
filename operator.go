package rxgo

import (
	"log"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// Emits the single value at the specified index in a sequence of emissions from the source Observable.
func ElementAt[T any](pos uint, defaultValue ...T) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index    uint
			notEmpty bool
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if index == pos {
					obs.Next(v)
					obs.Complete()
					notEmpty = true
					return
				}
				index++
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				if notEmpty {
					return
				}

				if len(defaultValue) > 0 {
					obs.Next(defaultValue[0])
					obs.Complete()
					return
				}

				obs.Error(ErrArgumentOutOfRange)
			},
		)
	}
}

// Emits only the first value (or the first value that meets some condition)
// emitted by the source Observable.
func First[T any](predicate PredicateFunc[T], defaultValue ...T) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index    uint
			hasValue bool
		)
		cb := skipPredicate[T]
		if predicate != nil {
			cb = predicate
		}
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if !hasValue && cb(v, index) {
					hasValue = true
					obs.Next(v)
					obs.Complete()
					return
				}
				index++
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				if !hasValue {
					if len(defaultValue) == 0 {
						obs.Error(ErrEmpty)
						return
					}

					obs.Next(defaultValue[0])
				}
				obs.Complete()
			},
		)
	}
}

// Returns an Observable that emits only the last item emitted by the source Observable.
// It optionally takes a predicate function as a parameter, in which case,
// rather than emitting the last item from the source Observable,
// the resulting Observable will emit the last item from the source Observable
// that satisfies the predicate.
func Last[T any](predicate PredicateFunc[T], defaultValue ...T) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index       uint
			hasValue    bool
			latestValue T
			found       bool
		)
		cb := skipPredicate[T]
		if predicate != nil {
			cb = predicate
		}
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				hasValue = true
				if cb(v, index) {
					found = true
					latestValue = v
				}
				index++
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				if found {
					obs.Next(latestValue)
				} else {
					if !hasValue {
						if len(defaultValue) == 0 {
							obs.Error(ErrEmpty)
							return
						}

						obs.Next(defaultValue[0])
					} else {
						obs.Error(ErrNotFound)
						return
					}
				}
				obs.Complete()
			},
		)
	}
}

// Ignores all items emitted by the source Observable and only passes calls of complete or error.
func IgnoreElements[T any]() OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Complete()
			},
		)
	}
}

// Returns an Observable that will resubscribe to the source
// stream when the source stream completes.
func Repeat[T any, N constraints.Unsigned](count N) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			// source.SubscribeSync(
			// 	func(t T) {
			// 		for i := N(0); i < count; i++ {
			// 			subscriber.Next(t)
			// 		}
			// 	},
			// 	subscriber.Error,
			// 	subscriber.Complete,
			// )
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

// Emits the most recently emitted value from the source Observable whenever
// another Observable, the notifier, emits.
func Sample[A any, B any](notifier Observable[B]) OperatorFunc[A, A] {
	return func(source Observable[A]) Observable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			var (
				wg           = new(sync.WaitGroup)
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = notifier.SubscribeOn(wg.Done)
				latestValue  = Next(*new(A))
			)

			wg.Add(2)

		observe:
			for {
				select {
				case <-subscriber.Closed():
					return
				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if item.Done() {
						notifyStream.Stop()
						subscriber.Send() <- item
						break observe
					}

					latestValue = item
				case <-notifyStream.ForEach():
					subscriber.Send() <- latestValue
				}
			}

			wg.Wait()
			log.Println("ALL DONE")
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

// Emits a value from the source Observable, then ignores subsequent source values
// for duration milliseconds, then repeats this process.
func Throttle[T any, R any](durationSelector func(v T) Observable[R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				durationSelector(v)
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

// Emits a notification from the source Observable only after a particular time span
// has passed without another source emission.
func DebounceTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			timer *time.Timer
		)
		// https://github.com/ReactiveX/RxGo/blob/35328a75073980197d938cf235158a0654024de5/observable_operator.go#L670
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(duration, func() {
					obs.Next(v)
				})
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

// Catches errors on the observable to be handled by returning a new observable
// or throwing an error.
func CatchError[T any](catch func(error, Observable[T]) Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
				// subscription Subscription
				// subscribe func(Observable[T])
			)

			// unsubscribe := func() {
			// 	if subscription != nil {
			// 		subscription.Unsubscribe()
			// 	}
			// 	subscription = nil
			// }

			// subscribe = func(stream Observable[T]) {
			// 	subscription = stream.Subscribe(
			// 		subscriber.Next,
			// 		func(err error) {
			// 			obs := catch(err, source)
			// 			unsubscribe() // unsubscribe the previous stream and start another one
			// 			subscribe(obs)
			// 		},
			// 		func() {
			// 			unsubscribe()
			// 			wg.Done()
			// 		},
			// 	)
			// }

			wg.Add(1)
			// subscribe(source)
			wg.Wait()

			subscriber.Send() <- Complete[T]()
		})
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
