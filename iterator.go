package rx

// Iterator type is implemented by Iterable.
type Iterator interface {
	Next() (interface{}, error)
}
