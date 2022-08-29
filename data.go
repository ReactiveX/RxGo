package rxgo

type DataValuer[T any] interface {
	Value() T
	Err() error
}

type streamData[T any] struct {
	v   T
	err error
}

var _ DataValuer[any] = (*streamData[any])(nil)

func (d streamData[T]) Value() T {
	return d.v
}

func (d streamData[T]) Err() error {
	return d.err
}

func emitData[T any](v T, ch chan<- DataValuer[T]) {
	ch <- &streamData[T]{v: v}
}

func emitError[T any](err error, ch chan<- DataValuer[T]) {
	ch <- &streamData[T]{err: err}
}
