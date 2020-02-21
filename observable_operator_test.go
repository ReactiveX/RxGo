package rxgo

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/stretchr/testify/assert"
)

var predicateAllInt = func(i interface{}) bool {
	switch i.(type) {
	case int:
		return true
	default:
		return false
	}
}

func Test_Observable_All_True(t *testing.T) {
	Assert(context.Background(), t, Range(1, 10000).All(predicateAllInt),
		HasItem(true), HasNotRaisedError())
}

func Test_Observable_All_False(t *testing.T) {
	Assert(context.Background(), t, testObservable(1, "x", 3).All(predicateAllInt),
		HasItem(false), HasNotRaisedError())
}

func Test_Observable_All_Parallel_True(t *testing.T) {
	Assert(context.Background(), t, Range(1, 10000).All(predicateAllInt, WithCPUPool()),
		HasItem(true), HasNotRaisedError())
}

func Test_Observable_All_Parallel_False(t *testing.T) {
	Assert(context.Background(), t, testObservable(1, "x", 3).All(predicateAllInt, WithCPUPool()),
		HasItem(false), HasNotRaisedError())
}

func Test_Observable_All_Parallel_Error(t *testing.T) {
	Assert(context.Background(), t, testObservable(1, errFoo, 3).All(predicateAllInt, WithCPUPool()),
		HasNoItem(), HasRaisedError(errFoo))
}

func Test_Observable_AverageFloat32(t *testing.T) {
	Assert(context.Background(), t, testObservable(float32(1), float32(20)).AverageFloat32(), HasItem(float32(10.5)))
}

func Test_Observable_AverageFloat32_Empty(t *testing.T) {
	Assert(context.Background(), t, Empty().AverageFloat32(), HasItem(0))
}

func Test_Observable_AverageFloat32_Error(t *testing.T) {
	Assert(context.Background(), t, testObservable("x").AverageFloat32(), HasRaisedAnError())
}

func Test_Observable_AverageFloat32_Parallel(t *testing.T) {
	Assert(context.Background(), t, testObservable(float32(1), float32(20)).AverageFloat32(), HasItem(float32(10.5)))
}

func Test_Observable_AverageFloat32_Parallel_Empty(t *testing.T) {
	Assert(context.Background(), t, Empty().AverageFloat32(WithCPUPool()),
		HasItem(0))
}

func Test_Observable_AverageFloat32_Parallel_Error(t *testing.T) {
	Assert(context.Background(), t, testObservable("x").AverageFloat32(WithCPUPool()),
		HasRaisedAnError())
}

func Test_Observable_AverageFloat64(t *testing.T) {
	Assert(context.Background(), t, testObservable(float64(1), float64(2), float64(3)).AverageFloat64(), HasItem(float64(2)))
	Assert(context.Background(), t, testObservable(float64(1), float64(20)).AverageFloat64(), HasItem(10.5))
	Assert(context.Background(), t, Empty().AverageFloat64(), HasItem(0))
	Assert(context.Background(), t, testObservable("x").AverageFloat64(), HasRaisedAnError())
}

func Test_Observable_AverageInt(t *testing.T) {
	Assert(context.Background(), t, testObservable(1, 2, 3).AverageInt(), HasItem(2))
	Assert(context.Background(), t, testObservable(1, 20).AverageInt(), HasItem(10))
	Assert(context.Background(), t, Empty().AverageInt(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt(), HasRaisedAnError())
}

func Test_Observable_AverageInt8(t *testing.T) {
	Assert(context.Background(), t, testObservable(int8(1), int8(2), int8(3)).AverageInt8(), HasItem(int8(2)))
	Assert(context.Background(), t, testObservable(int8(1), int8(20)).AverageInt8(), HasItem(int8(10)))
	Assert(context.Background(), t, Empty().AverageInt8(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt8(), HasRaisedAnError())
}

func Test_Observable_AverageInt16(t *testing.T) {
	Assert(context.Background(), t, testObservable(int16(1), int16(2), int16(3)).AverageInt16(), HasItem(int16(2)))
	Assert(context.Background(), t, testObservable(int16(1), int16(20)).AverageInt16(), HasItem(int16(10)))
	Assert(context.Background(), t, Empty().AverageInt16(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt16(), HasRaisedAnError())
}

func Test_Observable_AverageInt32(t *testing.T) {
	Assert(context.Background(), t, testObservable(int32(1), int32(2), int32(3)).AverageInt32(), HasItem(int32(2)))
	Assert(context.Background(), t, testObservable(int32(1), int32(20)).AverageInt32(), HasItem(int32(10)))
	Assert(context.Background(), t, Empty().AverageInt32(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt32(), HasRaisedAnError())
}

func Test_Observable_AverageInt64(t *testing.T) {
	Assert(context.Background(), t, testObservable(int64(1), int64(2), int64(3)).AverageInt64(), HasItem(int64(2)))
	Assert(context.Background(), t, testObservable(int64(1), int64(20)).AverageInt64(), HasItem(int64(10)))
	Assert(context.Background(), t, Empty().AverageInt64(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt64(), HasRaisedAnError())
}

func Test_Observable_BackOffRetry(t *testing.T) {
	i := 0
	backOffCfg := backoff.NewExponentialBackOff()
	backOffCfg.InitialInterval = time.Nanosecond
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item, done func()) {
		next <- Of(1)
		next <- Of(2)
		if i == 2 {
			next <- Of(3)
			done()
		} else {
			i++
			next <- Error(errFoo)
			done()
		}
	}}).BackOffRetry(backoff.WithMaxRetries(backOffCfg, 3))
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 2, 1, 2, 3), HasNotRaisedError())
}

func Test_Observable_BackOffRetry_Error(t *testing.T) {
	backOffCfg := backoff.NewExponentialBackOff()
	backOffCfg.InitialInterval = time.Nanosecond
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item, done func()) {
		next <- Of(1)
		next <- Of(2)
		next <- Error(errFoo)
		done()
	}}).BackOffRetry(backoff.WithMaxRetries(backOffCfg, 3))
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 2, 1, 2, 1, 2), HasRaisedError(errFoo))
}

func Test_Observable_BufferWithCount_CountAndSkipEqual(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5, 6).BufferWithCount(3, 3)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4, 5, 6}))
}

func Test_Observable_BufferWithCount_CountAndSkipNotEqual(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5, 6).BufferWithCount(2, 3)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2}, []interface{}{4, 5}))
}

func Test_Observable_BufferWithCount_IncompleteLastItem(t *testing.T) {
	obs := testObservable(1, 2, 3, 4).BufferWithCount(2, 3)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2}, []interface{}{4}))
}

func Test_Observable_BufferWithCount_Error(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, errFoo).BufferWithCount(3, 3)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4}), HasRaisedError(errFoo))
}

func Test_Observable_BufferWithCount_InvalidInputs(t *testing.T) {
	obs := testObservable(1, 2, 3, 4).BufferWithCount(0, 5)
	Assert(context.Background(), t, obs, HasRaisedAnError())

	obs = testObservable(1, 2, 3, 4).BufferWithCount(5, 0)
	Assert(context.Background(), t, obs, HasRaisedAnError())
}

func Test_Observable_BufferWithTime_MockedTime(t *testing.T) {
	timespan := new(mockDuration)
	timespan.On("duration").Return(10 * time.Second)

	timeshift := new(mockDuration)
	timeshift.On("duration").Return(10 * time.Second)

	obs := testObservable(1, 2, 3).BufferWithTime(timespan, timeshift)

	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2, 3}))
	timespan.AssertCalled(t, "duration")
	timeshift.AssertNotCalled(t, "duration")
}

func Test_Observable_BufferWithTime_MinorMockedTime(t *testing.T) {
	ch := make(chan Item)
	from := FromChannel(ch)

	timespan := new(mockDuration)
	timespan.On("duration").Return(1 * time.Millisecond)

	timeshift := new(mockDuration)
	timeshift.On("duration").Return(1 * time.Millisecond)

	obs := from.BufferWithTime(timespan, timeshift)

	ch <- Of(1)
	close(ch)

	<-obs.Observe()
	timespan.AssertCalled(t, "duration")
}

func Test_Observable_BufferWithTime_IllegalInput(t *testing.T) {
	Assert(context.Background(), t, Empty().BufferWithTime(nil, nil), HasRaisedAnError())
	Assert(context.Background(), t, Empty().BufferWithTime(WithDuration(0*time.Second), nil), HasRaisedAnError())
}

func Test_Observable_BufferWithTime_NilTimeshift(t *testing.T) {
	testObservable := testObservable(1, 2, 3)
	obs := testObservable.BufferWithTime(WithDuration(1*time.Second), nil)
	Assert(context.Background(), t, obs, HasSomeItems())
}

func Test_Observable_BufferWithTime_Error(t *testing.T) {
	testObservable := testObservable(1, 2, 3, errFoo)
	obs := testObservable.BufferWithTime(WithDuration(1*time.Second), nil)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2, 3}), HasRaisedError(errFoo))
}

func Test_Observable_BufferWithTimeOrCount_InvalidInputs(t *testing.T) {
	obs := Empty().BufferWithTimeOrCount(nil, 5)
	Assert(context.Background(), t, obs, HasRaisedAnError())

	obs = Empty().BufferWithTimeOrCount(WithDuration(0), 5)
	Assert(context.Background(), t, obs, HasRaisedAnError())

	obs = Empty().BufferWithTimeOrCount(WithDuration(time.Millisecond*5), 0)
	Assert(context.Background(), t, obs, HasRaisedAnError())
}

func Test_Observable_BufferWithTimeOrCount_Count(t *testing.T) {
	testObservable := testObservable(1, 2, 3)
	obs := testObservable.BufferWithTimeOrCount(WithDuration(1*time.Second), 2)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2}, []interface{}{3}))
}

func Test_Observable_BufferWithTimeOrCount_MockedTime(t *testing.T) {
	ch := make(chan Item)
	from := FromChannel(ch)

	timespan := new(mockDuration)
	timespan.On("duration").Return(1 * time.Millisecond)

	obs := from.BufferWithTimeOrCount(timespan, 5)

	time.Sleep(50 * time.Millisecond)
	ch <- Of(1)
	close(ch)

	<-obs.Observe()
	timespan.AssertCalled(t, "duration")
}

func Test_Observable_BufferWithTimeOrCount_Error(t *testing.T) {
	testObservable := testObservable(1, 2, 3, errFoo, 4)
	obs := testObservable.BufferWithTimeOrCount(WithDuration(10*time.Second), 2)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2}, []interface{}{3}),
		HasRaisedError(errFoo))
}

func Test_Observable_Contain(t *testing.T) {
	predicate := func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i == 2
		default:
			return false
		}
	}

	Assert(context.Background(), t,
		testObservable(1, 2, 3).Contains(predicate),
		HasItem(true))
	Assert(context.Background(), t,
		testObservable(1, 5, 3).Contains(predicate),
		HasItem(false))
}

func Test_Observable_Count(t *testing.T) {
	single := testObservable(1, 2, 3, "foo", "bar", errFoo).Count()
	Assert(context.Background(), t, single, HasItem(int64(6)))
}

func Test_Observable_DefaultIfEmpty_Empty(t *testing.T) {
	obs := Empty().DefaultIfEmpty(3)
	Assert(context.Background(), t, obs, HasItems(3))
}

func Test_Observable_DefaultIfEmpty_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2).DefaultIfEmpty(3)
	Assert(context.Background(), t, obs, HasItems(1, 2))
}

func Test_Observable_Distinct(t *testing.T) {
	obs := testObservable(1, 2, 2, 1, 3).Distinct(func(item interface{}) (interface{}, error) {
		return item, nil
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_Observable_Distinct_Error(t *testing.T) {
	obs := testObservable(1, 2, 2, errFoo, 3).Distinct(func(item interface{}) (interface{}, error) {
		return item, nil
	})
	Assert(context.Background(), t, obs, HasItems(1, 2), HasRaisedError(errFoo))
}

func Test_Observable_Distinct_Error2(t *testing.T) {
	obs := testObservable(1, 2, 2, 2, 3, 4).Distinct(func(item interface{}) (interface{}, error) {
		if item.(int) == 3 {
			return nil, errFoo
		}
		return item, nil
	})
	Assert(context.Background(), t, obs, HasItems(1, 2), HasRaisedError(errFoo))
}

func Test_Observable_DistinctUntilChanged(t *testing.T) {
	obs := testObservable(1, 2, 2, 1, 3).DistinctUntilChanged(func(item interface{}) (interface{}, error) {
		return item, nil
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 3))
}

func Test_Observable_DoOnCompleted_NoError(t *testing.T) {
	called := false
	<-testObservable(1, 2, 3).DoOnCompleted(func() {
		called = true
	})
	assert.True(t, called)
}

func Test_Observable_DoOnCompleted_Error(t *testing.T) {
	called := false
	<-testObservable(1, errFoo, 3).DoOnCompleted(func() {
		called = true
	})
	assert.True(t, called)
}

func Test_Observable_DoOnError_NoError(t *testing.T) {
	var got error
	<-testObservable(1, 2, 3).DoOnError(func(err error) {
		got = err
	})
	assert.Nil(t, got)
}

func Test_Observable_DoOnError_Error(t *testing.T) {
	var got error
	<-testObservable(1, errFoo, 3).DoOnError(func(err error) {
		got = err
	})
	assert.Equal(t, errFoo, got)
}

func Test_Observable_DoOnNext_NoError(t *testing.T) {
	s := make([]interface{}, 0)
	<-testObservable(1, 2, 3).DoOnNext(func(i interface{}) {
		s = append(s, i)
	})
	assert.Equal(t, []interface{}{1, 2, 3}, s)
}

func Test_Observable_DoOnNext_Error(t *testing.T) {
	s := make([]interface{}, 0)
	<-testObservable(1, errFoo, 3).DoOnNext(func(i interface{}) {
		s = append(s, i)
	})
	assert.Equal(t, []interface{}{1}, s)
}

func Test_Observable_ElementAt(t *testing.T) {
	obs := testObservable(0, 1, 2, 3, 4).ElementAt(2)
	Assert(context.Background(), t, obs, HasItems(2))
}

func Test_Observable_ElementAt_Error(t *testing.T) {
	obs := testObservable(0, 1, 2, 3, 4).ElementAt(10)
	Assert(context.Background(), t, obs, HasNoItem(), HasRaisedAnError())
}

func Test_Observable_Error_NoError(t *testing.T) {
	assert.NoError(t, testObservable(1, 2, 3).Error())
}

func Test_Observable_Error_Error(t *testing.T) {
	assert.Equal(t, errFoo, testObservable(1, errFoo, 3).Error())
}

func Test_Observable_Errors_NoError(t *testing.T) {
	assert.Equal(t, 0, len(testObservable(1, 2, 3).Errors()))
}

func Test_Observable_Errors_OneError(t *testing.T) {
	assert.Equal(t, 1, len(testObservable(1, errFoo, 3).Errors()))
}

func Test_Observable_Errors_MultipleError(t *testing.T) {
	assert.Equal(t, 2, len(testObservable(1, errFoo, errBar).Errors()))
}

func Test_Observable_Errors_MultipleErrorFromMap(t *testing.T) {
	errs := testObservable(1, 2, 3, 4).Map(func(i interface{}) (interface{}, error) {
		if i == 2 {
			return nil, errFoo
		}
		if i == 3 {
			return nil, errBar
		}
		return i, nil
	}, WithErrorStrategy(Continue)).Errors()
	assert.Equal(t, 2, len(errs))
}

func Test_Observable_Filter(t *testing.T) {
	obs := testObservable(1, 2, 3, 4).Filter(
		func(i interface{}) bool {
			return i.(int)%2 == 0
		})
	Assert(context.Background(), t, obs, HasItems(2, 4), HasNotRaisedError())
}

func Test_Observable_First_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).First()
	Assert(context.Background(), t, obs, HasItem(1))
}

func Test_Observable_First_Empty(t *testing.T) {
	obs := Empty().First()
	Assert(context.Background(), t, obs, HasNoItem())
}

func Test_Observable_FirstOrDefault_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).FirstOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(1))
}

func Test_Observable_FirstOrDefault_Empty(t *testing.T) {
	obs := Empty().FirstOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(10))
}

func Test_Observable_FlatMap(t *testing.T) {
	obs := testObservable(1, 2, 3).FlatMap(func(i Item) Observable {
		return testObservable(i.V.(int)+1, i.V.(int)*10)
	})
	Assert(context.Background(), t, obs, HasItems(2, 10, 3, 20, 4, 30))
}

func Test_Observable_FlatMap_Error1(t *testing.T) {
	obs := testObservable(1, 2, 3).FlatMap(func(i Item) Observable {
		if i.V == 2 {
			return testObservable(errFoo)
		}
		return testObservable(i.V.(int)+1, i.V.(int)*10)
	})
	Assert(context.Background(), t, obs, HasItems(2, 10), HasRaisedError(errFoo))
}

func Test_Observable_FlatMap_Error2(t *testing.T) {
	obs := testObservable(1, errFoo, 3).FlatMap(func(i Item) Observable {
		if i.Error() {
			return testObservable(0)
		}
		return testObservable(i.V.(int)+1, i.V.(int)*10)
	})
	Assert(context.Background(), t, obs, HasItems(2, 10, 0, 4, 30), HasNotRaisedError())
}

func Test_Observable_ForEach_Error(t *testing.T) {
	count := 0
	var gotErr error
	done := make(chan struct{})

	obs := testObservable(1, 2, 3, errFoo)
	obs.ForEach(func(i interface{}) {
		count += i.(int)
	}, func(err error) {
		gotErr = err
		done <- struct{}{}
	}, func() {
		done <- struct{}{}
	})

	// We avoid using the assertion API on purpose
	<-done
	assert.Equal(t, 6, count)
	assert.Equal(t, errFoo, gotErr)
}

func Test_Observable_ForEach_Done(t *testing.T) {
	count := 0
	var gotErr error
	done := make(chan struct{})

	obs := testObservable(1, 2, 3)
	obs.ForEach(func(i interface{}) {
		count += i.(int)
	}, func(err error) {
		gotErr = err
		done <- struct{}{}
	}, func() {
		done <- struct{}{}
	})

	// We avoid using the assertion API on purpose
	<-done
	assert.Equal(t, 6, count)
	assert.Nil(t, gotErr)
}

func Test_Observable_IgnoreElements(t *testing.T) {
	obs := testObservable(1, 2, 3).IgnoreElements()
	Assert(context.Background(), t, obs, HasNoItem())
}

func Test_Observable_IgnoreElements_Error(t *testing.T) {
	obs := testObservable(1, errFoo, 3).IgnoreElements()
	Assert(context.Background(), t, obs, HasNoItem(), HasRaisedError(errFoo))
}

func Test_Observable_GroupBy(t *testing.T) {
	count := 3
	max := 10

	obs := Range(0, max).GroupBy(count, func(item Item) int {
		return item.V.(int) % count
	}, WithBufferedChannel(max))
	s, err := obs.ToSlice(0)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if len(s) != count {
		assert.FailNow(t, "length", "got=%d, expected=%d", len(s), count)
	}

	Assert(context.Background(), t, s[0].(Observable), HasItems(0, 3, 6, 9), HasNotRaisedError())
	Assert(context.Background(), t, s[1].(Observable), HasItems(1, 4, 7, 10), HasNotRaisedError())
	Assert(context.Background(), t, s[2].(Observable), HasItems(2, 5, 8), HasNotRaisedError())
}

func Test_Observable_GroupBy_Error(t *testing.T) {
	count := 3
	max := 10

	obs := Range(0, max).GroupBy(count, func(item Item) int {
		return 4
	}, WithBufferedChannel(max))
	s, err := obs.ToSlice(0)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if len(s) != count {
		assert.FailNow(t, "length", "got=%d, expected=%d", len(s), count)
	}

	Assert(context.Background(), t, s[0].(Observable), HasRaisedAnError())
	Assert(context.Background(), t, s[1].(Observable), HasRaisedAnError())
	Assert(context.Background(), t, s[2].(Observable), HasRaisedAnError())
}

func Test_Observable_Last_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).Last()
	Assert(context.Background(), t, obs, HasItem(3))
}

func Test_Observable_Last_Empty(t *testing.T) {
	obs := Empty().Last()
	Assert(context.Background(), t, obs, HasNoItem())
}

func Test_Observable_LastOrDefault_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).LastOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(3))
}

func Test_Observable_LastOrDefault_Empty(t *testing.T) {
	obs := Empty().LastOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(10))
}

func Test_Observable_Map_One(t *testing.T) {
	obs := testObservable(1, 2, 3).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, obs, HasItems(2, 3, 4), HasNotRaisedError())
}

func Test_Observable_Map_Multiple(t *testing.T) {
	obs := testObservable(1, 2, 3).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}).Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
	Assert(context.Background(), t, obs, HasItems(20, 30, 40), HasNotRaisedError())
}

func Test_Observable_Map_Error(t *testing.T) {
	obs := testObservable(1, 2, 3, errFoo).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, obs, HasItems(2, 3, 4), HasRaisedError(errFoo))
}

func Test_Observable_Map_ReturnValueAndError(t *testing.T) {
	obs := testObservable(1).Map(func(i interface{}) (interface{}, error) {
		return 2, errFoo
	})
	Assert(context.Background(), t, obs, HasNoItems(), HasRaisedError(errFoo))
}

func Test_Observable_Map_Multiple_Error(t *testing.T) {
	called := false
	obs := testObservable(1, 2, 3).Map(func(i interface{}) (interface{}, error) {
		return nil, errFoo
	}).Map(func(i interface{}) (interface{}, error) {
		called = true
		return nil, nil
	})
	Assert(context.Background(), t, obs, HasNoItems(), HasRaisedError(errFoo))
	assert.False(t, called)
}

func Test_Observable_Map_Cancel(t *testing.T) {
	next := make(chan Item)

	ctx, cancel := context.WithCancel(context.Background())
	obs := FromChannel(next).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}, WithContext(ctx))
	cancel()
	Assert(context.Background(), t, obs, HasNoItems(), HasNotRaisedError())
}

func Test_Observable_Map_Parallel(t *testing.T) {
	const len = 10
	ch := make(chan Item, len)
	go func() {
		for i := 0; i < len; i++ {
			ch <- Of(i)
		}
		close(ch)
	}()

	obs := FromChannel(ch).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}, WithPool(len))
	Assert(context.Background(), t, obs, HasItemsNoParticularOrder(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), HasNotRaisedError())
}

func Test_Observable_Marshal(t *testing.T) {
	obs := testObservable(testStruct{
		ID: 1,
	}, testStruct{
		ID: 2,
	}).Marshal(json.Marshal)
	Assert(context.Background(), t, obs, HasItems([]byte(`{"id":1}`), []byte(`{"id":2}`)))
}

func Test_Observable_Max(t *testing.T) {
	obs := testObservable(1, 5, 2, -1, 3).Max(func(e1 interface{}, e2 interface{}) int {
		i1 := e1.(int)
		i2 := e2.(int)
		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		} else {
			return 0
		}
	})
	Assert(context.Background(), t, obs, HasItem(5))
}

func Test_Observable_Min(t *testing.T) {
	obs := testObservable(1, 5, 2, -1, 3).Min(func(e1 interface{}, e2 interface{}) int {
		i1 := e1.(int)
		i2 := e2.(int)
		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		} else {
			return 0
		}
	})
	Assert(context.Background(), t, obs, HasItem(-1))
}

// TODO Fix race
// func Test_Observable_Min_Parallel(t *testing.T) {
//	obs := Range(1, 10000).Min(func(e1 interface{}, e2 interface{}) int {
//		i1 := e1.(int)
//		i2 := e2.(int)
//		if i1 > i2 {
//			return 1
//		} else if i1 < i2 {
//			return -1
//		} else {
//			return 0
//		}
//	}, WithPool(4))
//	Assert(context.Background(), t, obs, HasItem(-1))
//}

func Test_Observable_Observe(t *testing.T) {
	got := make([]int, 0)
	ch := testObservable(1, 2, 3).Observe()
	for item := range ch {
		got = append(got, item.V.(int))
	}
	assert.Equal(t, []int{1, 2, 3}, got)
}

func Test_Observable_OnErrorResumeNext(t *testing.T) {
	obs := testObservable(1, 2, errFoo, 4).OnErrorResumeNext(func(e error) Observable {
		return testObservable(10, 20)
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 10, 20), HasNotRaisedError())
}

func Test_Observable_OnErrorReturn(t *testing.T) {
	obs := testObservable(1, 2, errFoo, 4, errBar, 6).OnErrorReturn(func(err error) interface{} {
		return err.Error()
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, "foo", 4, "bar", 6), HasNotRaisedError())
}

func Test_Observable_OnErrorReturnItem(t *testing.T) {
	obs := testObservable(1, 2, errFoo, 4, errBar, 6).OnErrorReturnItem("foo")
	Assert(context.Background(), t, obs, HasItems(1, 2, "foo", 4, "foo", 6), HasNotRaisedError())
}

func Test_Observable_Reduce(t *testing.T) {
	obs := Range(1, 10000).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	})
	Assert(context.Background(), t, obs, HasItem(50015001), HasNotRaisedError())
}

func Test_Observable_Reduce_Parallel(t *testing.T) {
	obs := Range(1, 10000).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	}, WithCPUPool())
	Assert(context.Background(), t, obs, HasItem(50015001), HasNotRaisedError())
}

func Test_Observable_Reduce_Parallel_Error(t *testing.T) {
	obs := Range(1, 10000).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		if elem == 1000 {
			return nil, errFoo
		}
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	}, WithCPUPool())
	Assert(context.Background(), t, obs, HasRaisedError(errFoo))
}

func Test_Observable_Reduce_Parallel_WithErrorStrategy(t *testing.T) {
	obs := Range(1, 10000).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		if elem == 1 {
			return nil, errFoo
		}
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	}, WithCPUPool(), WithErrorStrategy(Continue))
	Assert(context.Background(), t, obs, HasItem(50015000), HasRaisedError(errFoo))
}

func Test_Observable_Reduce_Empty(t *testing.T) {
	obs := Empty().Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		return 0, nil
	})
	Assert(context.Background(), t, obs, HasNoItem(), HasNotRaisedError())
}

func Test_Observable_Reduce_Error(t *testing.T) {
	obs := testObservable(1, 2, errFoo, 4, 5).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		return 0, nil
	})
	Assert(context.Background(), t, obs, HasNoItem(), HasRaisedError(errFoo))
}

func Test_Observable_Repeat(t *testing.T) {
	repeat := testObservable(1, 2, 3).Repeat(1, nil)
	Assert(context.Background(), t, repeat, HasItems(1, 2, 3, 1, 2, 3))
}

func Test_Observable_Repeat_Zero(t *testing.T) {
	repeat := testObservable(1, 2, 3).Repeat(0, nil)
	Assert(context.Background(), t, repeat, HasItems(1, 2, 3))
}

func Test_Observable_Repeat_NegativeCount(t *testing.T) {
	repeat := testObservable(1, 2, 3).Repeat(-2, nil)
	Assert(context.Background(), t, repeat, HasNoItem(), HasRaisedAnError())
}

func Test_Observable_Repeat_Infinite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	repeat := testObservable(1, 2, 3).Repeat(Infinite, nil, WithContext(ctx))
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	Assert(context.Background(), t, repeat, HasNotRaisedError(), CustomPredicate(func(items []interface{}) error {
		if len(items) == 0 {
			return errors.New("no items")
		}
		return nil
	}))
}

func Test_Observable_Repeat_Frequency(t *testing.T) {
	frequency := new(mockDuration)
	frequency.On("duration").Return(time.Millisecond)

	repeat := testObservable(1, 2, 3).Repeat(1, frequency)
	Assert(context.Background(), t, repeat, HasItems(1, 2, 3, 1, 2, 3))
	frequency.AssertNumberOfCalls(t, "duration", 1)
	frequency.AssertExpectations(t)
}

func Test_Observable_ReturnError(t *testing.T) {
	obs := testObservable(1, 2, 3).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		if elem == 2 {
			return 0, errFoo
		}
		return elem, nil
	})
	Assert(context.Background(), t, obs, HasNoItem(), HasRaisedError(errFoo))
}

func Test_Observable_Retry(t *testing.T) {
	i := 0
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item, done func()) {
		next <- Of(1)
		next <- Of(2)
		if i == 2 {
			next <- Of(3)
			done()
		} else {
			i++
			next <- Error(errFoo)
			done()
		}
	}}).Retry(3)
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 2, 1, 2, 3), HasNotRaisedError())
}

func Test_Observable_Retry_Error(t *testing.T) {
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item, done func()) {
		next <- Of(1)
		next <- Of(2)
		next <- Error(errFoo)
		done()
	}}).Retry(3)
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 2, 1, 2, 1, 2), HasRaisedError(errFoo))
}

func Test_Observable_Run(t *testing.T) {
	s := make([]int, 0)
	<-testObservable(1, 2, 3).Map(func(i interface{}) (interface{}, error) {
		s = append(s, i.(int))
		return i, nil
	}).Run()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func Test_Observable_Run_Error(t *testing.T) {
	s := make([]int, 0)
	<-testObservable(1, errFoo).Map(func(i interface{}) (interface{}, error) {
		s = append(s, i.(int))
		return i, nil
	}).Run()
	assert.Equal(t, []int{1}, s)
}

func Test_Observable_Sample(t *testing.T) {
	obs := testObservable(1).Sample(Empty())
	Assert(context.Background(), t, obs, HasNoItem(), HasNotRaisedError())
}

func Test_Observable_Scan(t *testing.T) {
	obs := testObservable(0, 1, 3, 5, 1, 8).Scan(func(x interface{}, y interface{}) (interface{}, error) {
		var v1, v2 int

		if x, ok := x.(int); ok {
			v1 = x
		}

		if y, ok := y.(int); ok {
			v2 = y
		}

		return v1 + v2, nil
	})
	Assert(context.Background(), t, obs, HasItems(0, 1, 4, 9, 10, 18))
}

func Test_Observable_Send(t *testing.T) {
	ch := make(chan Item, 10)
	testObservable(1, 2, 3, errFoo).Send(ch)
	assert.Equal(t, Of(1), <-ch)
	assert.Equal(t, Of(2), <-ch)
	assert.Equal(t, Of(3), <-ch)
	assert.Equal(t, Error(errFoo), <-ch)
}

func Test_Observable_SequenceEqual_EvenSequence(t *testing.T) {
	sequence := testObservable(2, 5, 12, 43, 98, 100, 213)
	result := testObservable(2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequence)
	Assert(context.Background(), t, result, HasItem(true))
}

func Test_Observable_SequenceEqual_UnevenSequence(t *testing.T) {
	sequence := testObservable(2, 5, 12, 43, 98, 100, 213)
	result := testObservable(2, 5, 12, 43, 15, 100, 213).SequenceEqual(sequence)
	Assert(context.Background(), t, result, HasItem(false))
}

func Test_Observable_SequenceEqual_DifferentLengthSequence(t *testing.T) {
	sequenceShorter := testObservable(2, 5, 12, 43, 98, 100)
	sequenceLonger := testObservable(2, 5, 12, 43, 98, 100, 213, 512)

	resultForShorter := testObservable(2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequenceShorter)
	Assert(context.Background(), t, resultForShorter, HasItem(false))

	resultForLonger := testObservable(2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequenceLonger)
	Assert(context.Background(), t, resultForLonger, HasItem(false))
}

func Test_Observable_SequenceEqual_Empty(t *testing.T) {
	result := Empty().SequenceEqual(Empty())
	Assert(context.Background(), t, result, HasItem(true))
}

type message struct {
	id int
}

func Test_Observable_Serialize(t *testing.T) {
	obs := testObservable(message{3}, message{5}, message{1}, message{2}, message{4}).
		Serialize(1, func(i interface{}) int {
			return i.(message).id
		})
	Assert(context.Background(), t, obs, HasItems(message{1}, message{2}, message{3}, message{4}, message{5}))
}

func Test_Observable_Serialize_DifferentFrom(t *testing.T) {
	obs := testObservable(message{13}, message{15}, message{11}, message{12}, message{14}).
		Serialize(11, func(i interface{}) int {
			return i.(message).id
		})
	Assert(context.Background(), t, obs, HasItems(message{11}, message{12}, message{13}, message{14}, message{15}))
}

func Test_Observable_Serialize_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	obs := Never().Serialize(1, func(i interface{}) int {
		return i.(message).id
	}, WithContext(ctx))
	Assert(context.Background(), t, obs, HasNoItems(), HasNotRaisedError())
}

func Test_Observable_Serialize_Empty(t *testing.T) {
	obs := testObservable(message{3}, message{5}, message{7}, message{2}, message{4}).
		Serialize(1, func(i interface{}) int {
			return i.(message).id
		})
	Assert(context.Background(), t, obs, HasNoItems())
}

func Test_Observable_Serialize_Error(t *testing.T) {
	obs := testObservable(message{3}, message{1}, errFoo, message{2}, message{4}).
		Serialize(1, func(i interface{}) int {
			return i.(message).id
		})
	Assert(context.Background(), t, obs, HasItems(message{1}), HasRaisedError(errFoo))
}

func Test_Observable_Skip(t *testing.T) {
	obs := testObservable(0, 1, 2, 3, 4, 5).Skip(3)
	Assert(context.Background(), t, obs, HasItems(3, 4, 5))
}

func Test_Observable_SkipLast(t *testing.T) {
	obs := testObservable(0, 1, 2, 3, 4, 5).SkipLast(3)
	Assert(context.Background(), t, obs, HasItems(0, 1, 2))
}

func Test_Observable_SkipWhile(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5).SkipWhile(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i != 3
		default:
			return true
		}
	})

	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNotRaisedError())
}

func Test_Observable_StartWithIterable(t *testing.T) {
	obs := testObservable(4, 5, 6).StartWithIterable(testObservable(1, 2, 3))
	Assert(context.Background(), t, obs, HasItems(1, 2, 3, 4, 5, 6), HasNotRaisedError())
}

func Test_Observable_StartWithIterable_Error1(t *testing.T) {
	obs := testObservable(4, 5, 6).StartWithIterable(testObservable(1, errFoo, 3))
	Assert(context.Background(), t, obs, HasItems(1), HasRaisedError(errFoo))
}

func Test_Observable_StartWithIterable_Error2(t *testing.T) {
	obs := testObservable(4, errFoo, 6).StartWithIterable(testObservable(1, 2, 3))
	Assert(context.Background(), t, obs, HasItems(1, 2, 3, 4), HasRaisedError(errFoo))
}

func Test_Observable_SumFloat32(t *testing.T) {
	Assert(context.Background(), t, testObservable(float32(1.0), float32(2.0), float32(3.0)).SumFloat32(),
		HasItem(float32(6.)))
	Assert(context.Background(), t, testObservable(float32(1.1), 2, int8(3), int16(1), int32(1), int64(1)).SumFloat32(),
		HasItem(float32(9.1)))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).SumFloat32(), HasRaisedAnError())
	Assert(context.Background(), t, Empty().SumFloat32(), HasItem(float32(0)))
}

func Test_Observable_SumFloat64(t *testing.T) {
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).SumFloat64(),
		HasItem(6.6))
	Assert(context.Background(), t, testObservable(float32(1.0), 2, int8(3), 4., int16(1), int32(1), int64(1)).SumFloat64(),
		HasItem(13.))
	Assert(context.Background(), t, testObservable("x").SumFloat64(), HasRaisedAnError())
	Assert(context.Background(), t, Empty().SumFloat64(), HasItem(float64(0)))
}

func Test_Observable_SumInt64(t *testing.T) {
	Assert(context.Background(), t, testObservable(1, 2, 3).SumInt64(), HasItem(int64(6)))
	Assert(context.Background(), t, testObservable(int8(1), int(2), int16(3), int32(4), int64(5)).SumInt64(),
		HasItem(int64(15)))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).SumInt64(), HasRaisedAnError())
	Assert(context.Background(), t, Empty().SumInt64(), HasItem(int64(0)))
}

// TODO Fix race
// func Test_Observable_SumInt64_Parallel(t *testing.T) {
//	Assert(context.Background(), t, Range(1, 10000).SumInt64(WithPool(8)), HasItem(int64(50015001)))
//}

func Test_Observable_Take(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5).Take(3)
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_Observable_TakeLast(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5).TakeLast(3)
	Assert(context.Background(), t, obs, HasItems(3, 4, 5))
}

func Test_Observable_TakeLast_LessThanNth(t *testing.T) {
	obs := testObservable(4, 5).TakeLast(3)
	Assert(context.Background(), t, obs, HasItems(4, 5))
}

func Test_Observable_TakeLast_LessThanNth2(t *testing.T) {
	obs := testObservable(4, 5).TakeLast(100000)
	Assert(context.Background(), t, obs, HasItems(4, 5))
}

func Test_Observable_TakeUntil(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5).TakeUntil(func(item interface{}) bool {
		return item == 3
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_Observable_TakeWhile(t *testing.T) {
	obs := testObservable(1, 2, 3, 4, 5).TakeWhile(func(item interface{}) bool {
		return item != 3
	})
	Assert(context.Background(), t, obs, HasItems(1, 2))
}

func Test_Observable_ToMap(t *testing.T) {
	obs := testObservable(3, 4, 5, true, false).ToMap(func(i interface{}) (interface{}, error) {
		switch v := i.(type) {
		case int:
			return v, nil
		case bool:
			if v {
				return 0, nil
			}
			return 1, nil
		default:
			return i, nil
		}
	})
	Assert(context.Background(), t, obs, HasItem(map[interface{}]interface{}{
		3: 3,
		4: 4,
		5: 5,
		0: true,
		1: false,
	}))
}

func Test_Observable_ToMapWithValueSelector(t *testing.T) {
	keySelector := func(i interface{}) (interface{}, error) {
		switch v := i.(type) {
		case int:
			return v, nil
		case bool:
			if v {
				return 0, nil
			}
			return 1, nil
		default:
			return i, nil
		}
	}
	valueSelector := func(i interface{}) (interface{}, error) {
		switch v := i.(type) {
		case int:
			return v * 10, nil
		case bool:
			return v, nil
		default:
			return i, nil
		}
	}
	single := testObservable(3, 4, 5, true, false).ToMapWithValueSelector(keySelector, valueSelector)
	Assert(context.Background(), t, single, HasItem(map[interface{}]interface{}{
		3: 30,
		4: 40,
		5: 50,
		0: true,
		1: false,
	}))
}

func Test_Observable_ToSlice(t *testing.T) {
	s, err := testObservable(1, 2, 3).ToSlice(5)
	assert.Equal(t, []interface{}{1, 2, 3}, s)
	assert.Equal(t, 5, cap(s))
	assert.NoError(t, err)
}

func Test_Observable_ToSlice_Error(t *testing.T) {
	s, err := testObservable(1, 2, errFoo, 3).ToSlice(0)
	assert.Equal(t, []interface{}{1, 2}, s)
	assert.Equal(t, errFoo, err)
}

func Test_Observable_Unmarshal(t *testing.T) {
	obs := testObservable([]byte(`{"id":1}`), []byte(`{"id":2}`)).Unmarshal(json.Unmarshal,
		func() interface{} {
			return &testStruct{}
		})
	Assert(context.Background(), t, obs, HasItems(&testStruct{
		ID: 1,
	}, &testStruct{
		ID: 2,
	}))
}

func Test_Observable_Unmarshal_Error(t *testing.T) {
	obs := testObservable([]byte(`{"id":1`), []byte(`{"id":2}`)).Unmarshal(json.Unmarshal,
		func() interface{} {
			return &testStruct{}
		})
	Assert(context.Background(), t, obs, HasRaisedAnError())
}

func Test_Observable_ZipFromObservable(t *testing.T) {
	obs1 := testObservable(1, 2, 3)
	obs2 := testObservable(10, 20, 30)
	zipper := func(elem1 interface{}, elem2 interface{}) (interface{}, error) {
		switch v1 := elem1.(type) {
		case int:
			switch v2 := elem2.(type) {
			case int:
				return v1 + v2, nil
			}
		}
		return 0, nil
	}
	zip := obs1.ZipFromIterable(obs2, zipper)
	Assert(context.Background(), t, zip, HasItems(11, 22, 33))
}

func Test_Observable_ZipFromObservable_DifferentLength1(t *testing.T) {
	obs1 := testObservable(1, 2, 3)
	obs2 := testObservable(10, 20)
	zipper := func(elem1 interface{}, elem2 interface{}) (interface{}, error) {
		switch v1 := elem1.(type) {
		case int:
			switch v2 := elem2.(type) {
			case int:
				return v1 + v2, nil
			}
		}
		return 0, nil
	}
	zip := obs1.ZipFromIterable(obs2, zipper)
	Assert(context.Background(), t, zip, HasItems(11, 22))
}

func Test_Observable_ZipFromObservable_DifferentLength2(t *testing.T) {
	obs1 := testObservable(1, 2)
	obs2 := testObservable(10, 20, 30)
	zipper := func(elem1 interface{}, elem2 interface{}) (interface{}, error) {
		switch v1 := elem1.(type) {
		case int:
			switch v2 := elem2.(type) {
			case int:
				return v1 + v2, nil
			}
		}
		return 0, nil
	}
	zip := obs1.ZipFromIterable(obs2, zipper)
	Assert(context.Background(), t, zip, HasItems(11, 22))
}

func Test_Observable_Option_WithOnErrorStrategy_Single(t *testing.T) {
	obs := testObservable(1, 2, 3).
		Map(func(i interface{}) (interface{}, error) {
			if i == 2 {
				return nil, errFoo
			}
			return i, nil
		}, WithErrorStrategy(Continue))
	Assert(context.Background(), t, obs, HasItems(1, 3), HasRaisedError(errFoo))
}

func Test_Observable_Option_WithOnErrorStrategy_Propagate(t *testing.T) {
	obs := testObservable(1, 2, 3).
		Map(func(i interface{}) (interface{}, error) {
			if i == 1 {
				return nil, errFoo
			}
			return i, nil
		}).
		Map(func(i interface{}) (interface{}, error) {
			if i == 2 {
				return nil, errBar
			}
			return i, nil
		}, WithErrorStrategy(Continue))
	Assert(context.Background(), t, obs, HasItems(3), HasRaisedErrors(errFoo, errBar))
}

func Test_Observable_Option_SimpleCapacity(t *testing.T) {
	ch := Just(1, WithBufferedChannel(5)).Observe()
	assert.Equal(t, 5, cap(ch))
}

func Test_Observable_Option_ComposedCapacity(t *testing.T) {
	obs1 := Just(1).Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(11))
	obs2 := obs1.Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(12))

	assert.Equal(t, 11, cap(obs1.Observe()))
	assert.Equal(t, 12, cap(obs2.Observe()))
}
