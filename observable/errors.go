package observable

import (
	"github.com/jochasinga/grx"
)

type ObservableError struct {
	grx.BaseError
}

func NewError(code grx.ErrorType) {
	return &ObservableError{
		Code:    code,
		Message: code.String(),
	}
}
