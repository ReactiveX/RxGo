package rxgo

type Optional[T any] interface {
	MustGet() T
	OrElse(fallback T) T
	IsNone() bool
	Get() (T, bool)
}

type optional[T any] struct {
	none bool
	v    T
}

var _ Optional[any] = (*optional[any])(nil)

func (o optional[T]) MustGet() T {
	if o.none {
		panic("rxgo: option has no value")
	}

	return o.v
}

func (o optional[T]) OrElse(fallback T) T {
	if o.none {
		return fallback
	}
	return o.v
}

func (o optional[T]) IsNone() bool {
	return o.none
}

func (o optional[T]) Get() (T, bool) {
	if o.none {
		return *new(T), false
	}

	return o.v, true
}

func Some[T any](v T) Optional[T] {
	return optional[T]{v: v}
}

func None[T any]() Optional[T] {
	return optional[T]{none: true}
}
