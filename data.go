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

func newData[T any](v T) DataValuer[T] {
	return &streamData[T]{v: v}
}

func newError[T any](err error) DataValuer[T] {
	return &streamData[T]{err: err}
}

func newComplete[T any]() DataValuer[T] {
	return &streamData[T]{done: true}
}
