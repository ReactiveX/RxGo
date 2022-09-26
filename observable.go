package rxgo

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// An Observable that emits no items to the Observer and never completes.
func Never[T any]() Observable[T] {
	return newObservable(func(sub Subscriber[T]) {})
}

// A simple Observable that emits no items to the Observer and immediately emits a complete notification.
func Empty[T any]() Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		Complete[T]().Send(subscriber)
	})
}

// Creates an Observable that, on subscribe, calls an Observable factory to make an Observable for each new Observer.
func Defer[T any](factory func() Observable[T]) Observable[T] {
	// `Defer` allows you to create an Observable only when the Observer subscribes. It waits until an Observer subscribes to it, calls the given factory function to get an Observable -- where a factory function typically generates a new Observable -- and subscribes the Observer to this Observable. In case the factory function returns a falsy value, then Empty is used as Observable instead. Last but not least, an exception during the factory function call is transferred to the Observer by calling error.
	return newObservable(func(subscriber Subscriber[T]) {
		var (
			wg     = new(sync.WaitGroup)
			stream = factory()
		)

		if stream == nil {
			stream = Empty[T]()
		}

		wg.Add(1)

		var (
			upStream = stream.SubscribeOn(wg.Done)
			ended    bool
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

				ended = item.Done() || item.Err() != nil
				item.Send(subscriber)
				if ended {
					break loop
				}
			}
		}

		wg.Wait()
	})
}

// Creates an Observable that emits a sequence of numbers within a specified range.
func Range[T constraints.Unsigned](start, count T) Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		var (
			end    = start + count
			closed bool
		)

	loop:
		for i := start; i < end; i++ {
			select {
			case <-subscriber.Closed():
				closed = true
				break loop

			case subscriber.Send() <- Next(i):
			}
		}

		if !closed {
			Complete[T]().Send(subscriber)
		}
	})
}

// Interval creates an Observable emitting incremental integers infinitely between each given time interval.
func Interval(duration time.Duration) Observable[uint] {
	return newObservable(func(subscriber Subscriber[uint]) {
		var (
			index uint
		)

	loop:
		for {
			select {
			// If receiver notify stop, we should terminate the operation
			case <-subscriber.Closed():
				break loop
			case <-time.After(duration):
				if Next(index).Send(subscriber) {
					index++
				}
			}
		}
	})
}

// FIXME: rename me to `Of`
func Of2[T any](item T, items ...T) Observable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(subscriber Subscriber[T]) {
		for _, item := range items {
			select {
			// If receiver notify stop, we should terminate the operation
			case <-subscriber.Closed():
				return
			case subscriber.Send() <- Next(item):
			}
		}

		Complete[T]().Send(subscriber)
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

// Creates an observable that will create an error instance and push it to the consumer as an error immediately upon subscription. This creation function is useful for creating an observable that will create an error and error every time it is subscribed to. Generally, inside of most operators when you might want to return an errored observable, this is unnecessary. In most cases, such as in the inner return of `ConcatMap`, `MergeMap`, `Defer`, and many others, you can simply throw the error, and RxGo will pick that up and notify the consumer of the error.
func Throw[T any](factory ErrorFunc) Observable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		Error[T](factory()).Send(subscriber)
	})
}

// Creates an observable that will wait for a specified time period before emitting the number 0.
func Timer[N constraints.Unsigned](startDue time.Duration, intervalDuration ...time.Duration) Observable[N] {
	return newObservable(func(subscriber Subscriber[N]) {
		var (
			index = N(0)
		)

		time.Sleep(startDue)
		Next(index).Send(subscriber)
		index++

		if len(intervalDuration) > 0 {
			startDue = intervalDuration[0]
			timeout := time.After(startDue)

			for {
				select {
				case <-subscriber.Closed():
					return
				case <-timeout:
					Next(index).Send(subscriber)
					index++
					timeout = time.After(startDue)
				}
			}
		}

		Complete[N]().Send(subscriber)
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

// Splits the source Observable into two, one with values that satisfy a predicate, and another with values that don't satisfy the predicate.
// FIXME: redesign the API
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
