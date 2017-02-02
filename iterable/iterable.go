package iterable

import "github.com/jochasinga/grx/errors"

// Emittable is a high-level alias for empty interface which may any underlying type including error
type Iterable <-chan interface{}

func (it Iterable) Next() (interface{}, error) {
	if next, ok := <-it; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

func From(any interface{}) (Iterable, error) {
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
