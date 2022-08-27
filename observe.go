package rxgo

type ObservableFunc[T any] func(obs Subscriber[T])

func newObservable[T any](obs ObservableFunc[T]) IObservable[T] {
	return &observableWrapper[T]{source: obs}
}

type observableWrapper[T any] struct {
	source ObservableFunc[T]
}
