package handlers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jochasinga/rx"
	"github.com/stretchr/testify/assert"
)

func TestHandlersImplementEventHandler(t *testing.T) {
	assert := assert.New(t)
	assert.Implements((*rx.EventHandler)(nil), (*NextFunc)(nil))
	assert.Implements((*rx.EventHandler)(nil), (*ErrFunc)(nil))
	assert.Implements((*rx.EventHandler)(nil), (*DoneFunc)(nil))
}

func TestNextFuncHandleMethod(t *testing.T) {
	assert := assert.New(t)

	var (
		num  = 1
		word = "Hello"
		pi   = 3.1416
		ru   = 'ด'
		err  = errors.New("Anonymous error")
	)

	samples := []interface{}{}

	// Append each item to samples slice
	nextf := NextFunc(func(item interface{}) {
		samples = append(samples, item)
	})

	nextHandleTests := []interface{}{num, word, pi, ru, err}

	for _, tt := range nextHandleTests {
		nextf.Handle(tt)
	}

	expected := []interface{}{num, word, pi, ru}
	assert.Exactly(samples, expected)
}

func TestErrFuncHandleMethod(t *testing.T) {
	assert := assert.New(t)

	var (
		num  = 1
		word = "Hello"
		pi   = 3.1416
		ru   = 'ด'
	)

	samples := []error{}

	errf := ErrFunc(func(err error) {
		samples = append(samples, err)
	})

	errHandleTests := []interface{}{
		fmt.Errorf("Integer %d error", num),
		fmt.Errorf("String %s error", word),
		pi,
		fmt.Errorf("Rune %U error", ru),
	}

	for _, tt := range errHandleTests {
		errf.Handle(tt)
	}

	expected := []error{
		fmt.Errorf("Integer %d error", num),
		fmt.Errorf("String %s error", word),
		fmt.Errorf("Rune %U error", ru),
	}

	assert.Exactly(samples, expected)
}

func TestDoneFuncHandleMethod(t *testing.T) {
	assert := assert.New(t)

	var (
		num  = 1
		word = "Hello"
		pi   = 3.1416
		ru   = 'ด'
	)

	samples := []string{}

	errf := DoneFunc(func() {
		samples = append(samples, "DONE")
	})

	doneHandleTests := []interface{}{num, word, pi, ru}

	for n, tt := range doneHandleTests {
		errf.Handle(tt)
		assert.Equal(samples[n], "DONE")
	}
}
