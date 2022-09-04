package rxgo

type DataValuer[T any] interface {
	Value() T
	Err() error
	Done() bool
}

type streamData[T any] struct {
	v    T
	err  error
	done bool
}

var _ DataValuer[any] = (*streamData[any])(nil)

func (d streamData[T]) Value() T {
	return d.v
}

func (d streamData[T]) Err() error {
	return d.err
}

func (d streamData[T]) Done() bool {
	return d.done
}

func emitData[T any](v T, ch chan<- DataValuer[T]) {
	ch <- &streamData[T]{v: v}
}

func emitError[T any](err error, ch chan<- DataValuer[T]) {
	ch <- &streamData[T]{err: err}
}

func emitDone[T any](ch chan<- DataValuer[T]) {
	ch <- &streamData[T]{done: true}
}
