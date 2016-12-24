package single

import (
	"github.com/jochasinga/grx"
)

type SingleError struct {
	grx.BaseError
}

func NewError(code grx.ErrorType) {
	return &SingleError{
		Code:    code,
		Message: code.String(),
	}
}
