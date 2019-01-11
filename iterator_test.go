package rxgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIteratorFromChannel(t *testing.T) {
	ch := make(chan interface{}, 1)
	it := newIteratorFromChannel(ch)

	ch <- 1
	next, err := it.Next()
	assert.Nil(t, err)
	assert.Equal(t, 1, next)

	ch <- 2
	next, err = it.Next()
	assert.Nil(t, err)
	assert.Equal(t, 2, next)

	close(ch)
	_, err = it.Next()
	assert.NotNil(t, err)
}

func TestIteratorFromSlice(t *testing.T) {
	it := newIteratorFromSlice([]interface{}{1, 2, 3})

	next, err := it.Next()
	assert.Nil(t, err)
	assert.Equal(t, 1, next)

	next, err = it.Next()
	assert.Nil(t, err)
	assert.Equal(t, 2, next)

	next, err = it.Next()
	assert.Nil(t, err)
	assert.Equal(t, 3, next)

	_, err = it.Next()
	assert.NotNil(t, err)
}
