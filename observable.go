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
		subscriber.Send() <- newComplete[T]()
	})
}

// Creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
func Defer[T any](factory func() IObservable[T]) IObservable[T] {
	return factory()
}

func ThrownError[T any](factory func() error) IObservable[T] {
	return newObservable(func(subscriber Subscriber[T]) {
		select {
		case <-subscriber.Closed():
			return
		case subscriber.Send() <- newError[T](factory()):
		}
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
			case subscriber.Send() <- newData(i):
			}
		}

		subscriber.Send() <- newComplete[T]()
	})
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(duration time.Duration) IObservable[uint] {
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
				subscriber.Send() <- newData(index)
				index++
			}
		}
	})
}

func Scheduled[T any](item T, items ...T) IObservable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(subscriber Subscriber[T]) {
		for _, item := range items {
			select {
			// If receiver tell sender to stop, we should terminate the send operation
			case <-subscriber.Closed():
				return
			case subscriber.Send() <- newData(item):
			}
		}

		subscriber.Send() <- newComplete[T]()
	})
}

func Timer[T any, N constraints.Unsigned](start, interval N) IObservable[N] {
	return newObservable(func(subscriber Subscriber[N]) {
		latest := start

		for {
			select {
			case <-subscriber.Closed():
				return
			case <-time.After(time.Duration(latest)):
				subscriber.Send() <- newData(latest)
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
		// first.subscribeOn(func(a A) {
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
		// second.subscribeOn(func(b B) {
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
