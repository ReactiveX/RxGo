package single

import (
	"github.com/jochasinga/grx/errors"
)

type SingleError struct {
	errors.BaseError
}

func NewError(code errors.ErrorCode) SingleError {
	baseError := errors.New(code)
	return SingleError{baseError}
}
