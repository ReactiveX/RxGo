package rxgo

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

func skip[T any](v T) {}

// Emits only the first count values emitted by the source Observable.
func Take[T any, N constraints.Unsigned](count N) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		if count == 0 {
			return EMPTY[T]()
		}

		return newObservable(func(subscriber Subscriber[T]) {
			seen := N(0)
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

// Emits the values emitted by the source Observable until a notifier Observable emits a value.
func TakeUntil[T any, R any](notifier IObservable[R]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var isComplete bool

			subscription := source.Subscribe(
				func(v T) {
					subscriber.Next(v)
				},
				subscriber.Error,
				func() {
					isComplete = true
					subscriber.Complete()
				},
			)

			notifier.SubscribeSync(func(v R) {
				if !isComplete {
					subscription.Unsubscribe()
				}
			}, func(err error) {
				subscription.Unsubscribe()
			}, func() {
				subscription.Unsubscribe()
			})
		})
	}
}

// Emits values emitted by the source Observable so long as each value satisfies the given predicate,
// and then completes as soon as this predicate is not satisfied.
func TakeWhile[T any](predicate func(value T, index uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var index uint
			source.SubscribeSync(
				func(v T) {
					if predicate(v, index) {
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

// Waits for the source to complete, then emits the last N values from the source,
// as specified by the count argument.
func TakeLast[T any, N constraints.Unsigned](count N) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			values := make([]T, count)
			source.SubscribeSync(
				func(v T) {
					if N(len(values)) >= count {
						// shift the item from queue
						values = values[1:]
					}
					values = append(values, v)
				},
				subscriber.Error,
				func() {
					for _, v := range values {
						subscriber.Next(v)
					}
					subscriber.Complete()
				},
			)
		})
	}
}

// Returns an Observable that skips the first count items emitted by the source Observable.
func Skip[T any](count uint) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				index uint
			)
			source.SubscribeSync(
				func(v T) {
					index++
					if count >= index {
						return
					}
					subscriber.Next(v)
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Emits the single value at the specified index in a sequence of emissions from the source Observable.
func ElementAt[T any](index uint, defaultValue ...T) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				i       uint
				emitted bool
			)
			source.SubscribeSync(
				func(v T) {
					if index == i {
						subscriber.Next(v)
						subscriber.Complete()
						emitted = true
						return
					}
					i++
				},
				subscriber.Error,
				func() {
					defer subscriber.Complete()
					if emitted {
						return
					}

					if len(defaultValue) > 0 {
						subscriber.Next(defaultValue[0])
						return
					}

					subscriber.Error(ErrArgumentOutOfRange)
				},
			)
		})
	}
}

// Emits only the first value (or the first value that meets some condition)
// emitted by the source Observable.
func First[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return Pipe1(source, Take[T, uint](1))
	}
}

// Returns an Observable that emits only the last item emitted by the source Observable.
// It optionally takes a predicate function as a parameter, in which case,
// rather than emitting the last item from the source Observable,
// the resulting Observable will emit the last item from the source Observable
// that satisfies the predicate.
func Last[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return Pipe1(source, TakeLast[T, uint](1))
	}
}

// Emits only the first value emitted by the source Observable that meets some condition.
func Find[T any](predicate func(T, uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var index uint
			source.SubscribeSync(
				func(v T) {
					if predicate(v, index) {
						subscriber.Next(v)
						subscriber.Complete()
					}
					index++
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Emits only the index of the first value emitted by the source Observable that meets some condition.
func FindIndex[T any](predicate func(T, uint) bool) OperatorFunc[T, int] {
	return func(source IObservable[T]) IObservable[int] {
		return newObservable(func(subscriber Subscriber[int]) {
			var (
				index uint
				ok    bool
			)
			source.SubscribeSync(
				func(v T) {
					if predicate(v, index) {
						ok = true
						subscriber.Next(int(index))
						subscriber.Complete()
					}
					index++
				},
				subscriber.Error,
				func() {
					if !ok {
						subscriber.Next(-1)
					}
					subscriber.Complete()
				},
			)
		})
	}
}

// The Min operator operates on an Observable that emits numbers
// (or items that can be compared with a provided function),
// and when source Observable completes it emits a single item: the item with the smallest value.
func Min[T any](comparer func(a T, b T) int8) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				lastValue T
				first     = true
			)
			source.SubscribeSync(
				func(v T) {
					if first {
						lastValue = v
						first = false
						return
					}

					switch comparer(lastValue, v) {
					case 1:
						fallthrough
					default:
						lastValue = v
					}
				},
				subscriber.Error,
				func() {
					subscriber.Next(lastValue)
					subscriber.Complete()
				},
			)
		})
	}
}

// The Max operator operates on an Observable that emits numbers
// (or items that can be compared with a provided function),
// and when source Observable completes it emits a single item: the item with the largest value.
func Max[T any](comparer func(a T, b T) int8) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				lastValue T
				first     = true
			)
			source.SubscribeSync(
				func(v T) {
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
				subscriber.Error,
				func() {
					subscriber.Next(lastValue)
					subscriber.Complete()
				},
			)
		})
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
		return newObservable(func(subscriber Subscriber[uint]) {
			var (
				count uint
				index uint
			)
			source.SubscribeSync(
				func(v T) {
					if cb(v, index) {
						count++
					}
					index++
				},
				subscriber.Error,
				func() {
					subscriber.Next(count)
					subscriber.Complete()
				},
			)
		})
	}
}

// Ignores all items emitted by the source Observable and only passes calls of complete or error.
func IgnoreElements[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				skip[T],
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Returns an Observable that emits whether or not every item of the
// source satisfies the condition specified.
func Every[T any](predicate func(value T, count uint) bool) OperatorFunc[T, bool] {
	return func(source IObservable[T]) IObservable[bool] {
		return newObservable(func(subscriber Subscriber[bool]) {
			var (
				allOk = true
				index uint
			)
			source.SubscribeSync(
				func(t T) {
					allOk = allOk && predicate(t, index)
				},
				subscriber.Error,
				func() {
					subscriber.Next(allOk)
					subscriber.Complete()
				},
			)
		})
	}
}

// Returns an Observable that will resubscribe to the source
// stream when the source stream completes.
func Repeat[T any, N constraints.Unsigned](count N) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(t T) {
					for i := N(0); i < count; i++ {
						subscriber.Next(t)
					}
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Emits false if the input Observable emits any values,
// or emits true if the input Observable completes without emitting any values.
func IsEmpty[T any]() OperatorFunc[T, bool] {
	return func(source IObservable[T]) IObservable[bool] {
		return newObservable(func(subscriber Subscriber[bool]) {
			var isEmpty = true
			source.SubscribeSync(
				func(t T) {
					isEmpty = false
				},
				subscriber.Error,
				func() {
					subscriber.Next(isEmpty)
					subscriber.Complete()
				},
			)
		})
	}
}

// Emits a given value if the source Observable completes without emitting any
// next value, otherwise mirrors the source Observable.
func DefaultIfEmpty[T any](defaultValue T) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			hasValue := false
			source.SubscribeSync(
				func(t T) {
					hasValue = true
					subscriber.Next(t)
				},
				subscriber.Error,
				func() {
					if !hasValue {
						subscriber.Next(defaultValue)
					}
					subscriber.Complete()
				},
			)
		})
	}
}

// Returns a result Observable that emits all values pushed by the source observable
// if they are distinct in comparison to the last value the result observable emitted.
func DistinctUntilChanged[T any](comparator func(prev T, current T) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				lastValue T
				first     = true
			)
			source.SubscribeSync(
				func(v T) {
					if first || !comparator(lastValue, v) {
						subscriber.Next(v)
						first = false
					}
					lastValue = v
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Filter emits only those items from an Observable that pass a predicate test.
func Filter[T any](filter func(T, uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var index uint
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

// Map transforms the items emitted by an Observable by applying a function to each item.
func Map[T any, R any](mapper func(T, uint) (R, error)) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var index uint
			source.SubscribeSync(
				func(v T) {
					output, err := mapper(v, index)
					index++
					if err != nil {
						subscriber.Error(err)
						return
					}
					subscriber.Next(output)
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Returns an observable that asserts that only one value is emitted from the observable
// that matches the predicate. If no predicate is provided, then it will assert that the
// observable only emits one value.
// FIXME: should rename `Single2` to `Single`
func Single2[T any](predicate func(v T, index uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				index    uint
				found    bool
				matches  uint
				hasValue bool
			)
			source.SubscribeSync(
				func(v T) {
					result := predicate(v, index)
					if result {
						found = result
						matches++
					}
					hasValue = true
					index++
				},
				subscriber.Error,
				func() {
					if !hasValue {
						subscriber.Error(ErrEmpty)
					} else if !found {
						subscriber.Error(ErrNotFound)
					} else if matches > 1 {
						subscriber.Error(ErrSequence)
					}
					subscriber.Complete()
				},
			)
		})
	}
}

// Returns an Observable that skips all items emitted by the source Observable
// as long as a specified condition holds true, but emits all further source items
// as soon as the condition becomes false.
func SkipWhile[T any](predicate func(v T, index uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				index uint
				skip  = true
			)
			source.SubscribeSync(
				func(v T) {
					if predicate(v, index) {
						skip = false
					}
					if !skip {
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

// Projects each source value to an Observable which is merged in the output Observable,
// in a serialized fashion waiting for each one to complete before merging the next.
func ConcatMap[T any, R any](project func(value T, index uint) IObservable[R]) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				index              uint
				buffer             = make([]T, 0)
				concurrent         = uint(1)
				isComplete         bool
				activeSubscription uint
			)

			checkComplete := func() {
				if isComplete && len(buffer) <= 0 {
					subscriber.Complete()
				}
			}

			var innerNext func(T)
			innerNext = func(outerV T) {
				activeSubscription++

				stream := project(outerV, index)
				index++

				// var subscription Subscription
				stream.SubscribeSync(
					func(innerV R) {
						subscriber.Next(innerV)
					},
					subscriber.Error,
					func() {
						activeSubscription--
						for len(buffer) > 0 {
							innerNext(buffer[0])
							buffer = buffer[1:]
						}
						checkComplete()
					},
				)
			}

			source.SubscribeSync(
				func(v T) {
					if activeSubscription >= concurrent {
						buffer = append(buffer, v)
						return
					}
					innerNext(v)
				},
				subscriber.Error,
				func() {
					isComplete = true
					checkComplete()
				},
			)
		})
	}
}

// Projects each source value to an Observable which is merged in the output
// Observable only if the previous projected Observable has completed.
func ExhaustMap[T any, R any](project func(value T, index uint) IObservable[R]) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				index        uint
				isComplete   bool
				subscription Subscription
			)
			source.SubscribeSync(
				func(v T) {
					if subscription == nil {
						wg := new(sync.WaitGroup)
						subscription = project(v, index).Subscribe(
							func(v R) {
								subscriber.Next(v)
							},
							func(error) {},
							func() {
								defer wg.Done()
								subscription.Unsubscribe()
								subscription = nil
								if isComplete {
									subscriber.Complete()
								}
							},
						)
						wg.Wait()
					}
					index++
				},
				subscriber.Error,
				func() {
					isComplete = true
					if subscription == nil {
						subscriber.Complete()
					}
				},
			)

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

// // Merge the values from all observables to a single observable result.
// func MergeWith1[A any, B any](IObservable[B]) OperatorFunc[A, B] {
// 	return func(source IObservable[A]) IObservable[B] {
// 		return newObservable(func(subscriber Subscriber[B]) {

// 		})
// 	}
// }

// Useful for encapsulating and managing state. Applies an accumulator (or "reducer function")
// to each value from the source after an initial state is established --
// either via a seed value (second argument), or from the first value from the source.
func Scan[V any, A any](accumulator func(acc A, v V, index uint) A, seed A) OperatorFunc[V, A] {
	return func(source IObservable[V]) IObservable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			var (
				index uint
			)
			source.SubscribeSync(
				func(v V) {
					seed = accumulator(seed, v, index)
					subscriber.Next(seed)
					index++
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Applies an accumulator function over the source Observable, and returns
// the accumulated result when the source completes, given an optional seed value.
func Reduce[V any, A any](accumulator func(acc A, v V, index uint) A, seed A) OperatorFunc[V, A] {
	return func(source IObservable[V]) IObservable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			var (
				index uint
			)
			source.SubscribeSync(
				func(v V) {
					seed = accumulator(seed, v, index)
					index++
				},
				subscriber.Error,
				func() {
					subscriber.Next(seed)
					subscriber.Complete()
				},
			)
		})
	}
}

// Delays the emission of items from the source Observable by a given timeout.
func Delay[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(v T) {
					time.Sleep(duration)
					subscriber.Next(v)
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

// Emits a value from the source Observable, then ignores subsequent source values
// for duration milliseconds, then repeats this process.
func Throttle[T any, R any](durationSelector func(v T) IObservable[R]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(v T) {
					durationSelector(v)
					subscriber.Next(v)
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

func DebounceTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		// var (
		// 	lastTime  time.Time
		// 	lastValue T
		// )
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(v T) {
					// lastValue = v
					// lastTime = time.Now()
					subscriber.Next(v)
				},
				subscriber.Error,
				subscriber.Complete,
			)
		})
	}
}

func CatchError[T any](catch func(error) IObservable[T]) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				subscriber.Next,
				func(err error) {
					catch(err)
				},
				subscriber.Complete,
			)
		})
	}
}

// Collects all source emissions and emits them as an array when the source completes.
func ToArray[T any]() OperatorFunc[T, []T] {
	return func(source IObservable[T]) IObservable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			result := make([]T, 0)
			source.SubscribeSync(
				func(v T) {
					result = append(result, v)
				},
				subscriber.Error,
				func() {
					subscriber.Next(result)
					subscriber.Complete()
				},
			)
		})
	}
}
