package errors

import "fmt"

// ErrorType serves as error code for the error enum
type Code uint32

const (
	EndOfIteratorError Code = iota
	NilObservableError
	NilEventStreamError
	NilObserverError
	UndefinedError
)

// BaseError provides a base template for more package-specific errors
type BaseError struct {
	code    Code
	message string
}

func New(code Code, msg ...string) BaseError {
	err := BaseError{
		code: code,
	}
	if len(msg) > 0 {
		err.message = msg[len(msg)-1]
		return err
	}
	err.message = code.String()
	return err
}

// Error returns an error string to implement the error interface
func (err *BaseError) Error() string {
	return fmt.Sprintf("%d - %s", err.code, err.message)
}

func (err *BaseError) Code() Code {
	return err.code
}
