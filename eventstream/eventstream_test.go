package eventstream

import (
	"testing"

	. "github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/errors"
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

var (
	num  myInt    = 1
	word myString = "Hello"
	pi   myFloat  = 3.1416
	ru   myRune   = 'ด'
)

type List struct {
	inner   []Emitter
	counter int
}

func (li *List) Next() (Emitter, error) {
	if li.counter < len(li.inner) {
		e := li.inner[li.counter]
		li.counter++
		return e, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

func TestEventStreamImplementIterator(t *testing.T) {
	evs := New(num, word, pi, ru)
	assert.Implements(t, (*Iterator)(nil), evs)
}

func TestCreateEventStreamWithConstructor(t *testing.T) {
	evs := New(num, word, pi, ru)
	assert.IsType(t, EventStream(nil), evs)
}

func TestCreateEventStreamWithFrom(t *testing.T) {
	li := &List{
		inner: []Emitter{num, word, pi, ru},
	}
	evs := From(li)
	assert.IsType(t, EventStream(nil), evs)
}

func TestEventStreamNextMethod(t *testing.T) {
	assert := assert.New(t)
	li := &List{
		inner: []Emitter{num, word, pi, ru},
	}

	evsA := From(li)
	evsB := New(num, word, pi, ru)

	for _, expected := range li.inner {
		emitter, err := evsA.Next()
		if err != nil {
			assert.IsType(EventStreamError{}, err)
		}
		assert.Equal(expected, emitter)
	}

	for _, expected := range li.inner {
		emitter, err := evsB.Next()
		if err != nil {
			assert.IsType(EventStreamError{}, err)
		}
		assert.Equal(expected, emitter)
	}
}
