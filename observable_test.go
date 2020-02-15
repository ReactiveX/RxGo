package rxgo

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ForEach(t *testing.T) {
	count := 0
	var gotErr error

	next := make(chan Item)

	ctx, cancel := context.WithCancel(context.Background())
	expectedErr := errors.New("foo")
	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		next <- FromError(expectedErr)
		cancel()
	}()

	obs := FromChannel(ctx, next)
	obs.ForEach(ctx, func(i interface{}) {
		count += i.(int)
	}, func(err error) {
		gotErr = err
	}, func() {})

	// We avoid using the assertion API on purpose
	<-ctx.Done()
	assert.Equal(t, 6, count)
	assert.Equal(t, expectedErr, gotErr)
}

func Test_Map(t *testing.T) {

}
