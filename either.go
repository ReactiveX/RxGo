package rxgo

type Either[L, R any] interface {
	IsLeft() bool
	IsRight() bool
	Left() (L, bool)
	Right() (R, bool)
}

func Left[L, R any](value L) Either[L, R] {
	return &either[L, R]{left: value, isLeft: true}
}

type either[L, R any] struct {
	isLeft bool
	left   L
	right  R
}

func (e *either[L, R]) IsLeft() bool {
	return e.isLeft
}

func (e *either[L, R]) IsRight() bool {
	return !e.isLeft
}

func (e *either[L, R]) Left() (L, bool) {
	return e.left, e.IsLeft()
}

func (e *either[L, R]) Right() (R, bool) {
	return e.right, e.IsRight()
}
