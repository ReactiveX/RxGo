package rxgo

import "context"

type (
	// BackpressureStrategy is the backpressure strategy type
	BackpressureStrategy uint32

	operatorOptions struct {
		stop          func()
		resetIterable func(Iterable)
	}

	// Comparator defines a func that returns:
	// - 0 if two elements are equals
	// - A negative value if the first argument is less than the second
	// - A positive value if the first argument is greater than the second
	Comparator func(interface{}, interface{}) int
	// ErrorToObservable defines a function that computes an observable from an error.
	ErrorToObservable func(error) Observable
	// Func defines a function that computes a value from an input value.
	Func func(interface{}) (interface{}, error)
	// ErrorFunc defines a function that computes a value from an error.
	ErrorFunc func(error) interface{}
	// Predicate defines a func that returns a bool from an input value.
	Predicate func(interface{}) bool
	// Marshaler defines a marshaler type (interface{} to []byte).
	Marshaler func(interface{}) ([]byte, error)
	// Unmarshaler defines an unmarshaler type ([]byte to interface).
	Unmarshaler func([]byte, interface{}) error
	// Scatter defines a scatter implementation
	Scatter func(ctx context.Context, next chan<- Item, done func())

	// NextFunc handles a next item in a stream.
	NextFunc func(interface{})
	// ErrFunc handles an error in a stream.
	ErrFunc func(error)
	// DoneFunc handles the end of a stream.
	DoneFunc func()

	operatorItem func(item Item, dst chan<- Item, operatorOpts operatorOptions)
	operatorEnd  func(dst chan<- Item)

	// Item is a wrapper having either a value or an error.
	Item struct {
		Value interface{}
		Err   error
	}
)

const (
	// Block blocks until the channel is available.
	Block BackpressureStrategy = iota
	// Drop drops the message.
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
