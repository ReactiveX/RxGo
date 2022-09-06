package rxgo

import (
	"context"
	"sync"
)

type ObservableFunc[T any] func(subscriber Subscriber[T])

func newObservable[T any](obs ObservableFunc[T]) IObservable[T] {
	return &observableWrapper[T]{source: obs}
}

type observableWrapper[T any] struct {
	source ObservableFunc[T]
}

var _ IObservable[any] = (*observableWrapper[any])(nil)

func (o *observableWrapper[T]) SubscribeOn(cb ...func()) Subscriber[T] {
	subscriber := NewSubscriber[T]()
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

func (o *observableWrapper[T]) SubscribeSync(
	onNext func(T),
	onError func(error),
	onComplete func(),
) {
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
