package observable

import (
	"github.com/jochasinga/grx/errors"
)

type ObservableError struct {
	errors.BaseError
}

func NewError(code errors.ErrorCode) ObservableError {
	baseError := errors.New(code)
	return ObservableError{baseError}
}
