package rxgo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Item_SendBlocking(t *testing.T) {
	ch := make(chan Item, 1)
	defer close(ch)
	Of(5).SendBlocking(ch)
	assert.Equal(t, 5, (<-ch).V)
}

func Test_Item_SendContext_True(t *testing.T) {
	ch := make(chan Item, 1)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.True(t, Of(5).SendContext(ctx, ch))
}

func Test_Item_SendNonBlocking(t *testing.T) {
	ch := make(chan Item, 1)
	defer close(ch)
	assert.True(t, Of(5).SendNonBlocking(ch))
	assert.False(t, Of(5).SendNonBlocking(ch))
}
