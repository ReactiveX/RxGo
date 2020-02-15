package rxgo

import "context"

type (
	Func    func(interface{}) (interface{}, error)
	Handler func(ctx context.Context, src <-chan Item, dst chan<- Item)

	NextFunc func(interface{})
	ErrFunc  func(error)
	DoneFunc func()

	Item struct {
		Value interface{}
		Err   error
	}
)

func (i Item) IsError() bool {
	return i.Err != nil
}

func FromValue(i interface{}) Item {
	return Item{Value: i}
}

func FromError(err error) Item {
	return Item{Err: err}
}
