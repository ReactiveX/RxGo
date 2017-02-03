package errors

import "fmt"

// ErrorType serves as error code for the error enum
type ErrorCode uint32

const (
	EndOfIteratorError ErrorCode = iota + 1
	HandlerError
	ObservableError
	ObserverError
	IterableError
	UndefinedError
)

// BaseError provides a base template for more package-specific errors
type BaseError struct {
	code    ErrorCode
	message string
}

func New(code ErrorCode, msg ...string) BaseError {
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
func (err BaseError) Error() string {
	return fmt.Sprintf("%d - %s", err.code, err.message)
}

func (err BaseError) Code() int {
	return int(err.code)
}
