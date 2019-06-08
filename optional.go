package rxgo

import "github.com/pkg/errors"

var emptyOptional = new(empty)

// Optional defines a container for empty values
type Optional interface {
	// Get returns the content and an optional error is the optional is empty
	Get() (interface{}, error)
	// IsEmpty returns whether the optional is empty
	IsEmpty() bool
}

type some struct {
	content interface{}
}

type empty struct {
}

// Get returns the content and an optional error is the optional is empty
func (s *some) Get() (interface{}, error) {
	return s.content, nil
}

// IsEmpty returns whether the optional is empty
func (s *some) IsEmpty() bool {
	return false
}

// Get returns the content and an optional error is the optional is empty
func (e *empty) Get() (interface{}, error) {
	return nil, errors.Wrap(&NoSuchElementError{}, "empty does not contain any element")
}

// IsEmpty returns whether the optional is empty
func (e *empty) IsEmpty() bool {
	return true
}

// Of returns a non-empty optional
func Of(data interface{}) Optional {
	return &some{
		content: data,
	}
}

// Empty returns an empty optional
func EmptyOptional() Optional {
	return emptyOptional
}
