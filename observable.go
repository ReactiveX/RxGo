package rxgo

import (
	"context"
	"time"

	"golang.org/x/exp/constraints"
)

type ObservableFunc[T any] func(subscriber Subscriber[T])

func newObservable[T any](obs ObservableFunc[T]) IObservable[T] {
	return &observableWrapper[T]{source: obs}
}

type observableWrapper[T any] struct {
	source ObservableFunc[T]
}

func (o *observableWrapper[T]) Subscribe(
	onNext func(T),
	onError func(error),
	onComplete func(),
) Subscription {
	ctx := context.Background()
	subcriber := NewSafeSubscriber(onNext, onError, onComplete)
	go o.source(subcriber)
	go consumeStreamUntil(ctx, subcriber, func() {})
	return subcriber
}

func (o *observableWrapper[T]) SubscribeSync(
	onNext func(T),
	onError func(error),
	onComplete func(),
) {
	ctx := context.Background()
	dispose := make(chan struct{})
	subcriber := NewSafeSubscriber(onNext, onError, onComplete)
	go o.source(subcriber)
	go consumeStreamUntil(ctx, subcriber, func() { close(dispose) })
	<-dispose
}

func consumeStreamUntil[T any](ctx context.Context, sub *safeSubscriber[T], finalizer func()) {
	defer finalizer()
	defer sub.Unsubscribe()

observe:
	for {
		select {
		// If context cancelled, shut down everything
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				sub.dst.Error(ctx.Err())
			}
			break observe
		case item, ok := <-sub.ForEach():
			if !ok {
				sub.dst.Complete()
				return
			}

			if err := item.Err(); err != nil {
				sub.dst.Error(err)
			} else {
				sub.dst.Next(item.Value())
			}
		}
	}
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
