package rxgo

// Comparison is the results of the Comparator function.
type Comparison uint32

const (
	// Equals represents an equality
	Equals Comparison = iota
	// Greater represents an item greater than another
	Greater
	// Smaller represents an item smaller than another
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

	// FunctionN defines a function that computes a value from N input values.
	FunctionN func(...interface{}) interface{}

	// Supplier defines a function that supplies a result from nothing.
	Supplier func() interface{}

	// Channel defines a type representing chan interface{}
	Channel chan interface{}
)
