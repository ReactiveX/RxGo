package rxgo

import (
	"time"
)

// Emits only the first count values emitted by the source Observable.
func Take[T any](count uint) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		if count == 0 {
			return EMPTY[T]()
		}

		seen := uint(0)
		return newObservable(func(subscriber Subscriber[T]) {
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

// Emits only the first value emitted by the source Observable that meets some condition.
func Find[T any](predicate func(T, uint) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		index := uint(0)
		return newObservable(func(subscriber Subscriber[T]) {
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
	var isEmpty = true
	return func(source IObservable[T]) IObservable[bool] {
		return newObservable(func(subscriber Subscriber[bool]) {
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
	var (
		lastValue T
		first     = true
	)
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
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
	var (
		lastValue T
		first     = true
	)
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
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

// Returns an Observable that emits whether or not every item of the
// source satisfies the condition specified.
func Every[T any](cb func(value T, count uint) bool) OperatorFunc[T, bool] {
	var (
		allOk = true
		index uint
	)
	return func(source IObservable[T]) IObservable[bool] {
		return newObservable(func(subscriber Subscriber[bool]) {
			source.SubscribeSync(
				func(t T) {
					allOk = allOk && cb(t, index)
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
func Repeat[T any, N Number](count N) OperatorFunc[T, T] {
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
	hasValue := false
	return func(source IObservable[T]) IObservable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
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
				Take[T](1),
				DefaultIfEmpty(defaultValue[0]),
			)
		}
	}

	return func(source IObservable[T]) IObservable[T] {
		return Pipe2(source,
			Filter(func(_ T, i uint) bool {
				return i == index
			}),
			Take[T](1),
		)
	}
}

// Emits only the first value (or the first value that meets some condition)
// emitted by the source Observable.
func First[T any]() OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		return Pipe1(source, Take[T](1))
	}
}

// Returns a result Observable that emits all values pushed by the source observable
// if they are distinct in comparison to the last value the result observable emitted.
func DistinctUntilChanged[T any](comparator func(prev T, current T) bool) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		var (
			lastValue T
			first     = true
		)
		return newObservable(func(subscriber Subscriber[T]) {
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

func ExhaustMap[T any, R any](project func(value T, index uint) R) OperatorFunc[T, T] {
	return func(source IObservable[T]) IObservable[T] {
		index := uint(0)
		return newObservable(func(subscriber Subscriber[T]) {
			source.SubscribeSync(
				func(v T) {
					// if filter(v, index) {
					// 	subscriber.Next(v)
					// }
					index++
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
		index := uint(0)
		return newObservable(func(subscriber Subscriber[T]) {
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
		index := uint(0)
		return newObservable(func(subscriber Subscriber[R]) {
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
