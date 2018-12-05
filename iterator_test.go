package rxgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorFromChannel(t *testing.T) {
	ch := make(chan interface{}, 1)
	it := NewIteratorFromChannel(ch)

	ch <- 1
	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())

	ch <- 2
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())

	close(ch)
	assert.False(t, it.Next())
}
