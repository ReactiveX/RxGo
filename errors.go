package rxgo

type CancelledIteratorError struct {
}

type EndOfIteratorError struct {
}

type IllegalInputError struct {
	reason string
}

type IndexOutOfBoundError struct {
}

type NoSuchElementError struct {
}

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
