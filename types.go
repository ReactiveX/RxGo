package rxgo

import "context"

type (
	// BackpressureStrategy is the backpressure strategy type.
	BackpressureStrategy uint32
	// OnErrorStrategy is the Observable error strategy.
	OnErrorStrategy uint32

	operatorOptions struct {
		stop          func()
		resetIterable func(Iterable)
	}

	// Comparator defines a func that returns an int:
	// - 0 if two elements are equals
	// - A negative value if the first argument is less than the second
	// - A positive value if the first argument is greater than the second
	Comparator func(interface{}, interface{}) int
	// ItemToObservable defines a function that computes an observable from an item.
	ItemToObservable func(Item) Observable
	// ErrorToObservable defines a function that transforms an observable from an error.
	ErrorToObservable func(error) Observable
	// Func defines a function that computes a value from an input value.
	Func func(interface{}) (interface{}, error)
	// Func2 defines a function that computes a value from two input values.
	Func2 func(interface{}, interface{}) (interface{}, error)
	// FuncN defines a function that computes a value from N input values.
	FuncN func(...interface{}) interface{}
	// ErrorFunc defines a function that computes a value from an error.
	ErrorFunc func(error) interface{}
	// Predicate defines a func that returns a bool from an input value.
	Predicate func(interface{}) bool
	// Marshaller defines a marshaller type (interface{} to []byte).
	Marshaller func(interface{}) ([]byte, error)
	// Unmarshaller defines an unmarshaller type ([]byte to interface).
	Unmarshaller func([]byte, interface{}) error
	// Producer defines a producer implementation.
	Producer func(ctx context.Context, next chan<- Item, done func())
	// Supplier defines a function that supplies a result from nothing.
	Supplier func(ctx context.Context) Item
	// Disposed is a notification channel indicating when an Observable is closed.
	Disposed <-chan struct{}

	// NextFunc handles a next item in a stream.
	NextFunc func(interface{})
	// ErrFunc handles an error in a stream.
	ErrFunc func(error)
	// CompletedFunc handles the end of a stream.
	CompletedFunc func()

	operatorItem func(ctx context.Context, item Item, dst chan<- Item, operator operatorOptions)
	operatorEnd  func(ctx context.Context, dst chan<- Item)
)

const (
	// Block blocks until the channel is available.
	Block BackpressureStrategy = iota
	// Drop drops the message.
	Drop
)

const (
	// Stop is the default error strategy.
	// An operator will stop processing items on error.
	Stop OnErrorStrategy = iota
	// Continue means an operator will continue processing items after an error.
	Continue
)
