package iterable

import (
	"errors"
	"testing"

	"github.com/jochasinga/grx/bases"
	"github.com/stretchr/testify/assert"
)

func TestIterableImplementsIterator(t *testing.T) {
	assert.Implements(t, (*bases.Iterator)(nil), Iterable(nil))
}

func TestCreateHomogenousIterable(t *testing.T) {
	ch := make(chan interface{})
	items := []interface{}{}

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for i := 0; i < 10; i++ {
		items = append(items, i)
	}

	it1, err := From(ch)
	if err != nil {
		t.Fail()
	}

	it2, err := From(items)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it1)
	assert.IsType(Iterable(nil), it2)

	for i := 0; i < 10; i++ {
		assert.Equal(i, it1.Next())
		assert.Equal(i, it2.Next())
	}
}

func TestCreateHeterogeneousIterable(t *testing.T) {
	ch := make(chan interface{})
	items := []interface{}{
		"foo", "bar", "baz", 'a', 'b', errors.New("bang"), 99,
	}

	go func() {
		for _, item := range items {
			ch <- item
		}
		close(ch)
	}()

	it1, err := From(ch)
	if err != nil {
		t.Fail()
	}

	it2, err := From(items)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it1)
	assert.IsType(Iterable(nil), it2)

	for _, item := range items {
		assert.Equal(item, it1.Next())
		assert.Equal(item, it2.Next())
	}
}
