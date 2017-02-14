// Package iterable provides an Iterable type that is capable of converting
// sequences of empty interface such as slice and channel to an Iterator.
package iterable

import "github.com/reactivex/rxgo/errors"

// Iterable converts channel and slice into an Iterator.
type Iterable <-chan interface{}

// Next returns the next element in an Iterable sequence and an
// error when it reaches the end. Next registers Iterable to Iterator.
func (it Iterable) Next() (interface{}, error) {
	if next, ok := <-it; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// New creates a new Iterable from a slice or a channel of empty interface.
func New(any interface{}) (Iterable, error) {
	switch any := any.(type) {
	case []interface{}:
		c := make(chan interface{}, len(any))
		go func() {
			for _, val := range any {
				c <- val
			}
			close(c)
		}()
		return Iterable(c), nil
	case chan interface{}:
		return Iterable(any), nil
	case <-chan interface{}:
		return Iterable(any), nil
	default:
		return nil, errors.New(errors.IterableError)
	}
}
