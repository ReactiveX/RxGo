package rxgo

import "time"

func NEVER[T any]() IObservable[T] {
	return newObservable(func(obs Subscriber[T]) {})
}

func EMPTY[T any]() IObservable[T] {
	return newObservable(func(obs Subscriber[T]) {
		obs.Complete()
	})
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(duration time.Duration) IObservable[uint] {
	return newObservable(func(obs Subscriber[uint]) {
		index := uint(0)
		for {
			time.Sleep(duration)
			obs.Next(index)
			index++
		}
	})
}

func Scheduled[T any](item T, items ...T) IObservable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(obs Subscriber[T]) {
		for _, item := range items {
			nextOrError(obs, item)
		}
		obs.Complete()
	})
}

func nextOrError[T any](sub Subscriber[T], v T) {
	switch vi := any(v).(type) {
	case error:
		sub.Error(vi)
	default:
		sub.Next(v)
	}
}
