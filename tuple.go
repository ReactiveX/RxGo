package rxgo

type Tuple[A any, B any] interface {
	First() A
	Second() B
}

type pair[A any, B any] struct {
	first  A
	second B
}

func (p pair[A, B]) First() A {
	return p.first
}

func (p pair[A, B]) Second() B {
	return p.second
}

// Create a tuple using first and second value.
func NewTuple[A any, B any](a A, b B) Tuple[A, B] {
	return &pair[A, B]{first: a, second: b}
}
