package rxgo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ForEach(t *testing.T) {
	count := 0
	var gotErr error
	next := make(chan Item)
	done := make(chan struct{})

	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		next <- FromError(fooErr)
	}()

	obs := FromChannel(next)
	obs.ForEach(context.Background(), func(i interface{}) {
		count += i.(int)
	}, func(err error) {
		gotErr = err
		done <- struct{}{}
	}, func() {})

	// We avoid using the assertion API on purpose
	<-done
	assert.Equal(t, 6, count)
	assert.Equal(t, fooErr, gotErr)
}

func Test_Map_One(t *testing.T) {
	next := make(chan Item)
	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		close(next)
	}()

	obs := FromChannel(next).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	assertObservable(t, context.Background(), obs, hasItems(2, 3, 4), hasNotRaisedError())
}

func Test_Map_Multiple(t *testing.T) {
	next := make(chan Item)
	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		close(next)
	}()

	obs := FromChannel(next).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
	assertObservable(t, context.Background(), obs, hasItems(20, 30, 40), hasNotRaisedError())
}

func Test_Map_Error(t *testing.T) {
	next := make(chan Item)
	go func() {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		next <- FromError(fooErr)
	}()

	obs := FromChannel(next).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	assertObservable(t, context.Background(), obs, hasItems(2, 3, 4), hasRaisedError(fooErr))
}

func Test_Map_Cancel(t *testing.T) {
	next := make(chan Item)

	ctx, cancel := context.WithCancel(context.Background())
	obs := FromChannel(next).Map(ctx, func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	cancel()
	assertObservable(t, context.Background(), obs, hasNoItems(), hasNotRaisedError())
}
