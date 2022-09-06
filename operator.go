package rxgo

import (
	"log"
	"reflect"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// Emits the single value at the specified index in a sequence of emissions from the source Observable.
func ElementAt[T any](pos uint, defaultValue ...T) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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
func First[T any](predicate func(value T, index uint) bool, defaultValue ...T) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			index    uint
			hasValue bool
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if !hasValue && predicate != nil && predicate(v, index) {
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
func Last[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return Pipe1(source, TakeLast[T](1))
	}
}

// Emits only the first value emitted by the source Observable that meets some condition.
func Find[T any](predicate func(T, uint) bool) OperatorFunc[T, Optional[T]] {
	return func(source IObservable[T]) IObservable[Optional[T]] {
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
func FindIndex[T any](predicate func(T, uint) bool) OperatorFunc[T, int] {
	var (
		index uint
		found bool
	)
	return func(source IObservable[T]) IObservable[int] {
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

// The Min operator operates on an Observable that emits numbers
// (or items that can be compared with a provided function),
// and when source Observable completes it emits a single item: the item with the smallest value.
func Min[T any](comparer func(a T, b T) int8) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			lastValue T
			first     = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if first {
					lastValue = v
					first = false
					return
				}

				switch comparer(lastValue, v) {
				case 1:
					lastValue = v
				default:
				}
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Next(lastValue)
				obs.Complete()
			},
		)
	}
}

// The Max operator operates on an Observable that emits numbers
// (or items that can be compared with a provided function),
// and when source Observable completes it emits a single item: the item with the largest value.
func Max[T any](comparer func(a T, b T) int8) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			lastValue T
			first     = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if first {
					lastValue = v
					first = false
					return
				}

				switch comparer(lastValue, v) {
				case -1:
					lastValue = v
				default:
					lastValue = v
				}
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				obs.Next(lastValue)
				obs.Complete()
			},
		)
	}
}

// Counts the number of emissions on the source and emits that number when the source completes.
func Count[T any](predicate ...func(v T, index uint) bool) OperatorFunc[T, uint] {
	cb := func(T, uint) bool {
		return true
	}
	if len(predicate) > 0 {
		cb = predicate[0]
	}

	return func(source IObservable[T]) IObservable[uint] {
		var (
			count uint
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[uint], v T) {
				if cb(v, index) {
					count++
				}
				index++
			},
			func(obs Observer[uint], err error) {
				obs.Error(err)
			},
			func(obs Observer[uint]) {
				obs.Next(count)
				obs.Complete()
			},
		)
	}
}

// Ignores all items emitted by the source Observable and only passes calls of complete or error.
func IgnoreElements[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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

// Returns an Observable that emits whether or not every item of the
// source satisfies the condition specified.
func Every[T any](predicate func(value T, index uint) bool) OperatorFunc[T, bool] {
	return func(source IObservable[T]) IObservable[bool] {
		var (
			allOk = true
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[bool], v T) {
				allOk = allOk && predicate(v, index)
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

// Returns an Observable that will resubscribe to the source
// stream when the source stream completes.
func Repeat[T any, N constraints.Unsigned](count N) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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

// Emits false if the input Observable emits any values,
// or emits true if the input Observable completes without emitting any values.
func IsEmpty[T any]() OperatorFunc[T, bool] {
	return func(source IObservable[T]) IObservable[bool] {
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

// Emits a given value if the source Observable completes without emitting any
// next value, otherwise mirrors the source Observable.
func DefaultIfEmpty[T any](defaultValue T) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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

// Returns a result Observable that emits all values pushed by the source observable
// if they are distinct in comparison to the last value the result observable emitted.
func DistinctUntilChanged[T any](comparator ...func(prev T, current T) bool) OperatorFunc[T, T] {
	cb := func(prev T, current T) bool {
		return reflect.DeepEqual(prev, current)
	}
	if len(comparator) > 0 {
		cb = comparator[0]
	}
	return func(source IObservable[T]) IObservable[T] {
		var (
			lastValue T
			first     = true
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if first || !cb(lastValue, v) {
					obs.Next(v)
					first = false
					lastValue = v
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

// Filter emits only those items from an Observable that pass a predicate test.
func Filter[T any](filter func(T, uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if filter(v, index) {
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

// Map transforms the items emitted by an Observable by applying a function to each item.
func Map[T any, R any](mapper func(T, uint) (R, error)) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
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

// Used to perform side-effects for notifications from the source observable
func Tap[T any](cb Observer[T]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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
	return func(source IObservable[T]) IObservable[T] {
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
func Sample[A any, B any](notifier IObservable[B]) OperatorFunc[A, A] {
	return func(source IObservable[A]) IObservable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			var (
				wg           = new(sync.WaitGroup)
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = notifier.SubscribeOn(wg.Done)
				latestValue  = NextNotification(*new(A))
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

// Projects each source value to an Observable which is merged in the output Observable,
// in a serialized fashion waiting for each one to complete before merging the next.
func ConcatMap[T any, R any](project func(value T, index uint) IObservable[R]) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			// var (
			// 	index              uint
			// 	buffer             = make([]T, 0)
			// 	concurrent         = uint(1)
			// 	isComplete         bool
			// 	activeSubscription uint
			// )

			// checkComplete := func() {
			// 	if isComplete && len(buffer) <= 0 {
			// 		subscriber.Complete()
			// 	}
			// }

			// var innerNext func(T)
			// innerNext = func(outerV T) {
			// 	activeSubscription++

			// 	stream := project(outerV, index)
			// 	index++

			// 	// var subscription Subscription
			// 	stream.SubscribeSync(
			// 		func(innerV R) {
			// 			subscriber.Next(innerV)
			// 		},
			// 		subscriber.Error,
			// 		func() {
			// 			activeSubscription--
			// 			for len(buffer) > 0 {
			// 				innerNext(buffer[0])
			// 				buffer = buffer[1:]
			// 			}
			// 			checkComplete()
			// 		},
			// 	)
			// }

			// source.SubscribeSync(
			// 	func(v T) {
			// 		if activeSubscription >= concurrent {
			// 			buffer = append(buffer, v)
			// 			return
			// 		}
			// 		innerNext(v)
			// 	},
			// 	subscriber.Error,
			// 	func() {
			// 		isComplete = true
			// 		checkComplete()
			// 	},
			// )
		})
	}
}

// Projects each source value to an Observable which is merged in the output
// Observable only if the previous projected Observable has completed.
func ExhaustMap[T any, R any](project func(value T, index uint) IObservable[R]) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			// var (
			// 	index        uint
			// 	isComplete   bool
			// 	subscription Subscription
			// )
			// source.SubscribeSync(
			// 	func(v T) {
			// 		if subscription == nil {
			// 			wg := new(sync.WaitGroup)
			// 			subscription = project(v, index).Subscribe(
			// 				func(v R) {
			// 					subscriber.Next(v)
			// 				},
			// 				func(error) {},
			// 				func() {
			// 					defer wg.Done()
			// 					subscription.Unsubscribe()
			// 					subscription = nil
			// 					if isComplete {
			// 						subscriber.Complete()
			// 					}
			// 				},
			// 			)
			// 			wg.Wait()
			// 		}
			// 		index++
			// 	},
			// 	subscriber.Error,
			// 	func() {
			// 		isComplete = true
			// 		if subscription == nil {
			// 			subscriber.Complete()
			// 		}
			// 	},
			// )

			// after collect the source
		})
	}
}

// // Merge the values from all observables to a single observable result.
// func ConcatAll[A any, B any](concurrent uint64) OperatorFunc[A, B] {
// 	return func(source IObservable[A]) IObservable[B] {
// 		return newObservable(func(subscriber Subscriber[B]) {
// 			source.SubscribeSync(func(a A) {}, func(err error) {}, func() {})
// 		})
// 	}
// }

// Merge the values from all observables to a single observable result.
func MergeAll[A any, B any](concurrent uint64) OperatorFunc[A, B] {
	return func(source IObservable[A]) IObservable[B] {
		return newObservable(func(subscriber Subscriber[B]) {

		})
	}
}

// Merge the values from all observables to a single observable result.
func Merge[T any](input IObservable[T]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				activeSubscription = 2
				wg                 = new(sync.WaitGroup)
				p1                 = source.SubscribeOn(wg.Done)
				p2                 = input.SubscribeOn(wg.Done)
				err                error
			)

			wg.Add(2)

			onNext := func(v Notification[T]) {
				if v == nil {
					return
				}

				// When any source errors, the resulting observable will error
				if err = v.Err(); err != nil {
					p1.Stop()
					p2.Stop()
					activeSubscription = 0
					return
				}

				if v.Done() {
					activeSubscription--
					return
				}

				subscriber.Send() <- v
			}

			for activeSubscription > 0 {
				select {
				case <-subscriber.Closed():
					return
				case v1 := <-p1.ForEach():
					onNext(v1)
				case v2 := <-p2.ForEach():
					onNext(v2)
				}
			}

			// Wait for all input streams to unsubscribe
			wg.Wait()

			if err != nil {
				subscriber.Send() <- ErrorNotification[T](err)
			} else {
				subscriber.Send() <- CompleteNotification[T]()
			}
		})
	}
}

// Useful for encapsulating and managing state. Applies an accumulator (or "reducer function")
// to each value from the source after an initial state is established --
// either via a seed value (second argument), or from the first value from the source.
func Scan[V any, A any](accumulator func(acc A, v V, index uint) (A, error), seed A) OperatorFunc[V, A] {
	return func(source IObservable[V]) IObservable[A] {
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

// Applies an accumulator function over the source Observable, and returns
// the accumulated result when the source completes, given an optional seed value.
func Reduce[V any, A any](accumulator func(acc A, v V, index uint) (A, error), seed A) OperatorFunc[V, A] {
	return func(source IObservable[V]) IObservable[A] {
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
				index++
			},
			func(obs Observer[A], err error) {
				obs.Error(err)
			},
			func(obs Observer[A]) {
				obs.Next(result)
				obs.Complete()
			},
		)
	}
}

// Delays the emission of items from the source Observable by a given timeout.
func Delay[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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
func Throttle[T any, R any](durationSelector func(v T) IObservable[R]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
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
	return func(source IObservable[T]) IObservable[T] {
		var (
			timer *time.Timer
		)
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

// Catches errors on the observable to be handled by returning a new observable or throwing an error.
func CatchError[T any](catch func(error, IObservable[T]) IObservable[T]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
				// subscription Subscription
				// subscribe func(IObservable[T])
			)

			// unsubscribe := func() {
			// 	if subscription != nil {
			// 		subscription.Unsubscribe()
			// 	}
			// 	subscription = nil
			// }

			// subscribe = func(stream IObservable[T]) {
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

			subscriber.Send() <- CompleteNotification[T]()
		})
	}
}

// Creates an Observable that mirrors the first source Observable to emit a
// next, error or complete notification from the combination of the Observable
// to which the operator is applied and supplied Observables.
func RaceWith[T any](input IObservable[T], inputs ...IObservable[T]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		inputs = append([]IObservable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				noOfInputs = len(inputs)
				// wg                  = new(sync.WaitGroup)
				fastestCh           = make(chan int, 1)
				activeSubscriptions = make([]Subscriber[T], noOfInputs)
				mu                  = new(sync.RWMutex)
				// unsubscribed        bool
			)
			// wg.Add(noOfInputs * 2)

			// unsubscribeAll := func(index int) {

			// 	var subscription Subscriber[T]

			// 	activeSubscriptions = []Subscriber[T]{subscription}
			// }

			// emit := func(index int, v Notification[T]) {
			// 	mu.RLock()
			// 	if unsubscribed {
			// 		mu.RUnlock()
			// 		return
			// 	}

			// 	log.Println("isThis", index)

			// 	mu.RUnlock()
			// 	mu.Lock()
			// 	unsubscribed = true
			// 	mu.Unlock()
			// 	// 	unsubscribeAll(index)

			// 	subscriber.Send() <- v
			// }

			for i, v := range inputs {
				activeSubscriptions[i] = v.SubscribeOn(func() {
					log.Println("DONE here")
					// wg.Done()
				})
				go func(idx int, obs Subscriber[T]) {
					// defer wg.Done()
					defer log.Println("closing routine", idx)

					for {
						select {
						case <-subscriber.Closed():
							log.Println("downstream closing ", idx)
							return
						case <-obs.Closed():
							log.Println("upstream closing ", idx)
							return
						case item := <-obs.ForEach():
							// mu.Lock()
							// defer mu.Unlock()
							// for _, sub := range activeSubscriptions {
							// 	sub.Stop()
							// }
							// activeSubscriptions = []Subscriber[T]{}
							log.Println("ForEach ah", idx, item)
							fastestCh <- idx
							// obs.Stop()
						}
					}
				}(i, activeSubscriptions[i])
			}

			log.Println("Fastest", <-fastestCh)
			mu.Lock()
			for _, v := range activeSubscriptions {
				v.Stop()
				log.Println(v)
			}
			mu.Unlock()
			// wg.Wait()

			log.Println("END")

			subscriber.Send() <- CompleteNotification[T]()
		})
	}
}

// Combines the source Observable with other Observables to create an Observable
// whose values are calculated from the latest values of each, only when the source emits.
func WithLatestFrom[A any, B any](input IObservable[B]) OperatorFunc[A, Tuple[A, B]] {
	return func(source IObservable[A]) IObservable[Tuple[A, B]] {
		return newObservable(func(subscriber Subscriber[Tuple[A, B]]) {
			// var (
			// 	allOk        [2]bool
			// 	latestA      A
			// 	latestB      B
			// 	subscription Subscription
			// )

			// // subscription = input.SubscribeOn(func(b B) {
			// // 	latestB = b
			// // 	allOk[1] = true
			// // }, func(err error) {}, func() {
			// // 	subscription.Unsubscribe()
			// // }, func() {})

			// source.SubscribeSync(
			// 	func(a A) {
			// 		latestA = a
			// 		allOk[0] = true
			// 		if allOk[0] && allOk[1] {
			// 			subscriber.Next(NewTuple(latestA, latestB))
			// 		}
			// 	},
			// 	subscriber.Error,
			// 	func() {
			// 		subscription.Unsubscribe()
			// 		subscriber.Complete()
			// 	},
			// )
		})
	}
}

// Collects all source emissions and emits them as an array when the source completes.
func ToArray[T any]() OperatorFunc[T, []T] {
	return func(source IObservable[T]) IObservable[[]T] {
		var (
			result = make([]T, 0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[[]T], v T) {
				result = append(result, v)
			},
			func(obs Observer[[]T], err error) {
				obs.Error(err)
			},
			func(obs Observer[[]T]) {
				obs.Next(result)
				obs.Complete()
			},
		)
	}
}
