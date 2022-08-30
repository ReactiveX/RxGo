package rxgo

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// An Observable that emits no items to the Observer and never completes.
func NEVER[T any]() IObservable[T] {
	return newObservable(func(sub Subscriber[T]) {})
}

// A simple Observable that emits no items to the Observer and immediately
// emits a complete notification.
func EMPTY[T any]() IObservable[T] {
	return newObservable(func(sub Subscriber[T]) {
		sub.Complete()
	})
}

func ThrownError[T any](factory func() error) IObservable[T] {
	return newObservable(func(sub Subscriber[T]) {
		sub.Error(factory())
	})
}

// Creates an Observable that emits a sequence of numbers within a specified range.
func Range[T constraints.Unsigned](start, count T) IObservable[T] {
	end := start + count
	return newObservable(func(sub Subscriber[T]) {
		var index uint
		for i := start; i < end; i++ {
			sub.Next(i)
			index++
		}
		sub.Complete()
	})
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(duration time.Duration) IObservable[uint] {
	return newObservable(func(sub Subscriber[uint]) {
		var index uint
		for {
			time.Sleep(duration)
			sub.Next(index)
			index++
		}
	})
}

func Scheduled[T any](item T, items ...T) IObservable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(sub Subscriber[T]) {
		for _, item := range items {
			nextOrError(sub, item)
		}
		sub.Complete()
	})
}

func Timer[T any, N constraints.Unsigned](start, due N) IObservable[N] {
	return newObservable(func(sub Subscriber[N]) {
		end := start + due
		for i := N(0); i < end; i++ {
			sub.Next(end)
			time.Sleep(time.Duration(due))
		}
		// sub.Complete()
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
