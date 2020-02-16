package rxgo

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/pkg/errors"

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

	Assert(context.Background(), t, testObservable(1, 2, 3).All(predicateAllInt),
		HasItem(true), HasNotRaisedError())
	Assert(context.Background(), t, testObservable(1, "x", 3).All(predicateAllInt),
		HasItem(false), HasNotRaisedError())
}

func Test_Observable_AverageFloat32(t *testing.T) {
	Assert(context.Background(), t, testObservable(float32(1), float32(2), float32(3)).AverageFloat32(), HasItem(float32(2)))
	Assert(context.Background(), t, testObservable(float32(1), float32(20)).AverageFloat32(), HasItem(float32(10.5)))
	Assert(context.Background(), t, FromEmpty().AverageFloat32(), HasItem(0))
	Assert(context.Background(), t, testObservable("x").AverageFloat32(), HasRaisedAnError())
}

func Test_Observable_AverageFloat64(t *testing.T) {
	Assert(context.Background(), t, testObservable(float64(1), float64(2), float64(3)).AverageFloat64(), HasItem(float64(2)))
	Assert(context.Background(), t, testObservable(float64(1), float64(20)).AverageFloat64(), HasItem(10.5))
	Assert(context.Background(), t, FromEmpty().AverageFloat64(), HasItem(0))
	Assert(context.Background(), t, testObservable("x").AverageFloat64(), HasRaisedAnError())
}

func Test_Observable_AverageInt(t *testing.T) {
	Assert(context.Background(), t, testObservable(1, 2, 3).AverageInt(), HasItem(2))
	Assert(context.Background(), t, testObservable(1, 20).AverageInt(), HasItem(10))
	Assert(context.Background(), t, FromEmpty().AverageInt(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt(), HasRaisedAnError())
}

func Test_Observable_AverageInt8(t *testing.T) {
	Assert(context.Background(), t, testObservable(int8(1), int8(2), int8(3)).AverageInt8(), HasItem(int8(2)))
	Assert(context.Background(), t, testObservable(int8(1), int8(20)).AverageInt8(), HasItem(int8(10)))
	Assert(context.Background(), t, FromEmpty().AverageInt8(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt8(), HasRaisedAnError())
}

func Test_Observable_AverageInt16(t *testing.T) {
	Assert(context.Background(), t, testObservable(int16(1), int16(2), int16(3)).AverageInt16(), HasItem(int16(2)))
	Assert(context.Background(), t, testObservable(int16(1), int16(20)).AverageInt16(), HasItem(int16(10)))
	Assert(context.Background(), t, FromEmpty().AverageInt16(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt16(), HasRaisedAnError())
}

func Test_Observable_AverageInt32(t *testing.T) {
	Assert(context.Background(), t, testObservable(int32(1), int32(2), int32(3)).AverageInt32(), HasItem(int32(2)))
	Assert(context.Background(), t, testObservable(int32(1), int32(20)).AverageInt32(), HasItem(int32(10)))
	Assert(context.Background(), t, FromEmpty().AverageInt32(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt32(), HasRaisedAnError())
}

func Test_Observable_AverageInt64(t *testing.T) {
	Assert(context.Background(), t, testObservable(int64(1), int64(2), int64(3)).AverageInt64(), HasItem(int64(2)))
	Assert(context.Background(), t, testObservable(int64(1), int64(20)).AverageInt64(), HasItem(int64(10)))
	Assert(context.Background(), t, FromEmpty().AverageInt64(), HasItem(0))
	Assert(context.Background(), t, testObservable(1.1, 2.2, 3.3).AverageInt64(), HasRaisedAnError())
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

	ch <- FromValue(1)
	close(ch)

	<-obs.Observe()
	timespan.AssertCalled(t, "duration")
}

func Test_Observable_BufferWithTime_IllegalInput(t *testing.T) {
	Assert(context.Background(), t, FromEmpty().BufferWithTime(nil, nil), HasRaisedAnError())
	Assert(context.Background(), t, FromEmpty().BufferWithTime(WithDuration(0*time.Second), nil), HasRaisedAnError())
}

func Test_Observable_BufferWithTime_NilTimeshift(t *testing.T) {
	just := testObservable(1, 2, 3)
	obs := just.BufferWithTime(WithDuration(1*time.Second), nil)
	Assert(context.Background(), t, obs, HasSomeItems())
}

func Test_Observable_BufferWithTime_Error(t *testing.T) {
	just := testObservable(1, 2, 3, errFoo)
	obs := just.BufferWithTime(WithDuration(1*time.Second), nil)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2, 3}), HasRaisedError(errFoo))
}

func Test_Observable_BufferWithTimeOrCount_InvalidInputs(t *testing.T) {
	obs := FromEmpty().BufferWithTimeOrCount(nil, 5)
	Assert(context.Background(), t, obs, HasRaisedAnError())

	obs = FromEmpty().BufferWithTimeOrCount(WithDuration(0), 5)
	Assert(context.Background(), t, obs, HasRaisedAnError())

	obs = FromEmpty().BufferWithTimeOrCount(WithDuration(time.Millisecond*5), 0)
	Assert(context.Background(), t, obs, HasRaisedAnError())
}

func Test_Observable_BufferWithTimeOrCount_Count(t *testing.T) {
	just := testObservable(1, 2, 3)
	obs := just.BufferWithTimeOrCount(WithDuration(1*time.Second), 2)
	Assert(context.Background(), t, obs, HasItems([]interface{}{1, 2}, []interface{}{3}))
}

func Test_Observable_BufferWithTimeOrCount_MockedTime(t *testing.T) {
	ch := make(chan Item)
	from := FromChannel(ch)

	timespan := new(mockDuration)
	timespan.On("duration").Return(1 * time.Millisecond)

	obs := from.BufferWithTimeOrCount(timespan, 5)

	time.Sleep(50 * time.Millisecond)
	ch <- FromValue(1)
	close(ch)

	<-obs.Observe()
	timespan.AssertCalled(t, "duration")
}

func Test_Observable_BufferWithTimeOrCount_Error(t *testing.T) {
	just := testObservable(1, 2, 3, errFoo, 4)
	obs := just.BufferWithTimeOrCount(WithDuration(10*time.Second), 2)
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
	obs := FromEmpty().DefaultIfEmpty(3)
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

func Test_Observable_ElementAt(t *testing.T) {
	obs := testObservable(0, 1, 2, 3, 4).ElementAt(2)
	Assert(context.Background(), t, obs, HasItems(2))
}

func Test_Observable_ElementAt_Error(t *testing.T) {
	obs := testObservable(0, 1, 2, 3, 4).ElementAt(10)
	Assert(context.Background(), t, obs, HasNoItem(), HasRaisedAnError())
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
	obs := FromEmpty().First()
	Assert(context.Background(), t, obs, HasNoItem())
}

func Test_Observable_FirstOrDefault_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).FirstOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(1))
}

func Test_Observable_FirstOrDefault_Empty(t *testing.T) {
	obs := FromEmpty().FirstOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(10))
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

func Test_Observable_Last_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).Last()
	Assert(context.Background(), t, obs, HasItem(3))
}

func Test_Observable_Last_Empty(t *testing.T) {
	obs := FromEmpty().Last()
	Assert(context.Background(), t, obs, HasNoItem())
}

func Test_Observable_LastOrDefault_NotEmpty(t *testing.T) {
	obs := testObservable(1, 2, 3).LastOrDefault(10)
	Assert(context.Background(), t, obs, HasItem(3))
}

func Test_Observable_LastOrDefault_Empty(t *testing.T) {
	obs := FromEmpty().LastOrDefault(10)
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
			ch <- FromValue(i)
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

func Test_Observable_Observe(t *testing.T) {
	got := make([]int, 0)
	ch := testObservable(1, 2, 3).Observe()
	for item := range ch {
		got = append(got, item.Value.(int))
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
	obs := testObservable(1, 2, 3, 4, 5).Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	})
	Assert(context.Background(), t, obs, HasItem(15), HasNotRaisedError())
}

func Test_Observable_Empty(t *testing.T) {
	obs := FromEmpty().Reduce(func(acc interface{}, elem interface{}) (interface{}, error) {
		return 0, nil
	})
	Assert(context.Background(), t, obs, HasNoItem(), HasNotRaisedError())
}

func Test_Observable_Error(t *testing.T) {
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
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		if i == 2 {
			next <- FromValue(3)
			done()
		} else {
			i++
			next <- FromError(errFoo)
			done()
		}
	}).Retry(3)
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 2, 1, 2, 3), HasNotRaisedError())
}

func Test_Observable_Retry_Error(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromError(errFoo)
		done()
	}).Retry(3)
	Assert(context.Background(), t, obs, HasItems(1, 2, 1, 2, 1, 2, 1, 2), HasRaisedError(errFoo))
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
	assert.Equal(t, FromValue(1), <-ch)
	assert.Equal(t, FromValue(2), <-ch)
	assert.Equal(t, FromValue(3), <-ch)
	assert.Equal(t, FromError(errFoo), <-ch)
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

func Test_Observable_ToSlice(t *testing.T) {
	single := testObservable(1, 2, 3).ToSlice()
	Assert(context.Background(), t, single, HasItem([]interface{}{1, 2, 3}))
}

func Test_Observable_ToSlice_Error(t *testing.T) {
	single := testObservable(1, 2, errFoo, 3).ToSlice()
	Assert(context.Background(), t, single, HasNoItem(), HasRaisedError(errFoo))
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
