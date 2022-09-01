package rxgo

import "fmt"

var (
	ErrEmpty              = fmt.Errorf("rxgo: empty value")
	ErrNotFound           = fmt.Errorf("rxgo: no values match")
	ErrSequence           = fmt.Errorf("rxgo: too many values match")
	ErrArgumentOutOfRange = fmt.Errorf("rxgo: argument out of range")
)

// IllegalInputError is triggered when the observable receives an illegal input.
type IllegalInputError struct {
	error string
}

func (e IllegalInputError) Error() string {
	return "illegal input: " + e.error
}

// IndexOutOfBoundError is triggered when the observable cannot access to the specified index.
type IndexOutOfBoundError struct {
	error string
}

func (e IndexOutOfBoundError) Error() string {
	return "index out of bound: " + e.error
}
