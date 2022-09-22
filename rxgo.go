package rxgo

import (
	"context"
	"sync"
)

type (
	OnNextFunc[T any] func(T)

	// OnErrorFunc defines a function that computes a value from an error.
	OnErrorFunc                   func(error)
	OnCompleteFunc                func()
	FinalizerFunc                 func()
	ErrorFunc                     func() error
	OperatorFunc[I any, O any]    func(source Observable[I]) Observable[O]
	DurationFunc[T any, R any]    func(value T) Observable[R]
	PredicateFunc[T any]          func(value T, index uint) bool
	ProjectionFunc[T any, R any]  func(value T, index uint) Observable[R]
	ComparerFunc[A any, B any]    func(prev A, curr B) int8
	ComparatorFunc[A any, B any]  func(prev A, curr B) bool
	AccumulatorFunc[A any, V any] func(acc A, value V, index uint) (A, error)
	ObservableFunc[T any]         func(subscriber Subscriber[T])
)

type Observable[T any] interface {
	SubscribeWith(subscriber Subscriber[T])
	SubscribeOn(finalizer ...func()) Subscriber[T]
	SubscribeSync(onNext func(v T), onError func(err error), onComplete func())
	// Subscribe(onNext func(T), onError func(error), onComplete func()) Subscription
}

type GroupedObservable[K comparable, R any] interface {
	Observable[R] // Inherit from observable
	Key() K
}

type Subscription interface {
	// to unsubscribe the stream
	Unsubscribe()
}

type Observer[T any] interface {
	Next(T)
	Error(error)
	Complete()
}

type Subscriber[T any] interface {
	Stop()
	Send() chan<- Notification[T]
	ForEach() <-chan Notification[T]
	Closed() <-chan struct{}
	// Unsubscribe()
	// Observer[T]
}

type Subject[T any] interface {
	Subscriber[T]
	Subscription
}

func newObservable[T any](obs ObservableFunc[T]) Observable[T] {
	return &observableWrapper[T]{source: obs}
}

type observableWrapper[T any] struct {
	source    ObservableFunc[T]
	connector func() Subject[T]
}

var _ Observable[any] = (*observableWrapper[any])(nil)

func (o *observableWrapper[T]) SubscribeWith(subscriber Subscriber[T]) {
	o.source(subscriber)
}

func (o *observableWrapper[T]) SubscribeOn(cb ...func()) Subscriber[T] {
	var subscriber Subject[T]
	if o.connector != nil {
		subscriber = o.connector()
	} else {
		subscriber = NewSubscriber[T]()
	}
	finalizer := func() {}
	if len(cb) > 0 {
		finalizer = cb[0]
	}
	go func() {
		defer subscriber.Unsubscribe()
		defer finalizer()
		o.source(subscriber)
	}()
	return subscriber
}

func (o *observableWrapper[T]) SubscribeSync(onNext func(T), onError func(error), onComplete func()) {
	ctx := context.Background()
	subscriber := NewSafeSubscriber(onNext, onError, onComplete)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		o.source(subscriber)
	}()
	go func() {
		defer wg.Done()
		consumeStreamUntil(ctx, subscriber, func() {})
	}()
	wg.Wait()
}

func consumeStreamUntil[T any](ctx context.Context, sub *safeSubscriber[T], finalizer FinalizerFunc) {
	defer sub.Unsubscribe()
	defer finalizer()

observe:
	for {
		select {
		// If context cancelled, shut down everything
		// if err := ctx.Err(); err != nil {
		// 	sub.dst.Error(ctx.Err())
		// }
		case <-sub.Closed():
			break observe

		case item, ok := <-sub.ForEach():
			if !ok {
				break observe
			}

			if item.Done() {
				sub.dst.Complete()
				break observe
			}

			if err := item.Err(); err != nil {
				sub.dst.Error(err)
				break observe
			}

			sub.dst.Next(item.Value())
		}
	}
}
