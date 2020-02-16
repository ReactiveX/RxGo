package rxgo

// IllegalInputError is triggered when the observable receives an illegal input.
type IllegalInputError struct {
}

func (e *IllegalInputError) Error() string {
	return "illegal input"
}

// IndexOutOfBoundError is triggered when the observable cannot access to the specified index
type IndexOutOfBoundError struct {
}

func (e *IndexOutOfBoundError) Error() string {
	return "index out of bound"
}
