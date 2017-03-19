package iterable

import (
	"errors"
	"strings"
	"testing"

	"github.com/reactivex/rxgo"
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

func TestCreateSliceIterable(t *testing.T) {
	items := []int{}

	for i := 0; i < 10; i++ {
		items = append(items, i)
	}

	it, err := New(items)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it)

	for i := 0; i < 10; i++ {
		if v, err := it.Next(); err == nil {
			assert.Equal(i, v)
		} else {
			t.Fail()
		}
	}
}

func TestCreateMapIterable(t *testing.T) {
	items := map[string][]string{
		"a": []string{"aa", "bb"},
		"b": []string{"ba", "bb"},
		"c": []string{"ca", "cb"},
	}

	it, err := New(items)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it)

	for i := 0; i < len(items); i++ {
		if elem, err := it.Next(); err == nil {
			vs, _ := elem.([]interface{})
			k := vs[0].(string)
			v := vs[1].([]string)
			val := items[k]
			assert.Equal(len(val), len(v))
			assert.Equal(strings.Join(val, ""), strings.Join(v, ""))
		} else {
			t.Fail()
		}
	}
}

func TestCreateChanIntIterable(t *testing.T) {
	ch := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	it, err := New(ch)
	if err != nil {
		t.Fail()
	}

	assert := assert.New(t)
	assert.IsType(Iterable(nil), it)

	for i := 0; i < 10; i++ {
		if v, err := it.Next(); err == nil {
			assert.Equal(i, v)
		} else {
			t.Fail()
		}
	}
}
