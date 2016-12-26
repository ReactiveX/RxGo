package handlers

import (
	"github.com/jochasinga/grx/errors"
)

type HandlerError struct {
	errors.BaseError
}

func NewError(code errors.ErrorCode) HandlerError {
	baseErr := errors.New(code)
	return HandlerError{baseErr}
}
