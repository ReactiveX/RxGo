package rxgo

import (
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
)

var (
	ErrEmpty    = fmt.Errorf("rxgo: empty value")
	ErrNotFound = fmt.Errorf("rxgo: no values match")
	ErrSequence = fmt.Errorf("rxgo: too many values match")
)

func noop[T any](v T) {}

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

// Ignores all items emitted by the source Observable and only passes calls of complete or error.
func IgnoreElements[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				noop[T],
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

// Emits the single value at the specified index in a sequence of emissions
// from the source Observable.
func ElementAt[T any](index uint, defaultValue ...T) OperatorFunc[T, T] {
	if len(defaultValue) > 0 {
		return func(source IObservable[T]) IObservable[T] {
			return Pipe3(source,
				Filter(func(_ T, i uint) bool {
					return i == index
				}),
				Take[T, uint](1),
				DefaultIfEmpty(defaultValue[0]),
			)
		}
	}

	return func(source IObservable[T]) IObservable[T] {
		return Pipe2(source,
			Filter(func(_ T, i uint) bool {
				return i == index
			}),
			Take[T, uint](1),
		)
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

func ExhaustMap[T any, R any](project func(value T, index uint) IObservable[R]) OperatorFunc[T, R] {
	return func(source IObservable[T]) IObservable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var index uint
			source.SubscribeSync(
				func(v T) {
					project(v, index).SubscribeSync(
						subscriber.Next,
						subscriber.Error,
						subscriber.Complete,
					)
					index++
				},
				subscriber.Error,
				subscriber.Complete,
			)

			// after collect the source
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
