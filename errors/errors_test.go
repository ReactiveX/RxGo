package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errcodes = []ErrorCode{
	EndOfIteratorError,
	HandlerError,
	ObservableError,
	ObserverError,
	IterableError,
	UndefinedError,
}

func TestErrorCodes(t *testing.T) {
	for i, errcode := range errcodes {
		assert.Equal(t, ErrorCode(i+1), errcode)
	}
}

func TestBaseErrorImplementsError(t *testing.T) {
	for _, errcode := range errcodes {
		err := New(errcode)
		assert.Implements(t, (*error)(nil), err)
	}
}

func TestBaseErrorWithDefaultMessage(t *testing.T) {
	for _, errcode := range errcodes {
		err := New(errcode)
		assert.Equal(t, err.code.String(), err.message)
	}
}

func TestBaseErrorWithCustomMessage(t *testing.T) {
	msg := "Custom error message"
	for _, errcode := range errcodes {
		err := New(errcode, msg)
		assert.Equal(t, msg, err.message)
	}
}

func TestErrorMethod(t *testing.T) {
	for i, errcode := range errcodes {
		err := New(errcode)
		msg := fmt.Sprintf("%d - %s", i+1, err.code.String())
		assert.Equal(t, msg, err.Error())
	}
}

func TestCodeMethod(t *testing.T) {
	for i, errcode := range errcodes {
		err := New(errcode)
		assert.EqualValues(t, i+1, err.Code())
	}
}
