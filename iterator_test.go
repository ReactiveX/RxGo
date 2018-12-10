package rxgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorFromChannel(t *testing.T) {
	ch := make(chan interface{}, 1)
	it := newIteratorFromChannel(ch)

	ch <- 1
	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())

	ch <- 2
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())

	close(ch)
	assert.False(t, it.Next())
}

func TestIteratorFromSlice(t *testing.T) {
	it := newIteratorFromSlice([]interface{}{1, 2, 3})

	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())

	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())

	assert.True(t, it.Next())
	assert.Equal(t, 3, it.Value())

	assert.False(t, it.Next())
}

func TestName(t *testing.T) {
	just := Just(1).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	}).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	})

	AssertThatObservable(t, just, HasItems(3))
	AssertThatObservable(t, just, HasItems(3))
}
