package rxgo

import "time"

type Timestamp[T any] interface {
	Value() T
	Time() time.Time
}

type TimeInterval[T any] interface {
	Value() T
	Elapsed() time.Duration
}

type ts[T any] struct {
	v T
	t time.Time
}

var _ Timestamp[any] = (*ts[any])(nil)

func NewTimestamp[T any](value T) Timestamp[T] {
	return &ts[T]{v: value, t: time.Now()}
}

func (t *ts[T]) Value() T {
	return t.v
}

func (t *ts[T]) Time() time.Time {
	return t.t
}

type ti[T any] struct {
	v       T
	elapsed time.Duration
}

var _ TimeInterval[any] = (*ti[any])(nil)

func NewTimeInterval[T any](value T, elasped time.Duration) TimeInterval[T] {
	return &ti[T]{v: value, elapsed: elasped}
}

func (t *ti[T]) Value() T {
	return t.v
}

func (t *ti[T]) Elapsed() time.Duration {
	return t.elapsed
}
