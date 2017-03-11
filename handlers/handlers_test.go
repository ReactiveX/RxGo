package handlers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNextFuncHandleMethod(t *testing.T) {
	assert := assert.New(t)

	var (
		num  = 1
		word = "Hello"
		pi   = 3.1416
		ru   = 'ด'
	)

	samples := []interface{}{}

	// Append each item to samples slice
	nextf := NextFunc(func(item interface{}) {
		samples = append(samples, item)
	})

	nextHandleTests := []interface{}{num, word, pi, ru}

	for _, tt := range nextHandleTests {
		nextf(tt)
	}

	expected := []interface{}{num, word, pi, ru}
	assert.Exactly(samples, expected)
}

func TestErrFuncHandleMethod(t *testing.T) {
	assert := assert.New(t)

	var (
		num  = 1
		word = "Hello"
		ru   = 'ด'
	)

	samples := []error{}

	errf := ErrFunc(func(err error) {
		samples = append(samples, err)
	})

	errHandleTests := []error{
		fmt.Errorf("Integer %d error", num),
		fmt.Errorf("String %s error", word),
		fmt.Errorf("Rune %U error", ru),
	}

	for _, tt := range errHandleTests {
		errf(tt)
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

	for n := range doneHandleTests {
		errf()
		assert.Equal(samples[n], "DONE")
	}
}

func TestReturnsAsDoneFunc(t *testing.T) {
	// where
	type data struct{
		function interface{}
		casted bool
	}
	testData := [...]data {
		{ function: func() {},								casted: true,},
		{ function: DoneFunc(func() {}),					casted: true,},
		{ function: func(err error) {},						casted: false,},
		{ function: func(item interface{}) {},				casted: false,},
		{ function: NextFunc(func(item interface{}) {}),	casted: false,},
	}

	for _, singleTestData := range testData {

		// when
		result, ok := AsDoneFunc(singleTestData.function)

		// then
		assert.IsType(t, DoneFunc(nil), result)
		assert.Equal(t, singleTestData.casted, ok)
	}
}

func TestReturnsAsNextFunc(t *testing.T) {
	// where
	type data struct{
		function interface{}
		casted bool
	}
	testData := [...]data {
		{ function: func(interface{}) {},			casted: true,},
		{ function: NextFunc(func(interface{}) {}),	casted: true,},
		{ function: func(err error) {},				casted: false,},
		{ function: func() {},						casted: false,},
		{ function: DoneFunc(func() {}),			casted: false,},
	}

	for _, singleTestData := range testData {

		// when
		result, ok := AsNextFunc(singleTestData.function)

		// then
		assert.IsType(t, NextFunc(nil), result)
		assert.Equal(t, singleTestData.casted, ok)
	}
}

func TestReturnsAsErrFunc(t *testing.T) {
	// where
	type data struct{
		function interface{}
		casted bool
	}
	testData := [...]data {
		{ function: func(error) {},				casted: true,},
		{ function: ErrFunc(func(error) {}),	casted: true,},
		{ function: func(interface{}) {},		casted: false,},
		{ function: func() {},					casted: false,},
		{ function: DoneFunc(func() {}),		casted: false,},
	}

	for _, singleTestData := range testData {

		// when
		result, ok := AsErrFunc(singleTestData.function)

		// then
		assert.IsType(t, ErrFunc(nil), result)
		assert.Equal(t, singleTestData.casted, ok)
	}
}
