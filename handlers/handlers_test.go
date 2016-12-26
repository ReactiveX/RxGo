package handlers

import (
	"errors"
	"testing"

	. "github.com/jochasinga/grx/bases"
	"github.com/stretchr/testify/assert"
)

type (
	myInt    int
	myString string
	myFloat  float64
	myRune   rune
)

func (n myInt) Emit() (Item, error) {
	return Item(n), nil
}
func (w myString) Emit() (Item, error) {
	return Item(w), nil
}
func (pi myFloat) Emit() (Item, error) {
	return Item(pi), nil
}
func (r myRune) Emit() (Item, error) {
	return Item(r), nil
}

func TestAllTypesImplementEventHandler(t *testing.T) {
	assert := assert.New(t)
	assert.Implements((*EventHandler)(nil), (*NextFunc)(nil))
	assert.Implements((*EventHandler)(nil), (*ErrFunc)(nil))
	assert.Implements((*EventHandler)(nil), (*DoneFunc)(nil))
}

func TestNextFuncHandleMethod(t *testing.T) {
	var (
		num  myInt    = 1
		word myString = "Hello"
		pi   myFloat  = 3.1416
		ru   myRune   = 'ด'
	)

	assert := assert.New(t)
	samples := []Item{}
	nextf := NextFunc(func(i Item) {
		samples = append(samples, i)
	})

	nextHandleTests := []struct {
		in       Emitter
		expected Item
	}{
		{num, Item(num)},
		{word, Item(word)},
		{pi, Item(pi)},
		{ru, Item(ru)},
	}

	for n, tt := range nextHandleTests {
		nextf.Handle(tt.in)
		assert.Equal(tt.expected, samples[n])
	}
}

type (
	someInt    int
	someString string
	someFloat  float64
	someRune   rune
)

func (n someInt) Emit() (Item, error) {
	return nil, errors.New("myInt error")
}
func (w someString) Emit() (Item, error) {
	return nil, errors.New("myString error")
}
func (pi someFloat) Emit() (Item, error) {
	return nil, errors.New("myFloat error")
}
func (r someRune) Emit() (Item, error) {
	return nil, errors.New("myRune error")
}

func TestErrFuncHandleMethod(t *testing.T) {
	var (
		num  someInt    = 1
		word someString = "Hello"
		pi   someFloat  = 3.1416
		ru   someRune   = 'ด'
	)

	assert := assert.New(t)
	samples := []error{}
	errf := ErrFunc(func(err error) {
		samples = append(samples, err)
	})

	errHandleTests := []struct {
		in       Emitter
		expected string
	}{
		{num, "myInt error"},
		{word, "myString error"},
		{pi, "myFloat error"},
		{ru, "myRune error"},
	}

	for n, tt := range errHandleTests {
		errf.Handle(tt.in)
		assert.Equal(tt.expected, samples[n].Error())
	}
}

func TestDoneFuncHandleMethod(t *testing.T) {
	var (
		num  someInt    = 1
		word someString = "Hello"
		pi   someFloat  = 3.1416
		ru   someRune   = 'ด'
	)

	assert := assert.New(t)
	samples := []string{}
	errf := DoneFunc(func() {
		samples = append(samples, "done")
	})

	doneHandleTests := []struct {
		in       Emitter
		expected string
	}{
		{num, "done"},
		{word, "done"},
		{pi, "done"},
		{ru, "done"},
	}

	for n, tt := range doneHandleTests {
		errf.Handle(tt.in)
		assert.Equal(tt.expected, samples[n])
	}
}
