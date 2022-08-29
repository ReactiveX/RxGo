package rxgo

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func NEVER[T any]() IObservable[T] {
	return newObservable(func(obs Subscriber[T]) {})
}

// A simple Observable that emits no items to the Observer and immediately emits a complete notification.
func EMPTY[T any]() IObservable[T] {
	return newObservable(func(obs Subscriber[T]) {
		obs.Complete()
	})
}

// Creates an Observable that emits a sequence of numbers within a specified range.
func Range[T constraints.Unsigned](start, count T) IObservable[T] {
	end := start + count
	return newObservable(func(obs Subscriber[T]) {
		index := uint(0)
		for i := start; i < end; i++ {
			obs.Next(i)
			index++
		}
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
