package handlers

import (
	"github.com/jochasinga/grx"
)

type HandlerError struct {
	grx.BaseError
}

func NewError(code grx.ErrorType) {
	return &HandlerError{
		Code:    code,
		Message: code.String(),
	}
}
