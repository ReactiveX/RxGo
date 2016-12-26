package observable

import (
	"github.com/jochasinga/grx/errors"
)

type ObservableError struct {
	errors.BaseError
}

func NewError(code grx.ErrorType) {
	return &ObservableError{
		Code:    code,
		Message: code.String(),
	}
}
