package rxgo

type Optional[T any] struct {
	none bool
	v    T
}

func (o Optional[T]) MustGet() T {
	if !o.none {
		panic("rxgo: option has no value")
	}

	return o.v
}

func (o Optional[T]) OrElse(fallback T) T {
	if o.none {
		return fallback
	}
	return o.v
}

func (o Optional[T]) None() bool {
	return o.none
}

func (o Optional[T]) Get() (T, bool) {
	if !o.none {
		return *new(T), false
	}

	return o.v, true
}
