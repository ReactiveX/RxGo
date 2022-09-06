package rxgo

type (
	OnNextFunc[T any] func(T)
	// OnErrorFunc defines a function that computes a value from an error.
	OnErrorFunc                func(error)
	OnCompleteFunc             func()
	FinalizerFunc              func()
	OperatorFunc[I any, O any] func(IObservable[I]) IObservable[O]
)

// FIXME: please rename it to `Observable`
type IObservable[T any] interface {
	SubscribeOn(...func()) Subscriber[T]
	SubscribeSync(onNext func(T), onError func(error), onComplete func())
	// Subscribe(onNext func(T), onError func(error), onComplete func()) Subscription
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
	Send() chan<- DataValuer[T]
	ForEach() <-chan DataValuer[T]
	Closed() <-chan struct{}
	// Unsubscribe()
	// Observer[T]
}

// Pipe
func Pipe[S any, O1 any](
	stream IObservable[S],
	f1 OperatorFunc[S, any],
	f ...OperatorFunc[any, any],
) IObservable[any] {
	result := f1(stream)
	for _, cb := range f {
		result = cb(result)
	}
	return result
}

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

func Pipe4[S any, O1 any, O2 any, O3 any, O4 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
) IObservable[O4] {
	return f4(f3(f2(f1(stream))))
}

func Pipe5[S any, O1 any, O2 any, O3 any, O4 any, O5 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
) IObservable[O5] {
	return f5(f4(f3(f2(f1(stream)))))
}

func Pipe6[S any, O1 any, O2 any, O3 any, O4 any, O5 any, O6 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
) IObservable[O6] {
	return f6(f5(f4(f3(f2(f1(stream))))))
}

func Pipe7[S any, O1 any, O2 any, O3 any, O4 any, O5 any, O6 any, O7 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
) IObservable[O7] {
	return f7(f6(f5(f4(f3(f2(f1(stream)))))))
}

func Pipe8[S any, O1 any, O2 any, O3 any, O4 any, O5 any, O6 any, O7 any, O8 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
	f8 OperatorFunc[O7, O8],
) IObservable[O8] {
	return f8(f7(f6(f5(f4(f3(f2(f1(stream))))))))
}

func Pipe9[S any, O1 any, O2 any, O3 any, O4 any, O5 any, O6 any, O7 any, O8 any, O9 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
	f8 OperatorFunc[O7, O8],
	f9 OperatorFunc[O8, O9],
) IObservable[O9] {
	return f9(f8(f7(f6(f5(f4(f3(f2(f1(stream)))))))))
}

func Pipe10[S any, O1 any, O2 any, O3 any, O4 any, O5 any, O6 any, O7 any, O8 any, O9 any, O10 any](
	stream IObservable[S],
	f1 OperatorFunc[S, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
	f8 OperatorFunc[O7, O8],
	f9 OperatorFunc[O8, O9],
	f10 OperatorFunc[O9, O10],
) IObservable[O10] {
	return f10(f9(f8(f7(f6(f5(f4(f3(f2(f1(stream))))))))))
}
