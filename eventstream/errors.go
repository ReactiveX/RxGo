package eventstream

import (
	"github.com/jochasinga/grx/errors"
)

type EventStreamError struct {
	errors.BaseError
}

func NewError(code errors.ErrorCode) EventStreamError {
	baseErr := errors.New(code)
	return EventStreamError{baseErr}
}
