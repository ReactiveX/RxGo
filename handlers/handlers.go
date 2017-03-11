// Package handlers provides handler types which implements EventHandler.
package handlers

type (
	// NextFunc handles a next item in a stream.
	NextFunc func(interface{})

	// ErrFunc handles an error in a stream.
	ErrFunc func(error)

	// DoneFunc handles the end of a stream.
	DoneFunc func()
)

func AsNextFunc(handler interface{}) (NextFunc, bool) {
	switch handler := handler.(type) {
	case NextFunc:
		return handler, true
	case func(interface{}):
		return NextFunc(handler), true
	}

	return nil, false
}

func AsErrFunc(handler interface{}) (ErrFunc, bool) {
	switch handler := handler.(type) {
	case ErrFunc:
		return handler, true
	case func(error):
		return ErrFunc(handler), true
	}

	return nil, false
}
func AsDoneFunc(handler interface{}) (DoneFunc, bool) {
	switch handler := handler.(type) {
	case DoneFunc:
		return handler, true
	case func():
		return DoneFunc(handler), true
	}

	return nil, false
}
