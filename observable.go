package rxgo

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// An Observable that emits no items to the Observer and never completes.
func NEVER[T any]() Observable[T] {
	return newObservable(func(sub Subscriber[T]) {})
}

// A simple Observable that emits no items to the Observer and immediately
// emits a complete notification.
func EMPTY[T any]() Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		Complete[T]().Send(subscriber)
	})
}

// Creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
func Defer[T any](factory func() Observable[T]) Observable[T] {
	// defer allows you to create an Observable only when the Observer subscribes.
	// It waits until an Observer subscribes to it, calls the given factory function
	// to get an Observable -- where a factory function typically generates a new
	// Observable -- and subscribes the Observer to this Observable. In case the factory
	// function returns a falsy value, then EMPTY is used as Observable instead.
	// Last but not least, an exception during the factory function call is transferred
	// to the Observer by calling error.
	obs := factory()
	if obs == nil {
		return EMPTY[T]()
	}
	return obs
}

// Creates an Observable that emits a sequence of numbers within a specified range.
func Range[T constraints.Unsigned](start, count T) Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		var (
			end = start + count
		)

		for i := start; i < end; i++ {
			select {
			case <-subscriber.Closed():
				return
			case subscriber.Send() <- Next(i):
			}
		}

		Complete[T]().Send(subscriber)
	})
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(duration time.Duration) Observable[uint] {
	return newObservable(func(subscriber Subscriber[uint]) {
		var (
			index uint
		)

		for {
			select {
			// If receiver notify stop, we should terminate the operation
			case <-subscriber.Closed():
				return
			case <-time.After(duration):
				if Next(index).Send(subscriber) {
					index++
				}
			}
		}
	})
}

func Scheduled[T any](item T, items ...T) Observable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(subscriber Subscriber[T]) {
		for _, item := range items {
			notice := Next(item)
			switch vi := any(item).(type) {
			case error:
				notice = Error[T](vi)
			}

			select {
			// If receiver notify stop, we should terminate the operation
			case <-subscriber.Closed():
				return
			case subscriber.Send() <- notice:
			}

			if err := notice.Err(); err != nil {
				return
			}
		}

		Complete[T]().Send(subscriber)
	})
}

func ThrownError[T any](factory func() error) Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		Error[T](factory()).Send(subscriber)
	})
}

func Timer[T any](start, interval time.Duration) Observable[float64] {
	return newObservable(func(subscriber Subscriber[float64]) {
		var (
			latest = start
		)

		for {
			select {
			case <-subscriber.Closed():
				return
			case <-time.After(interval):
				subscriber.Send() <- Next(latest.Seconds())
				latest = latest + interval
			}
		}
	})
}

// Checks a boolean at subscription time, and chooses between one of two observable sources
func Iif[T any](condition func() bool, trueObservable Observable[T], falseObservable Observable[T]) Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		var (
			wg         = new(sync.WaitGroup)
			observable = falseObservable
		)

		if condition() {
			observable = trueObservable
		}

		wg.Add(1)

		var (
			upStream = observable.SubscribeOn(wg.Done)
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
				item.Send(subscriber)
				if ended {
					break loop
				}
			}
		}

		wg.Wait()
	})
}

// Splits the source Observable into two, one with values that satisfy a predicate,
// and another with values that don't satisfy the predicate.
func Partition[T any](source Observable[T], predicate PredicateFunc[T]) {
	newObservable(func(subscriber Subscriber[Tuple[Observable[T], Observable[T]]]) {
		var (
			wg          = new(sync.WaitGroup)
			trueStream  = NewSubscriber[T]()
			falseStream = NewSubscriber[T]()
		)

		wg.Add(1)

		var (
			index    uint
			upStream = source.SubscribeOn(wg.Done)
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
				if predicate(item.Value(), index) {
					item.Send(trueStream)
				} else {
					item.Send(falseStream)
				}

				if ended {
					break loop
				}

				index++
			}
		}

		wg.Wait()
		// Next(NewTuple(trueStream, falseStream)).Send(subscriber)
	})
}
