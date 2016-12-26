package single

import "github.com/jochasinga/grx/bases"

// Single either emits an Item or an error
type Single struct {
	bases.Item
}

// Emit type-assert and return an Item type or an error
func (s Single) Emit() (bases.Item, error) {
	if err, ok := s.Item.(error); ok {
		return nil, err
	}
	return s.Item, nil
}

// Next always return an EndOfIteratorError
func (s Single) Next() (interface{}, error) {
	return NewError(grx.EndOfIteratorError)
}

// HasNext always return false
func (s Single) HasNext() bool {
	return false
}

// New creates a Single from Item
func New(fs ...func(interface{})) Single {
	s := new(Single)
	for _, f := range fs {
		fn(s)
	}
	return *s
}

// From creates a Single from an empty interface type
func From(value interface{}) Single {
	return Single{value}
}
