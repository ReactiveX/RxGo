package rxgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Observable_All(t *testing.T) {
	predicateAllInt := func(i interface{}) bool {
		switch i.(type) {
		case int:
			return true
		default:
			return false
		}
	}

	AssertSingle(context.Background(), t, FromItems(FromValue(1), FromValue(2), FromValue(3)).All(context.Background(), predicateAllInt),
		HasItem(true), HasNotRaisedError())
	AssertSingle(context.Background(), t, FromItems(FromValue(1), FromValue("x"), FromValue(3)).All(context.Background(), predicateAllInt),
		HasItem(false), HasNotRaisedError())
}

func Test_Observable_AverageFloat32(t *testing.T) {
	AssertSingle(context.Background(), t, FromItems(FromValue(float32(1)), FromValue(float32(2)), FromValue(float32(3))).AverageFloat32(context.Background()), HasItem(float32(2)))
	AssertSingle(context.Background(), t, FromItems(FromValue(float32(1)), FromValue(float32(20))).AverageFloat32(context.Background()), HasItem(float32(10.5)))
	AssertSingle(context.Background(), t, Empty().AverageFloat32(context.Background()), HasItem(0))
	AssertSingle(context.Background(), t, FromItems(FromValue("x")).AverageFloat32(context.Background()), HasRaisedAnError())
}

func Test_Observable_Filter(t *testing.T) {
	obs := FromChannel(channelValue(1, 2, 3, 4, closeCmd)).Filter(context.Background(),
		func(i interface{}) bool {
			return i.(int)%2 == 0
		})
	AssertObservable(context.Background(), t, obs, HasItems(2, 4), HasNotRaisedError())
}

func Test_Observable_ForEach(t *testing.T) {
	count := 0
	var gotErr error
	done := make(chan struct{})
	next := channelValue(1, 2, 3, errFoo)

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
	assert.Equal(t, errFoo, gotErr)
}

func Test_Observable_Map_One(t *testing.T) {
	next := channelValue(1, 2, 3, closeCmd)

	obs := FromChannel(next).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	AssertObservable(context.Background(), t, obs, HasItems(2, 3, 4), HasNotRaisedError())
}

func Test_Observable_Map_Multiple(t *testing.T) {
	next := channelValue(1, 2, 3, closeCmd)

	obs := FromChannel(next).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
	AssertObservable(context.Background(), t, obs, HasItems(20, 30, 40), HasNotRaisedError())
}

func Test_Observable_Map_Error(t *testing.T) {
	next := channelValue(1, 2, 3, errFoo)

	obs := FromChannel(next).Map(context.Background(), func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	AssertObservable(context.Background(), t, obs, HasItems(2, 3, 4), HasRaisedError(errFoo))
}

func Test_Observable_Map_Cancel(t *testing.T) {
	next := make(chan Item)

	ctx, cancel := context.WithCancel(context.Background())
	obs := FromChannel(next).Map(ctx, func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	cancel()
	AssertObservable(context.Background(), t, obs, HasNoItems(), HasNotRaisedError())
}

func Test_Observable_Observe(t *testing.T) {
	got := make([]int, 0)
	ch := FromChannel(channelValue(1, 2, 3, closeCmd)).Observe()
	for item := range ch {
		got = append(got, item.Value.(int))
	}
	assert.Equal(t, []int{1, 2, 3}, got)
}

func Test_Observable_SkipWhile(t *testing.T) {
	next := channelValue(1, 2, 3, 4, 5, closeCmd)

	obs := FromChannel(next).SkipWhile(context.Background(), func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i != 3
		default:
			return true
		}
	})

	AssertObservable(context.Background(), t, obs, HasItems(3, 4, 5), HasNotRaisedError())
}
