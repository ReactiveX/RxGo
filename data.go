package rxgo

type Data[T any] struct {
	v   T
	err error
}

func (d Data[T]) Value() T {
	return d.v
}

func (d Data[T]) Err() error {
	return d.err
}

type DataValuer[T any] interface {
	Value() T
	Err() error
}
