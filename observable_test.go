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

	next := make(chan interface{})
	errs := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	expectedErr := errors.New("foo")
	go func() {
		next <- 1
		next <- 2
		next <- 3
		errs <- expectedErr
		cancel()
	}()

	obs := FromChannel(ctx, next, errs)
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
