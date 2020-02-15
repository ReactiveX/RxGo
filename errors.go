package rxgo

// IllegalInputError is triggered when the observable receives an illegal input
type IllegalInputError struct {
}

func (e *IllegalInputError) Error() string {
	return "illegal input"
}
