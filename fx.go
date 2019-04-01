// Package fx provides predicate-like function types to be used with operators
// such as Map, Filter, Scan, and Start.
package rxgo

// Comparison is the results of the Comparator function.
type Comparison uint32

const (
	Equals Comparison = iota
	Greater
	Smaller
)

type (
	// Predicate defines a func that returns a bool from an input value.
	Predicate func(interface{}) bool

	// BiPredicate defines a func that returns a bool from two input values.
	BiPredicate func(interface{}, interface{}) bool

	// Comparator defines a func that returns:
	// Equals if two elements are equals
	// Greater if the first argument is greater than the second
	// Smaller if the first argument is smaller than the second
	Comparator func(interface{}, interface{}) Comparison

	// Consumer defines a func that accepts a single value
	Consumer func(interface{})

	// ErrorFunction defines a function that computes a value from an error.
	ErrorFunction func(error) interface{}

	// ErrorToObservableFunction defines a function that computes an observable from an error.
	ErrorToObservableFunction func(error) Observable

	// Function defines a function that computes a value from an input value.
	Function func(interface{}) interface{}

	// Function2 defines a function that computes a value from two input values.
	Function2 func(interface{}, interface{}) interface{}

	// Supplier defines a function that supplies a result from nothing.
	Supplier func() interface{}

	// ChanProducer defines a function that produces results into a channel.
	ChanProducer func(chan interface{})
)
