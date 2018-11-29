// Package fx provides predicate-like function types to be used with operators
// such as Map, Filter, Scan, and Start.
package fx

type (
	// Predicate defines a func that returns a bool from an input value.
	Predicate func(interface{}) bool

	// Consumer defines a func that accepts a single value
	Consumer func(interface{})

	// Supplier defines a function that supplies a result from nothing.
	Supplier func() interface{}

	// Function defines a function that computes a value from an input value.
	Function func(interface{}) interface{}

	// Function2 defines a function that computes a value from two input values.
	Function2 func(interface{}, interface{}) interface{}
)
