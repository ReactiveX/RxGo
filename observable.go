package rxgo

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// An Observable that emits no items to the Observer and never completes.
func NEVER[T any]() IObservable[T] {
	return newObservable(func(sub Subscriber[T]) {})
}

// A simple Observable that emits no items to the Observer and immediately
// emits a complete notification.
func EMPTY[T any]() IObservable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		CompleteNotification[T]().Send(subscriber)
	})
}

// Creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
func Defer[T any](factory func() IObservable[T]) IObservable[T] {
	return factory()
}

func ThrownError[T any](factory func() error) IObservable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		ErrorNotification[T](factory()).Send(subscriber)
	})
}

// Creates an Observable that emits a sequence of numbers within a specified range.
func Range[T constraints.Unsigned](start, count T) IObservable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		var (
			end = start + count
		)

		for i := start; i < end; i++ {
			select {
			case <-subscriber.Closed():
				return
			case subscriber.Send() <- NextNotification(i):
			}
		}

		CompleteNotification[T]().Send(subscriber)
	})
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(duration time.Duration) IObservable[uint] {
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
				if NextNotification(index).Send(subscriber) {
					index++
				}
			}
		}
	})
}

func Scheduled[T any](item T, items ...T) IObservable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(subscriber Subscriber[T]) {
		for _, item := range items {
			notice := NextNotification(item)
			switch vi := any(item).(type) {
			case error:
				notice = ErrorNotification[T](vi)
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

		CompleteNotification[T]().Send(subscriber)
	})
}

func Timer[T any](start, interval time.Duration) IObservable[float64] {
	return newObservable(func(subscriber Subscriber[float64]) {
		var (
			latest = start
		)

		for {
			select {
			case <-subscriber.Closed():
				return
			case <-time.After(interval):
				subscriber.Send() <- NextNotification(latest.Seconds())
				latest = latest + interval
			}
		}
	})
}

func CombineLatest[A any, B any](first IObservable[A], second IObservable[B]) IObservable[Tuple[A, B]] {
	return newObservable(func(subscriber Subscriber[Tuple[A, B]]) {
		// var (
		// 	mu       sync.Mutex
		// 	latestA  A
		// 	latestB  B
		// 	hasValue = [2]bool{}
		// 	allOk    = [2]bool{}
		// )

		// nextValue := func() {
		// 	if hasValue[0] && hasValue[1] {
		// 		subscriber.Next(NewTuple(latestA, latestB))
		// 	}
		// }
		// checkComplete := func() {
		// 	if allOk[0] && allOk[1] {
		// 		subscriber.Complete()
		// 	}
		// }

		wg := new(sync.WaitGroup)
		wg.Add(2)
		// first.SubscribeOn(func(a A) {
		// 	mu.Lock()
		// 	defer mu.Unlock()
		// 	latestA = a
		// 	hasValue[0] = true
		// 	nextValue()
		// }, subscriber.Error, func() {
		// 	mu.Lock()
		// 	defer mu.Unlock()
		// 	allOk[0] = true
		// 	checkComplete()
		// }, wg.Done)
		// second.SubscribeOn(func(b B) {
		// 	mu.Lock()
		// 	defer mu.Unlock()
		// 	latestB = b
		// 	hasValue[1] = true
		// 	nextValue()
		// }, subscriber.Error, func() {
		// 	mu.Lock()
		// 	defer mu.Unlock()
		// 	allOk[1] = true
		// 	checkComplete()
		// }, wg.Done)
		wg.Wait()
	})
}
