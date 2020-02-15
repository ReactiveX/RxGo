package rxgo

import "context"

type (
	Func    func(interface{}) (interface{}, error)
	Handler func(ctx context.Context, nextSrc <-chan interface{}, errsSrc <-chan error, nextDst chan<- interface{}, errsDst chan<- error)

	NextFunc func(interface{})
	ErrFunc  func(error)
	DoneFunc func()
)
