package rxgo

import "errors"

var (
	// An error thrown when an Observable or a sequence was queried but has no elements.
	ErrEmpty = errors.New("rxgo: empty value")
	// An error thrown when a value or values are missing from an observable sequence.
	ErrNotFound           = errors.New("rxgo: no values match")
	ErrSequence           = errors.New("rxgo: too many values match")
	ErrArgumentOutOfRange = errors.New("rxgo: argument out of range")
)
