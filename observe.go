package rxgo

import (
	"context"
)

type ObservableFunc[T any] func(subscriber Subscriber[T])

func newObservable[T any](obs ObservableFunc[T]) IObservable[T] {
	return &observableWrapper[T]{source: obs}
}

type observableWrapper[T any] struct {
	source ObservableFunc[T]
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
	go consumeStreamUntil(ctx, dispose, subcriber)
	<-dispose
}

func consumeStreamUntil[T any](ctx context.Context, dispose chan struct{}, sub *safeSubscriber[T]) {
	defer close(dispose)
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
