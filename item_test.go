package rxgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func Test_SendItems_Variadic(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, 1, 2, 3)
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 2, 3), HasNoError())
}

func Test_SendItems_VariadicWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, 1, errFoo, 3)
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 3), HasError(errFoo))
}

func Test_SendItems_Slice(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, []int{1, 2, 3})
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 2, 3), HasNoError())
}

func Test_SendItems_SliceWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, []interface{}{1, errFoo, 3})
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 3), HasError(errFoo))
}

func Test_Item_SendBlocking(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	Of(5).SendBlocking(ch)
	assert.Equal(t, 5, (<-ch).V)
}

func Test_Item_SendContext_True(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.True(t, Of(5).SendContext(ctx, ch))
}

func Test_Item_SendNonBlocking(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	assert.True(t, Of(5).SendNonBlocking(ch))
	assert.False(t, Of(5).SendNonBlocking(ch))
}
