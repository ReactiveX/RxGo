package rxgo

// CancelledIteratorError is triggered when an iterator is canceled
type CancelledIteratorError struct {
}

// EndOfIteratorError is triggered when an iterator is complete
type EndOfIteratorError struct {
}

// IllegalInputError is triggered when the observable receives an illegal input
type IllegalInputError struct {
	reason string
}

// IndexOutOfBoundError is triggered when the observable cannot access to the specified index
type IndexOutOfBoundError struct {
}

// NoSuchElementError is triggered when an optional does not contain any element
type NoSuchElementError struct {
}

// TimeoutError is triggered when a timeout occurs
type TimeoutError struct {
}

func (e *CancelledIteratorError) Error() string {
	return "CancelledIteratorError"
}

func (e *EndOfIteratorError) Error() string {
	return "EndOfIteratorError"
}

func (e *IllegalInputError) Error() string {
	return e.reason
}

func (e *IndexOutOfBoundError) Error() string {
	return "IndexOutOfBoundError"
}

func (e *NoSuchElementError) Error() string {
	return "NoSuchElementError"
}

func (e *TimeoutError) Error() string {
	return "TimeoutError"
}
