package rxgo

import "context"

type (
	// BackpressureStrategy is the backpressure strategy type
	BackpressureStrategy uint32

	// Function defines a function that computes a value from an input value.
	Function func(interface{}) (interface{}, error)
	// Handler defines a function implementing the handler logic for a stream.
	Handler func(ctx context.Context, src <-chan Item, dst chan<- Item)
	// Operator defines an operator function.
	Operator func(item Item, dst chan<- Item, stop func())
	// Predicate defines a func that returns a bool from an input value.
	Predicate func(interface{}) bool

	// NextFunc handles a next item in a stream.
	NextFunc func(interface{})
	// ErrFunc handles an error in a stream.
	ErrFunc func(error)
	// DoneFunc handles the end of a stream.
	DoneFunc func()

	// Item is a wrapper having either a value or an error.
	Item struct {
		Value interface{}
		Err   error
	}
)

const (
	// Block blocks until the channel is available
	Block BackpressureStrategy = iota
	// Drop drops the message
	Drop
)

// IsError checks if an item is an error.
func (i Item) IsError() bool {
	return i.Err != nil
}

// FromValue creates an item from a value.
func FromValue(i interface{}) Item {
	return Item{Value: i}
}

// FromError creates an item from an error.
func FromError(err error) Item {
	return Item{Err: err}
}
