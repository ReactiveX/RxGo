package rxgo

// CancelledIteratorError is triggered when an iterator is canceled
type CancelledIteratorError struct {
}

// IllegalInputError is triggered when the observable receives an illegal input
type IllegalInputError struct {
}

// IndexOutOfBoundError is triggered when the observable cannot access to the specified index
type IndexOutOfBoundError struct {
}

// NoSuchElementError is triggered when an element does not exist
type NoSuchElementError struct {
}

// CancelledSubscriptionError is triggered when a subscription was cancelled manually by an end user
type CancelledSubscriptionError struct {
}

func (e *CancelledIteratorError) Error() string {
	return "cancelled iterator"
}

func (e *IllegalInputError) Error() string {
	return "illegal input"
}

func (e *IndexOutOfBoundError) Error() string {
	return "index out of bound"
}

func (e *NoSuchElementError) Error() string {
	return "no such element"
}

func (e *CancelledSubscriptionError) Error() string {
	return "timeout"
}
