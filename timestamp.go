package rxgo

import "time"

type Timestamper[T any] interface {
	Value() T
	Time() time.Time
}

type ts[T any] struct {
	v T
	t time.Time
}

var _ Timestamper[any] = (*ts[any])(nil)

func NewTimestamp[T any](value T) Timestamper[T] {
	return &ts[T]{v: value, t: time.Now()}
}

func (t *ts[T]) Value() T {
	return t.v
}

func (t *ts[T]) Time() time.Time {
	return t.t
}
