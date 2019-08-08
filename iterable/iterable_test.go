package iterable

import (
	"errors"
	"testing"

	"rxgo"
	"github.com/stretchr/testify/assert"
)

func TestIterableImplementsIterator(t *testing.T) {
	assert.Implements(t, (*rx.Iterator)(nil), Iterable(nil))
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

	it1, err := New(ch)
	if err != nil {
		t.Fail()
	}

	it2, err := New(items)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it1)
	assert.IsType(Iterable(nil), it2)

	for i := 0; i < 10; i++ {
		if v, err := it1.Next(); err == nil {
			assert.Equal(i, v)
		} else {
			t.Fail()
		}

		if v, err := it2.Next(); err == nil {
			assert.Equal(i, v)
		} else {
			t.Fail()
		}
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

	it1, err := New(ch)
	if err != nil {
		t.Fail()
	}

	it2, err := New(items)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it1)
	assert.IsType(Iterable(nil), it2)

	for _, item := range items {
		if v, err := it1.Next(); err == nil {
			assert.Equal(item, v)
		} else {
			t.Fail()
		}

		if v, err := it2.Next(); err == nil {
			assert.Equal(item, v)
		} else {
			t.Fail()
		}
	}
}
