package rxgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func Test_OptionalSingle_Get_Item(t *testing.T) {
	defer goleak.VerifyNone(t)
	var os OptionalSingle = &OptionalSingleImpl{iterable: Just(1)()}
	get, err := os.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, get.V)
}

func Test_OptionalSingle_Get_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	var os OptionalSingle = &OptionalSingleImpl{iterable: Empty()}
	get, err := os.Get()
	assert.NoError(t, err)
	assert.Equal(t, OptionalSingleEmpty, get)
}

func Test_OptionalSingle_Get_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	var os OptionalSingle = &OptionalSingleImpl{iterable: Just(errFoo)()}
	get, err := os.Get()
	assert.NoError(t, err)
	assert.Equal(t, errFoo, get.E)
}

func Test_OptionalSingle_Get_ContextCanceled(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	var os OptionalSingle = &OptionalSingleImpl{iterable: Never()}
	cancel()
	_, err := os.Get(WithContext(ctx))
	assert.Equal(t, ctx.Err(), err)
}

func Test_OptionalSingle_Map(t *testing.T) {
	defer goleak.VerifyNone(t)
	single := Just(1)().Max(func(_, _ interface{}) int {
		return 1
	}).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, single, HasItem(2), HasNoError())
}

func Test_OptionalSingle_Observe(t *testing.T) {
	defer goleak.VerifyNone(t)
	os := JustItem(1).Filter(func(i interface{}) bool {
		return i == 1
	})
	Assert(context.Background(), t, os, HasItem(1), HasNoError())
}
