package eventstream

import (
	"github.com/jochasinga/grx"
)

type EventStreamError struct {
	grx.BaseError
}

func NewError(code grx.ErrorType) {
	return &EventStreamError{
		Code:    code,
		Message: code.String(),
	}
}
