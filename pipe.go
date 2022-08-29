package rxgo

type IObservable[T any] interface {
	// Subscribe(onNext func(T), onError func(error), onComplete func(), opts ...any) Subscription
	SubscribeSync(onNext func(T), onError func(error), onComplete func())
}

type Subscription interface {
	Unsubscribe()
	// Done() <-chan struct{}
}

type Observer[T any] interface {
	Next(T)
	Error(error)
	Complete()
}

type Subscriber[T any] interface {
	Closed() bool
	Observer[T]
}

type OperatorFunc[I any, O any] func(IObservable[I]) IObservable[O]

func Pipe1[S any, O1 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
) IObservable[O1] {
	return f1(stream)
}

func Pipe2[S any, O1 any, O2 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
) IObservable[O2] {
	return f2(f1(stream))
}

func Pipe3[S any, O1 any, O2 any, O3 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
) IObservable[O3] {
	return f3(f2(f1(stream)))
}
