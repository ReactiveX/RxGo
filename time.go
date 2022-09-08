package rxgo

import "time"

// Emits an object containing the current value, and the time that has passed
// between emitting the current value and the previous value, which is calculated by
// using the provided scheduler's now() method to retrieve the current time at each
// emission, then calculating the difference.
func WithTimeInterval[T any]() OperatorFunc[T, TimeInterval[T]] {
	return func(source IObservable[T]) IObservable[TimeInterval[T]] {
		var (
			pastTime = time.Now().UTC()
		)
		return createOperatorFunc(
			source,
			func(obs Observer[TimeInterval[T]], v T) {
				now := time.Now().UTC()
				obs.Next(NewTimeInterval(v, now.Sub(pastTime)))
				pastTime = now
			},
			func(obs Observer[TimeInterval[T]], err error) {
				obs.Error(err)
			},
			func(obs Observer[TimeInterval[T]]) {
				obs.Complete()
			},
		)
	}
}

// Attaches a UTC timestamp to each item emitted by an observable indicating
// when it was emitted
func WithTimestamp[T any]() OperatorFunc[T, Timestamp[T]] {
	return func(source IObservable[T]) IObservable[Timestamp[T]] {
		return createOperatorFunc(
			source,
			func(obs Observer[Timestamp[T]], v T) {
				obs.Next(NewTimestamp(v))
			},
			func(obs Observer[Timestamp[T]], err error) {
				obs.Error(err)
			},
			func(obs Observer[Timestamp[T]]) {
				obs.Complete()
			},
		)
	}
}

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
	return &ts[T]{v: value, t: time.Now().UTC()}
}

func (t ts[T]) Value() T {
	return t.v
}

func (t ts[T]) Time() time.Time {
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

func (t ti[T]) Value() T {
	return t.v
}

func (t ti[T]) Elapsed() time.Duration {
	return t.elapsed
}
