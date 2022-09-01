package rxgo

import (
	"context"
	"sync"
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

var _ IObservable[any] = (*observableWrapper[any])(nil)

//nolint:golint,unused
func (o *observableWrapper[T]) subscribeOn(
	onNext func(T),
	onError func(error),
	onComplete func(),
	finalizer func(),
) Subscription {
	ctx := context.Background()
	subcriber := NewSafeSubscriber(onNext, onError, onComplete)
	go o.source(subcriber)
	go consumeStreamUntil(ctx, subcriber, finalizer)
	return subcriber
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

// Creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
func Defer[T any](factory func() IObservable[T]) IObservable[T] {
	return factory()
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
	return newObservable(func(subscriber Subscriber[uint]) {
		var index uint
		for {
			time.Sleep(duration)
			subscriber.Next(index)
			index++
		}
	})
}

func Scheduled[T any](item T, items ...T) IObservable[T] {
	items = append([]T{item}, items...)
	return newObservable(func(subscriber Subscriber[T]) {
		for _, item := range items {
			nextOrError(subscriber, item)
		}
		subscriber.Complete()
	})
}

func Timer[T any, N constraints.Unsigned](start, interval N) IObservable[N] {
	return newObservable(func(subscriber Subscriber[N]) {
		latest := start

		for {
			subscriber.Next(latest)
			time.Sleep(time.Duration(latest))
			latest = latest + interval
		}
	})
}

func CombineLatest[A any, B any](first IObservable[A], second IObservable[B]) IObservable[Tuple[A, B]] {
	return newObservable(func(subscriber Subscriber[Tuple[A, B]]) {
		var (
			latestA  A
			latestB  B
			hasValue = [2]bool{}
			allOk    = [2]bool{}
		)

		nextValue := func() {
			if hasValue[0] && hasValue[1] {
				subscriber.Next(NewTuple(latestA, latestB))
			}
		}
		checkComplete := func() {
			if allOk[0] && allOk[1] {
				subscriber.Complete()
			}
		}

		wg := new(sync.WaitGroup)
		wg.Add(2)
		first.subscribeOn(func(a A) {
			latestA = a
			hasValue[0] = true
			nextValue()
		}, subscriber.Error, func() {
			allOk[0] = true
			checkComplete()
		}, wg.Done)
		second.subscribeOn(func(b B) {
			latestB = b
			hasValue[1] = true
			nextValue()
		}, subscriber.Error, func() {
			allOk[1] = true
			checkComplete()
		}, wg.Done)
		wg.Wait()
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
