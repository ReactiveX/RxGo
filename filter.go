package rxgo

import (
	"reflect"
	"sync"
	"time"
)

// Ignores source values for a duration determined by another Observable, then emits the most recent value from the source Observable, then repeats this process.
func Audit[T any, R any](durationSelector DurationFunc[T, R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream       = source.SubscribeOn(wg.Done)
				durationStream Subscriber[R]
				durationCh     <-chan Notification[R]
				latestValue    T
			)

			setValues := func() {
				durationStream = nil
				durationCh = make(<-chan Notification[R])
			}

			unsubscribeStream := func() {
				if durationStream != nil {
					durationStream.Stop()
				}
				setValues()
			}

			setValues()

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
						Error[T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						Complete[T]().Send(subscriber)
						break observe
					}

					latestValue = item.Value()
					if durationStream == nil {
						wg.Add(1)
						durationStream = durationSelector(latestValue).SubscribeOn(wg.Done)
						durationCh = durationStream.ForEach()
					}

				case item, ok := <-durationCh:
					if !ok {
						continue
					}

					// TODO: handle done?

					if err := item.Err(); err != nil {
						upStream.Stop()
						Error[T](err).Send(subscriber)
						break observe
					}

					Next(latestValue).Send(subscriber)

					// reset
					unsubscribeStream()
				}
			}

			// prevent leaking
			unsubscribeStream()

			wg.Wait()
		})
	}
}

// Ignores source values for duration milliseconds, then emits the most recent value from the source Observable, then repeats this process.
func AuditTime[T any, R any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return Pipe1(
			source,
			Debounce(func(value T) Observable[uint] {
				// FIXME: maybe replace it to timer
				return Interval(duration)
			}),
		)
	}
}

// Emits a notification from the source Observable only after a particular time span determined by another Observable has passed without another source emission.
func Debounce[T any, R any](durationSelector DurationFunc[T, R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream    = source.SubscribeOn(wg.Done)
				downStream  Subscriber[R]
				notifyCh    <-chan Notification[R]
				latestValue T
			)

			setValues := func() {
				downStream = nil
				notifyCh = make(<-chan Notification[R])
			}

			unsubscribeStream := func() {
				if downStream != nil {
					downStream.Stop()
				}
				setValues()
			}

			setValues()

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
						Error[T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						Complete[T]().Send(subscriber)
						break observe
					}

					// The notification is emitted only when the duration Observable emits a next notification, and if no other notification was emitted on the source Observable since the duration Observable was spawned. If a new notification appears before the duration Observable emits, the previous notification will not be emitted and a new duration is scheduled from durationSelector is scheduled.
					latestValue = item.Value()
					unsubscribeStream()
					if downStream == nil {
						wg.Add(1)
						downStream = durationSelector(latestValue).SubscribeOn(wg.Done)
						notifyCh = downStream.ForEach()
					}

				// TODO: goroutine selection is chosen via a uniform pseudo-random selection: https://go.dev/ref/spec#Select_statements
				case item, ok := <-notifyCh:
					if !ok {
						continue
					}

					// TODO: handle done?

					if err := item.Err(); err != nil {
						upStream.Stop()
						Error[T](err).Send(subscriber)
						break observe
					}

					Next(latestValue).Send(subscriber)

					// reset
					unsubscribeStream()
				}
			}

			// prevent leaking
			unsubscribeStream()

			wg.Wait()
		})
	}
}

// Emits a notification from the source Observable only after a particular time span
// has passed without another source emission.
func DebounceTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return Pipe1(
			source,
			Debounce(func(value T) Observable[uint] {
				// FIXME: maybe replace it to timer
				return Interval(duration)
			}),
		)
	}
}

// Returns an Observable that emits all items emitted by the source Observable
// that are distinct by comparison from previous items.
func Distinct[T any, K comparable](keySelector func(value T) K) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			keySet = make(map[K]bool)
			exists bool
			key    K
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				key = keySelector(v)
				_, exists = keySet[key]
				if !exists {
					keySet[key] = true
					obs.Next(v)
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

// Returns a result Observable that emits all values pushed by the source observable
// if they are distinct in comparison to the last value the result observable emitted.
func DistinctUntilChanged[T any](comparator ...ComparatorFunc[T, T]) OperatorFunc[T, T] {
	cb := func(prev T, current T) bool {
		return reflect.DeepEqual(prev, current)
	}
	if len(comparator) > 0 {
		cb = comparator[0]
	}
	return func(source Observable[T]) Observable[T] {
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

// Filter emits only those items from an Observable that pass a predicate test.
func Filter[T any](predicate PredicateFunc[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
		)
		cb := skipPredicate[T]
		if predicate != nil {
			cb = predicate
		}
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if cb(v, index) {
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

// Emits the most recently emitted value from the source Observable whenever
// another Observable, the notifier, emits.
func Sample[T any, R any](notifier Observable[R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				latestValue  Notification[T]
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = notifier.SubscribeOn(wg.Done)
			)

			unsubscribeAll := func() {
				upStream.Stop()
				notifyStream.Stop()
			}

		observe:
			for {
				select {
				case <-subscriber.Closed():
					unsubscribeAll()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						unsubscribeAll()
						break observe
					}

					if err := item.Err(); err != nil {
						item.Send(subscriber)
						unsubscribeAll()
						break observe
					}

					if item.Done() {
						item.Send(subscriber)
						unsubscribeAll()
						break observe
					}
					latestValue = item

				case <-notifyStream.ForEach():
					if latestValue != nil {
						latestValue.Send(subscriber)
					}
				}
			}

			wg.Wait()
		})
	}
}

// Emits the most recently emitted value from the source Observable within periodic time
// intervals.
func SampleTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return Pipe1(source, Sample[T](Interval(duration)))
	}
}

// Returns an observable that asserts that only one value is emitted from the observable
// that matches the predicate. If no predicate is provided, then it will assert that the
// observable only emits one value.
func Single[T any](predicate ...func(value T, index uint, source Observable[T]) bool) OperatorFunc[T, T] {
	cb := func(T, uint, Observable[T]) bool {
		return true
	}
	if len(predicate) > 0 {
		cb = predicate[0]
	}
	return func(source Observable[T]) Observable[T] {
		var (
			index    uint
			hasValue bool
			result   = make([]T, 0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				hasValue = true
				if cb(v, index, source) {
					result = append(result, v)
				}
				index++
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				noOfResult := len(result)
				if !hasValue {
					obs.Error(ErrEmpty)
					return
				} else if noOfResult > 1 {
					obs.Error(ErrSequence)
					return
				} else if noOfResult < 1 {
					obs.Error(ErrNotFound)
					return
				}
				obs.Next(result[0])
				obs.Complete()
			},
		)
	}
}

// Returns an Observable that skips the first count items emitted by the source Observable.
func Skip[T any](count uint) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				index++
				if count >= index {
					return
				}
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

// Skip a specified number of values before the completion of an observable.
func SkipLast[T any](skipCount uint) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			values = make([]T, 0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				values = append(values, v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				values = values[:uint(len(values))-skipCount]
				for _, v := range values {
					obs.Next(v)
				}
				obs.Complete()
			},
		)
	}
}

// Returns an Observable that skips items emitted by the source Observable until a
// second Observable emits an item.
func SkipUntil[T any, R any](notifier Observable[R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				skip         = true
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = notifier.SubscribeOn(wg.Done)
			)

			// It will never let the source observable emit any values if the
			// notifier completes or throws an error without emitting a value before.

		loop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					notifyStream.Stop()
					break loop

				case item, ok := <-upStream.ForEach():
					if !ok {
						notifyStream.Stop()
						break loop
					}

					ended := item.Err() != nil || item.Done()
					if ended {
						notifyStream.Stop()
						item.Send(subscriber)
						break loop
					}

					if !skip {
						item.Send(subscriber)
					}

				// Internally the skipUntil operator subscribes to the passed in observable
				// (in the following called notifier) in order to recognize the emission of
				// its first value. When this happens, the operator unsubscribes from the
				// notifier and starts emitting the values of the source observable.
				case <-notifyStream.ForEach():
					notifyStream.Stop()
					skip = false
				}
			}

			wg.Wait()
		})
	}
}

// Returns an Observable that skips all items emitted by the source Observable
// as long as a specified condition holds true, but emits all further source items
// as soon as the condition becomes false.
func SkipWhile[T any](predicate func(v T, index uint) bool) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
			pass  bool
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if pass {
					obs.Next(v)
					return
				}
				if !predicate(v, index) {
					pass = true
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

// Emits only the first count values emitted by the source Observable.
func Take[T any](count uint) OperatorFunc[T, T] {
	if count == 0 {
		return func(source Observable[T]) Observable[T] {
			return Empty[T]()
		}
	}

	return func(source Observable[T]) Observable[T] {
		var (
			seen = uint(0)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				seen++
				if seen <= count {
					obs.Next(v)
					if count <= seen {
						obs.Complete()
					}
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

// Waits for the source to complete, then emits the last N values from the source,
// as specified by the count argument.
func TakeLast[T any](count uint) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			values = make([]T, count)
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if uint(len(values)) >= count {
					// shift the item from queue
					values = values[1:]
				}
				values = append(values, v)
			},
			func(obs Observer[T], err error) {
				obs.Error(err)
			},
			func(obs Observer[T]) {
				for _, v := range values {
					obs.Next(v)
				}
				obs.Complete()
			},
		)
	}
}

// Emits the values emitted by the source Observable until a notifier Observable emits a value.
func TakeUntil[T any, R any](notifier Observable[R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = notifier.SubscribeOn(wg.Done)
			)

		loop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					notifyStream.Stop()
					break loop

				case item, ok := <-upStream.ForEach():
					if !ok {
						notifyStream.Stop()
						break loop
					}

					ended := item.Err() != nil || item.Done()
					item.Send(subscriber)
					if ended {
						notifyStream.Stop()
						break loop
					}

				// Lets values pass until notifier Observable emits a value.
				// Then, it completes.
				case <-notifyStream.ForEach():
					upStream.Stop()
					notifyStream.Stop()
					Complete[T]().Send(subscriber)
					break loop
				}
			}

			wg.Wait()
		})
	}
}

// Emits values emitted by the source Observable so long as each value satisfies the given predicate,
// and then completes as soon as this predicate is not satisfied.
func TakeWhile[T any](predicate func(value T, index uint) bool) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[T], v T) {
				if !predicate(v, index) {
					obs.Complete()
					return
				}
				obs.Next(v)
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

// Emits a value from the source Observable, then ignores subsequent source values
// for a duration determined by another Observable, then repeats this process.
func Throttle[T any, R any](durationSelector func(value T) Observable[R]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
				canEmit  = true
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

					ended := item.Err() != nil || item.Done()
					if ended {
						item.Send(subscriber)
						break loop
					}

					if canEmit {
						item.Send(subscriber)
						canEmit = false
					}

					wg.Add(1)
					durationSelector(item.Value()).SubscribeOn(wg.Done)
				}
			}

			wg.Wait()
		})
	}
}

// Emits a value from the source Observable, then ignores subsequent source
// values for duration milliseconds, then repeats this process
func ThrottleTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
				canEmit  = true
				timeout  = time.After(duration)
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

					ended := item.Err() != nil || item.Done()
					if ended {
						item.Send(subscriber)
						break loop
					}
					if canEmit {
						item.Send(subscriber)
						canEmit = false
					}

				case <-timeout:
					canEmit = true
					timeout = time.After(duration)
				}
			}

			wg.Wait()
		})
	}
}
