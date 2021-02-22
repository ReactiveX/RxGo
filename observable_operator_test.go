package rxgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"go.uber.org/goleak"

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
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Range(1, 10000).All(predicateAllInt),
		HasItem(true), HasNoError())
}

func Test_Observable_All_False(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1, "x", 3).All(predicateAllInt),
		HasItem(false), HasNoError())
}

func Test_Observable_All_Parallel_True(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Range(1, 10000).All(predicateAllInt, WithContext(ctx), WithCPUPool()),
		HasItem(true), HasNoError())
}

func Test_Observable_All_Parallel_False(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1, "x", 3).All(predicateAllInt, WithContext(ctx), WithCPUPool()),
		HasItem(false), HasNoError())
}

func Test_Observable_All_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1, errFoo, 3).All(predicateAllInt, WithContext(ctx), WithCPUPool()),
		HasError(errFoo))
}

func Test_Observable_AverageFloat32(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float32(1), float32(20)).AverageFloat32(), HasItem(float32(10.5)))
	Assert(ctx, t, testObservable(ctx, float64(1), float64(20)).AverageFloat32(), HasItem(float32(10.5)))
}

func Test_Observable_AverageFloat32_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().AverageFloat32(), HasItem(0))
}

func Test_Observable_AverageFloat32_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, "x").AverageFloat32(), HasAnError())
}

func Test_Observable_AverageFloat32_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float32(1), float32(20)).AverageFloat32(), HasItem(float32(10.5)))
	Assert(ctx, t, testObservable(ctx, float64(1), float64(20)).AverageFloat32(), HasItem(float32(10.5)))
}

func Test_Observable_AverageFloat32_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().AverageFloat32(WithContext(ctx), WithCPUPool()),
		HasItem(0))
}

func Test_Observable_AverageFloat32_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, "x").AverageFloat32(WithContext(ctx), WithCPUPool()),
		HasAnError())
}

func Test_Observable_AverageFloat64(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float64(1), float64(20)).AverageFloat64(), HasItem(10.5))
	Assert(ctx, t, testObservable(ctx, float32(1), float32(20)).AverageFloat64(), HasItem(10.5))
}

func Test_Observable_AverageFloat64_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().AverageFloat64(), HasItem(0))
}

func Test_Observable_AverageFloat64_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, "x").AverageFloat64(), HasAnError())
}

func Test_Observable_AverageFloat64_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float64(1), float64(20)).AverageFloat64(), HasItem(float64(10.5)))
}

func Test_Observable_AverageFloat64_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().AverageFloat64(WithContext(ctx), WithCPUPool()),
		HasItem(0))
}

func Test_Observable_AverageFloat64_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, "x").AverageFloat64(WithContext(ctx), WithCPUPool()),
		HasAnError())
}

func Test_Observable_AverageInt(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1, 2, 3).AverageInt(), HasItem(2))
	Assert(ctx, t, testObservable(ctx, 1, 20).AverageInt(), HasItem(10))
	Assert(ctx, t, Empty().AverageInt(), HasItem(0))
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).AverageInt(), HasAnError())
}

func Test_Observable_AverageInt8(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, int8(1), int8(2), int8(3)).AverageInt8(), HasItem(int8(2)))
	Assert(ctx, t, testObservable(ctx, int8(1), int8(20)).AverageInt8(), HasItem(int8(10)))
	Assert(ctx, t, Empty().AverageInt8(), HasItem(0))
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).AverageInt8(), HasAnError())
}

func Test_Observable_AverageInt16(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, int16(1), int16(2), int16(3)).AverageInt16(), HasItem(int16(2)))
	Assert(ctx, t, testObservable(ctx, int16(1), int16(20)).AverageInt16(), HasItem(int16(10)))
	Assert(ctx, t, Empty().AverageInt16(), HasItem(0))
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).AverageInt16(), HasAnError())
}

func Test_Observable_AverageInt32(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, int32(1), int32(2), int32(3)).AverageInt32(), HasItem(int32(2)))
	Assert(ctx, t, testObservable(ctx, int32(1), int32(20)).AverageInt32(), HasItem(int32(10)))
	Assert(ctx, t, Empty().AverageInt32(), HasItem(0))
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).AverageInt32(), HasAnError())
}

func Test_Observable_AverageInt64(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, int64(1), int64(2), int64(3)).AverageInt64(), HasItem(int64(2)))
	Assert(ctx, t, testObservable(ctx, int64(1), int64(20)).AverageInt64(), HasItem(int64(10)))
	Assert(ctx, t, Empty().AverageInt64(), HasItem(0))
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).AverageInt64(), HasAnError())
}

func Test_Observable_BackOffRetry(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := 0
	backOffCfg := backoff.NewExponentialBackOff()
	backOffCfg.InitialInterval = time.Nanosecond
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		if i == 2 {
			next <- Of(3)
		} else {
			i++
			next <- Error(errFoo)
		}
	}}).BackOffRetry(backoff.WithMaxRetries(backOffCfg, 3))
	Assert(ctx, t, obs, HasItems(1, 2, 1, 2, 1, 2, 3), HasNoError())
}

func Test_Observable_BackOffRetry_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	backOffCfg := backoff.NewExponentialBackOff()
	backOffCfg.InitialInterval = time.Nanosecond
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Error(errFoo)
	}}).BackOffRetry(backoff.WithMaxRetries(backOffCfg, 3))
	Assert(ctx, t, obs, HasItems(1, 2, 1, 2, 1, 2, 1, 2), HasError(errFoo))
}

func Test_Observable_BufferWithCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5, 6).BufferWithCount(3)
	Assert(ctx, t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4, 5, 6}))
}

func Test_Observable_BufferWithCount_IncompleteLastItem(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4).BufferWithCount(3)
	Assert(ctx, t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4}))
}

func Test_Observable_BufferWithCount_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, errFoo).BufferWithCount(3)
	Assert(ctx, t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4}), HasError(errFoo))
}

func Test_Observable_BufferWithCount_InputError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4).BufferWithCount(0)
	Assert(ctx, t, obs, HasAnError())
}

func Test_Observable_BufferWithTime_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Just(1, 2, 3)().BufferWithTime(WithDuration(30 * time.Millisecond))
	Assert(ctx, t, obs, HasItems(
		[]interface{}{1, 2, 3},
	))
}

func Test_Observable_BufferWithTime_Multiple(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan Item, 1)
	obs := FromChannel(ch)
	obs = obs.BufferWithTime(WithDuration(30 * time.Millisecond))
	go func() {
		for i := 0; i < 10; i++ {
			ch <- Of(i)
		}
		close(ch)
	}()
	Assert(ctx, t, obs, CustomPredicate(func(items []interface{}) error {
		if len(items) == 0 {
			return errors.New("items should not be nil")
		}
		return nil
	}))
}

func Test_Observable_BufferWithTimeOrCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan Item, 1)
	obs := FromChannel(ch)
	obs = obs.BufferWithTimeOrCount(WithDuration(30*time.Millisecond), 100)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- Of(i)
		}
		close(ch)
	}()
	Assert(ctx, t, obs, CustomPredicate(func(items []interface{}) error {
		if len(items) == 0 {
			return errors.New("items should not be nil")
		}
		return nil
	}))
}

func Test_Observable_Contain(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	predicate := func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i == 2
		default:
			return false
		}
	}

	Assert(ctx, t,
		testObservable(ctx, 1, 2, 3).Contains(predicate),
		HasItem(true))
	Assert(ctx, t,
		testObservable(ctx, 1, 5, 3).Contains(predicate),
		HasItem(false))
}

func Test_Observable_Contain_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	predicate := func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i == 2
		default:
			return false
		}
	}

	Assert(ctx, t,
		testObservable(ctx, 1, 2, 3).Contains(predicate, WithContext(ctx), WithCPUPool()),
		HasItem(true))
	Assert(ctx, t,
		testObservable(ctx, 1, 5, 3).Contains(predicate, WithContext(ctx), WithCPUPool()),
		HasItem(false))
}

func Test_Observable_Count(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Range(1, 10000).Count(),
		HasItem(int64(10001)))
}

func Test_Observable_Count_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Range(1, 10000).Count(WithCPUPool()),
		HasItem(int64(10001)))
}

// FIXME
//func Test_Observable_Debounce(t *testing.T) {
//	defer goleak.VerifyNone(t)
//	ctx, obs, d := timeCausality(1, tick, 2, tick, 3, 4, 5, tick, 6, tick)
//	Assert(ctx, t, obs.Debounce(d, WithBufferedChannel(10), WithContext(ctx)),
//		HasItems(1, 2, 5, 6))
//}

func Test_Observable_Debounce_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, obs, d := timeCausality(1, tick, 2, tick, 3, errFoo, 5, tick, 6, tick)
	Assert(ctx, t, obs.Debounce(d, WithBufferedChannel(10), WithContext(ctx)),
		HasItems(1, 2), HasError(errFoo))
}

func Test_Observable_DefaultIfEmpty_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().DefaultIfEmpty(3)
	Assert(ctx, t, obs, HasItems(3))
}

func Test_Observable_DefaultIfEmpty_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2).DefaultIfEmpty(3)
	Assert(ctx, t, obs, HasItems(1, 2))
}

func Test_Observable_DefaultIfEmpty_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().DefaultIfEmpty(3, WithCPUPool())
	Assert(ctx, t, obs, HasItems(3))
}

func Test_Observable_DefaultIfEmpty_Parallel_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2).DefaultIfEmpty(3, WithCPUPool())
	Assert(ctx, t, obs, HasItems(1, 2))
}

func Test_Observable_Distinct(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, 1, 3).Distinct(func(_ context.Context, item interface{}) (interface{}, error) {
		return item, nil
	})
	Assert(ctx, t, obs, HasItems(1, 2, 3), HasNoError())
}

func Test_Observable_Distinct_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, errFoo, 3).Distinct(func(_ context.Context, item interface{}) (interface{}, error) {
		return item, nil
	})
	Assert(ctx, t, obs, HasItems(1, 2), HasError(errFoo))
}

func Test_Observable_Distinct_Error2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, 2, 3, 4).Distinct(func(_ context.Context, item interface{}) (interface{}, error) {
		if item.(int) == 3 {
			return nil, errFoo
		}
		return item, nil
	})
	Assert(ctx, t, obs, HasItems(1, 2), HasError(errFoo))
}

func Test_Observable_Distinct_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, 1, 3).Distinct(func(_ context.Context, item interface{}) (interface{}, error) {
		return item, nil
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItemsNoOrder(1, 2, 3), HasNoError())
}

func Test_Observable_Distinct_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, errFoo).Distinct(func(_ context.Context, item interface{}) (interface{}, error) {
		return item, nil
	}, WithContext(ctx), WithCPUPool())
	Assert(ctx, t, obs, HasError(errFoo))
}

func Test_Observable_Distinct_Parallel_Error2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, 2, 3, 4).Distinct(func(_ context.Context, item interface{}) (interface{}, error) {
		if item.(int) == 3 {
			return nil, errFoo
		}
		return item, nil
	}, WithContext(ctx), WithCPUPool())
	Assert(ctx, t, obs, HasError(errFoo))
}

func Test_Observable_DistinctUntilChanged(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, 1, 3).DistinctUntilChanged(func(_ context.Context, item interface{}) (interface{}, error) {
		return item, nil
	})
	Assert(ctx, t, obs, HasItems(1, 2, 1, 3))
}

func Test_Observable_DistinctUntilChanged_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 2, 1, 3).DistinctUntilChanged(func(_ context.Context, item interface{}) (interface{}, error) {
		return item, nil
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItems(1, 2, 1, 3))
}

func Test_Observable_DoOnCompleted_NoError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := false
	<-testObservable(ctx, 1, 2, 3).DoOnCompleted(func() {
		called = true
	})
	assert.True(t, called)
}

func Test_Observable_DoOnCompleted_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := false
	<-testObservable(ctx, 1, errFoo, 3).DoOnCompleted(func() {
		called = true
	})
	assert.True(t, called)
}

func Test_Observable_DoOnError_NoError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var got error
	<-testObservable(ctx, 1, 2, 3).DoOnError(func(err error) {
		got = err
	})
	assert.Nil(t, got)
}

func Test_Observable_DoOnError_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var got error
	<-testObservable(ctx, 1, errFoo, 3).DoOnError(func(err error) {
		got = err
	})
	assert.Equal(t, errFoo, got)
}

func Test_Observable_DoOnNext_NoError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := make([]interface{}, 0)
	<-testObservable(ctx, 1, 2, 3).DoOnNext(func(i interface{}) {
		s = append(s, i)
	})
	assert.Equal(t, []interface{}{1, 2, 3}, s)
}

func Test_Observable_DoOnNext_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := make([]interface{}, 0)
	<-testObservable(ctx, 1, errFoo, 3).DoOnNext(func(i interface{}) {
		s = append(s, i)
	})
	assert.Equal(t, []interface{}{1}, s)
}

func Test_Observable_ElementAt(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(0, 10000).ElementAt(10000)
	Assert(ctx, t, obs, HasItems(10000))
}

func Test_Observable_ElementAt_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(0, 10000).ElementAt(10000, WithCPUPool())
	Assert(ctx, t, obs, HasItems(10000))
}

func Test_Observable_ElementAt_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 0, 1, 2, 3, 4).ElementAt(10)
	Assert(ctx, t, obs, IsEmpty(), HasAnError())
}

func Test_Observable_Error_NoError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.NoError(t, testObservable(ctx, 1, 2, 3).Error())
}

func Test_Observable_Error_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.Equal(t, errFoo, testObservable(ctx, 1, errFoo, 3).Error())
}

func Test_Observable_Errors_NoError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.Equal(t, 0, len(testObservable(ctx, 1, 2, 3).Errors()))
}

func Test_Observable_Errors_OneError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.Equal(t, 1, len(testObservable(ctx, 1, errFoo, 3).Errors()))
}

func Test_Observable_Errors_MultipleError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.Equal(t, 2, len(testObservable(ctx, 1, errFoo, errBar).Errors()))
}

func Test_Observable_Errors_MultipleErrorFromMap(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errs := testObservable(ctx, 1, 2, 3, 4).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		if i == 2 {
			return nil, errFoo
		}
		if i == 3 {
			return nil, errBar
		}
		return i, nil
	}, WithErrorStrategy(ContinueOnError)).Errors()
	assert.Equal(t, 2, len(errs))
}

func Test_Observable_Filter(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4).Filter(
		func(i interface{}) bool {
			return i.(int)%2 == 0
		})
	Assert(ctx, t, obs, HasItems(2, 4), HasNoError())
}

func Test_Observable_Filter_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4).Filter(
		func(i interface{}) bool {
			return i.(int)%2 == 0
		}, WithCPUPool())
	Assert(ctx, t, obs, HasItemsNoOrder(2, 4), HasNoError())
}

func Test_Observable_Find_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).Find(func(i interface{}) bool {
		return i == 2
	})
	Assert(ctx, t, obs, HasItem(2))
}

func Test_Observable_Find_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().Find(func(_ interface{}) bool {
		return true
	})
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_First_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).First()
	Assert(ctx, t, obs, HasItem(1))
}

func Test_Observable_First_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().First()
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_First_Parallel_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).First(WithCPUPool())
	Assert(ctx, t, obs, HasItem(1))
}

func Test_Observable_First_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().First(WithCPUPool())
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_FirstOrDefault_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).FirstOrDefault(10)
	Assert(ctx, t, obs, HasItem(1))
}

func Test_Observable_FirstOrDefault_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().FirstOrDefault(10)
	Assert(ctx, t, obs, HasItem(10))
}

func Test_Observable_FirstOrDefault_Parallel_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).FirstOrDefault(10, WithCPUPool())
	Assert(ctx, t, obs, HasItem(1))
}

func Test_Observable_FirstOrDefault_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().FirstOrDefault(10, WithCPUPool())
	Assert(ctx, t, obs, HasItem(10))
}

func Test_Observable_FlatMap(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).FlatMap(func(i Item) Observable {
		return testObservable(ctx, i.V.(int)+1, i.V.(int)*10)
	})
	Assert(ctx, t, obs, HasItems(2, 10, 3, 20, 4, 30))
}

func Test_Observable_FlatMap_Error1(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).FlatMap(func(i Item) Observable {
		if i.V == 2 {
			return testObservable(ctx, errFoo)
		}
		return testObservable(ctx, i.V.(int)+1, i.V.(int)*10)
	})
	Assert(ctx, t, obs, HasItems(2, 10), HasError(errFoo))
}

func Test_Observable_FlatMap_Error2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, errFoo, 3).FlatMap(func(i Item) Observable {
		if i.Error() {
			return testObservable(ctx, 0)
		}
		return testObservable(ctx, i.V.(int)+1, i.V.(int)*10)
	})
	Assert(ctx, t, obs, HasItems(2, 10, 0, 4, 30), HasNoError())
}

func Test_Observable_FlatMap_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).FlatMap(func(i Item) Observable {
		return testObservable(ctx, i.V.(int)+1, i.V.(int)*10)
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItemsNoOrder(2, 10, 3, 20, 4, 30))
}

func Test_Observable_FlatMap_Parallel_Error1(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).FlatMap(func(i Item) Observable {
		if i.V == 2 {
			return testObservable(ctx, errFoo)
		}
		return testObservable(ctx, i.V.(int)+1, i.V.(int)*10)
	})
	Assert(ctx, t, obs, HasError(errFoo))
}

func Test_Observable_ForEach_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	count := 0
	var gotErr error
	done := make(chan struct{})

	obs := testObservable(ctx, 1, 2, 3, errFoo)
	obs.ForEach(func(i interface{}) {
		count += i.(int)
	}, func(err error) {
		gotErr = err
		select {
		case <-ctx.Done():
			return
		case done <- struct{}{}:
		}
	}, func() {
		select {
		case <-ctx.Done():
			return
		case done <- struct{}{}:
		}
	}, WithContext(ctx))

	// We avoid using the assertion API on purpose
	<-done
	assert.Equal(t, 6, count)
	assert.Equal(t, errFoo, gotErr)
}

func Test_Observable_ForEach_Done(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	count := 0
	var gotErr error
	done := make(chan struct{})

	obs := testObservable(ctx, 1, 2, 3)
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
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).IgnoreElements()
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_IgnoreElements_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, errFoo, 3).IgnoreElements()
	Assert(ctx, t, obs, IsEmpty(), HasError(errFoo))
}

func Test_Observable_IgnoreElements_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).IgnoreElements(WithCPUPool())
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_IgnoreElements_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, errFoo, 3).IgnoreElements(WithCPUPool())
	Assert(ctx, t, obs, IsEmpty(), HasError(errFoo))
}

func Test_Observable_GroupBy(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	Assert(ctx, t, s[0].(Observable), HasItems(0, 3, 6, 9), HasNoError())
	Assert(ctx, t, s[1].(Observable), HasItems(1, 4, 7, 10), HasNoError())
	Assert(ctx, t, s[2].(Observable), HasItems(2, 5, 8), HasNoError())
}

func Test_Observable_GroupBy_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	Assert(ctx, t, s[0].(Observable), HasAnError())
	Assert(ctx, t, s[1].(Observable), HasAnError())
	Assert(ctx, t, s[2].(Observable), HasAnError())
}

func Test_Observable_GroupByDynamic(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	count := 3
	max := 10

	obs := Range(0, max).GroupByDynamic(func(item Item) string {
		if item.V == 10 {
			return "10"
		}
		return strconv.Itoa(item.V.(int) % count)
	}, WithBufferedChannel(max))
	s, err := obs.ToSlice(0)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if len(s) != 4 {
		assert.FailNow(t, "length", "got=%d, expected=%d", len(s), 4)
	}

	Assert(ctx, t, s[0].(GroupedObservable), HasItems(0, 3, 6, 9), HasNoError())
	assert.Equal(t, "0", s[0].(GroupedObservable).Key)
	Assert(ctx, t, s[1].(GroupedObservable), HasItems(1, 4, 7), HasNoError())
	assert.Equal(t, "1", s[1].(GroupedObservable).Key)
	Assert(ctx, t, s[2].(GroupedObservable), HasItems(2, 5, 8), HasNoError())
	assert.Equal(t, "2", s[2].(GroupedObservable).Key)
	Assert(ctx, t, s[3].(GroupedObservable), HasItems(10), HasNoError())
	assert.Equal(t, "10", s[3].(GroupedObservable).Key)
}

func joinTest(ctx context.Context, t *testing.T, left, right []interface{}, window Duration, expected []int64) {
	leftObs := testObservable(ctx, left...)
	rightObs := testObservable(ctx, right...)

	obs := leftObs.Join(func(ctx context.Context, l, r interface{}) (interface{}, error) {
		return map[string]interface{}{
			"l": l,
			"r": r,
		}, nil
	},
		rightObs,
		func(i interface{}) time.Time {
			return time.Unix(0, i.(map[string]int64)["tt"]*1000000)
		},
		window,
	)

	Assert(ctx, t, obs, CustomPredicate(func(items []interface{}) error {
		actuals := make([]int64, 0)
		for _, p := range items {
			val := p.(map[string]interface{})
			actuals = append(actuals, val["l"].(map[string]int64)["V"], val["r"].(map[string]int64)["V"])
		}
		assert.Equal(t, expected, actuals)
		return nil
	}))
}

func Test_Observable_Join1(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	left := []interface{}{
		map[string]int64{"tt": 1, "V": 1},
		map[string]int64{"tt": 4, "V": 2},
		map[string]int64{"tt": 7, "V": 3},
	}
	right := []interface{}{
		map[string]int64{"tt": 2, "V": 5},
		map[string]int64{"tt": 3, "V": 6},
		map[string]int64{"tt": 5, "V": 7},
	}
	window := WithDuration(2 * time.Millisecond)
	expected := []int64{
		1, 5,
		1, 6,
		2, 5,
		2, 6,
		2, 7,
		3, 7,
	}

	joinTest(ctx, t, left, right, window, expected)
}

func Test_Observable_Join2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	left := []interface{}{
		map[string]int64{"tt": 1, "V": 1},
		map[string]int64{"tt": 3, "V": 2},
		map[string]int64{"tt": 5, "V": 3},
		map[string]int64{"tt": 9, "V": 4},
	}
	right := []interface{}{
		map[string]int64{"tt": 2, "V": 1},
		map[string]int64{"tt": 7, "V": 2},
		map[string]int64{"tt": 10, "V": 3},
	}
	window := WithDuration(2 * time.Millisecond)
	expected := []int64{
		1, 1,
		2, 1,
		3, 2,
		4, 2,
		4, 3,
	}

	joinTest(ctx, t, left, right, window, expected)
}

func Test_Observable_Join3(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	left := []interface{}{
		map[string]int64{"tt": 1, "V": 1},
		map[string]int64{"tt": 2, "V": 2},
		map[string]int64{"tt": 3, "V": 3},
		map[string]int64{"tt": 4, "V": 4},
	}
	right := []interface{}{
		map[string]int64{"tt": 5, "V": 1},
		map[string]int64{"tt": 6, "V": 2},
		map[string]int64{"tt": 7, "V": 3},
	}
	window := WithDuration(3 * time.Millisecond)
	expected := []int64{
		2, 1,
		3, 1,
		3, 2,
		4, 1,
		4, 2,
		4, 3,
	}

	joinTest(ctx, t, left, right, window, expected)
}

func Test_Observable_Join_Error_OnLeft(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	left := []interface{}{
		map[string]int64{"tt": 1, "V": 1},
		map[string]int64{"tt": 3, "V": 2},
		errFoo,
		map[string]int64{"tt": 9, "V": 4},
	}
	right := []interface{}{
		map[string]int64{"tt": 2, "V": 1},
		map[string]int64{"tt": 7, "V": 2},
		map[string]int64{"tt": 10, "V": 3},
	}
	window := WithDuration(3 * time.Millisecond)
	expected := []int64{
		1, 1,
		2, 1,
	}

	joinTest(ctx, t, left, right, window, expected)
}

func Test_Observable_Join_Error_OnRight(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	left := []interface{}{
		map[string]int64{"tt": 1, "V": 1},
		map[string]int64{"tt": 3, "V": 2},
		map[string]int64{"tt": 5, "V": 3},
		map[string]int64{"tt": 9, "V": 4},
	}
	right := []interface{}{
		map[string]int64{"tt": 2, "V": 1},
		errFoo,
		map[string]int64{"tt": 10, "V": 3},
	}
	window := WithDuration(3 * time.Millisecond)
	expected := []int64{
		1, 1,
	}

	joinTest(ctx, t, left, right, window, expected)
}

func Test_Observable_Last_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).Last()
	Assert(ctx, t, obs, HasItem(3))
}

func Test_Observable_Last_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().Last()
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_Last_Parallel_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).Last(WithCPUPool())
	Assert(ctx, t, obs, HasItem(3))
}

func Test_Observable_Last_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().Last(WithCPUPool())
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_LastOrDefault_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).LastOrDefault(10)
	Assert(ctx, t, obs, HasItem(3))
}

func Test_Observable_LastOrDefault_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().LastOrDefault(10)
	Assert(ctx, t, obs, HasItem(10))
}

func Test_Observable_LastOrDefault_Parallel_NotEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).LastOrDefault(10, WithCPUPool())
	Assert(ctx, t, obs, HasItem(3))
}

func Test_Observable_LastOrDefault_Parallel_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().LastOrDefault(10, WithCPUPool())
	Assert(ctx, t, obs, HasItem(10))
}

func Test_Observable_Map_One(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(ctx, t, obs, HasItems(2, 3, 4), HasNoError())
}

func Test_Observable_Map_Multiple(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
	Assert(ctx, t, obs, HasItems(20, 30, 40), HasNoError())
}

func Test_Observable_Map_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, errFoo).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(ctx, t, obs, HasItems(2, 3, 4), HasError(errFoo))
}

func Test_Observable_Map_ReturnValueAndError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return 2, errFoo
	})
	Assert(ctx, t, obs, IsEmpty(), HasError(errFoo))
}

func Test_Observable_Map_Multiple_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := false
	obs := testObservable(ctx, 1, 2, 3).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return nil, errFoo
	}).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		called = true
		return nil, nil
	})
	Assert(ctx, t, obs, IsEmpty(), HasError(errFoo))
	assert.False(t, called)
}

func Test_Observable_Map_Cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	next := make(chan Item)

	ctx, cancel := context.WithCancel(context.Background())
	obs := FromChannel(next).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}, WithContext(ctx))
	cancel()
	Assert(ctx, t, obs, IsEmpty(), HasNoError())
}

func Test_Observable_Map_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const len = 10
	ch := make(chan Item, len)
	go func() {
		for i := 0; i < len; i++ {
			ch <- Of(i)
		}
		close(ch)
	}()

	obs := FromChannel(ch).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}, WithPool(len))
	Assert(ctx, t, obs, HasItemsNoOrder(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), HasNoError())
}

func Test_Observable_Marshal(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, testStruct{
		ID: 1,
	}, testStruct{
		ID: 2,
	}).Marshal(json.Marshal)
	Assert(ctx, t, obs, HasItems([]byte(`{"id":1}`), []byte(`{"id":2}`)))
}

func Test_Observable_Marshal_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, testStruct{
		ID: 1,
	}, testStruct{
		ID: 2,
	}).Marshal(json.Marshal,
		// We cannot use HasItemsNoOrder function with a []byte
		WithPool(1))
	Assert(ctx, t, obs, HasItems([]byte(`{"id":1}`), []byte(`{"id":2}`)))
}

func Test_Observable_Max(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(0, 10000).Max(func(e1, e2 interface{}) int {
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
	Assert(ctx, t, obs, HasItem(10000))
}

func Test_Observable_Max_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(0, 10000).Max(func(e1, e2 interface{}) int {
		var i1 int
		if e1 == nil {
			i1 = 0
		} else {
			i1 = e1.(int)
		}

		var i2 int
		if e2 == nil {
			i2 = 0
		} else {
			i2 = e2.(int)
		}

		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		} else {
			return 0
		}
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItem(10000))
}

func Test_Observable_Min(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(0, 10000).Min(func(e1, e2 interface{}) int {
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
	Assert(ctx, t, obs, HasItem(0))
}

func Test_Observable_Min_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(0, 10000).Min(func(e1, e2 interface{}) int {
		i1 := e1.(int)
		i2 := e2.(int)
		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		} else {
			return 0
		}
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItem(0))
}

func Test_Observable_Observe(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	got := make([]int, 0)
	ch := testObservable(ctx, 1, 2, 3).Observe()
	for item := range ch {
		got = append(got, item.V.(int))
	}
	assert.Equal(t, []int{1, 2, 3}, got)
}

func Test_Observable_OnErrorResumeNext(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, errFoo, 4).OnErrorResumeNext(func(e error) Observable {
		return testObservable(ctx, 10, 20)
	})
	Assert(ctx, t, obs, HasItems(1, 2, 10, 20), HasNoError())
}

func Test_Observable_OnErrorReturn(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, errFoo, 4, errBar, 6).OnErrorReturn(func(err error) interface{} {
		return err.Error()
	})
	Assert(ctx, t, obs, HasItems(1, 2, "foo", 4, "bar", 6), HasNoError())
}

func Test_Observable_OnErrorReturnItem(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, errFoo, 4, errBar, 6).OnErrorReturnItem("foo")
	Assert(ctx, t, obs, HasItems(1, 2, "foo", 4, "foo", 6), HasNoError())
}

func Test_Observable_Reduce(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(1, 10000).Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	})
	Assert(ctx, t, obs, HasItem(50015001), HasNoError())
}

func Test_Observable_Reduce_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		return 0, nil
	})
	Assert(ctx, t, obs, IsEmpty(), HasNoError())
}

func Test_Observable_Reduce_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, errFoo, 4, 5).Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		return 0, nil
	})
	Assert(ctx, t, obs, IsEmpty(), HasError(errFoo))
}

func Test_Observable_Reduce_ReturnError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		if elem == 2 {
			return 0, errFoo
		}
		return elem, nil
	})
	Assert(ctx, t, obs, IsEmpty(), HasError(errFoo))
}

func Test_Observable_Reduce_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(1, 10000).Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b, nil
			}
		} else {
			return elem.(int), nil
		}
		return 0, errFoo
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItem(50015001), HasNoError())
}

func Test_Observable_Reduce_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(1, 10000).Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
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
	}, WithContext(ctx), WithCPUPool())
	Assert(ctx, t, obs, HasError(errFoo))
}

func Test_Observable_Reduce_Parallel_WithErrorStrategy(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Range(1, 10000).Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
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
	}, WithCPUPool(), WithErrorStrategy(ContinueOnError))
	Assert(ctx, t, obs, HasItem(50015000), HasError(errFoo))
}

func Test_Observable_Repeat(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repeat := testObservable(ctx, 1, 2, 3).Repeat(1, nil)
	Assert(ctx, t, repeat, HasItems(1, 2, 3, 1, 2, 3))
}

func Test_Observable_Repeat_Zero(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repeat := testObservable(ctx, 1, 2, 3).Repeat(0, nil)
	Assert(ctx, t, repeat, HasItems(1, 2, 3))
}

func Test_Observable_Repeat_NegativeCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repeat := testObservable(ctx, 1, 2, 3).Repeat(-2, nil)
	Assert(ctx, t, repeat, IsEmpty(), HasAnError())
}

func Test_Observable_Repeat_Infinite(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	repeat := testObservable(ctx, 1, 2, 3).Repeat(Infinite, nil, WithContext(ctx))
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	Assert(ctx, t, repeat, HasNoError(), CustomPredicate(func(items []interface{}) error {
		if len(items) == 0 {
			return errors.New("no items")
		}
		return nil
	}))
}

func Test_Observable_Repeat_Frequency(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	frequency := new(mockDuration)
	frequency.On("duration").Return(time.Millisecond)

	repeat := testObservable(ctx, 1, 2, 3).Repeat(1, frequency)
	Assert(ctx, t, repeat, HasItems(1, 2, 3, 1, 2, 3))
	frequency.AssertNumberOfCalls(t, "duration", 1)
	frequency.AssertExpectations(t)
}

func Test_Observable_Retry(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := 0
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		if i == 2 {
			next <- Of(3)
		} else {
			i++
			next <- Error(errFoo)
		}
	}}).Retry(3, func(err error) bool {
		return true
	})
	Assert(ctx, t, obs, HasItems(1, 2, 1, 2, 1, 2, 3), HasNoError())
}

func Test_Observable_Retry_Error_ShouldRetry(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Error(errFoo)
	}}).Retry(3, func(err error) bool {
		return true
	})
	Assert(ctx, t, obs, HasItems(1, 2, 1, 2, 1, 2, 1, 2), HasError(errFoo))
}

func Test_Observable_Retry_Error_ShouldNotRetry(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Error(errFoo)
	}}).Retry(3, func(err error) bool {
		return false
	})
	Assert(ctx, t, obs, HasItems(1, 2), HasError(errFoo))
}

func Test_Observable_Run(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := make([]int, 0)
	<-testObservable(ctx, 1, 2, 3).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		s = append(s, i.(int))
		return i, nil
	}).Run()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func Test_Observable_Run_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := make([]int, 0)
	<-testObservable(ctx, 1, errFoo).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		s = append(s, i.(int))
		return i, nil
	}).Run()
	assert.Equal(t, []int{1}, s)
}

func Test_Observable_Sample_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1).Sample(Empty(), WithContext(ctx))
	Assert(ctx, t, obs, IsEmpty(), HasNoError())
}

func Test_Observable_Scan(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).Scan(func(_ context.Context, x, y interface{}) (interface{}, error) {
		if x == nil {
			return y, nil
		}
		return x.(int) + y.(int), nil
	})
	Assert(ctx, t, obs, HasItems(1, 3, 6, 10, 15))
}

func Test_Observable_Scan_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).Scan(func(_ context.Context, x, y interface{}) (interface{}, error) {
		if x == nil {
			return y, nil
		}
		return x.(int) + y.(int), nil
	}, WithCPUPool())
	Assert(ctx, t, obs, HasItemsNoOrder(1, 3, 6, 10, 15))
}

func Test_Observable_SequenceEqual_EvenSequence(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sequence := testObservable(ctx, 2, 5, 12, 43, 98, 100, 213)
	result := testObservable(ctx, 2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequence)
	Assert(ctx, t, result, HasItem(true))
}

func Test_Observable_SequenceEqual_UnevenSequence(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sequence := testObservable(ctx, 2, 5, 12, 43, 98, 100, 213)
	result := testObservable(ctx, 2, 5, 12, 43, 15, 100, 213).SequenceEqual(sequence, WithContext(ctx))
	Assert(ctx, t, result, HasItem(false))
}

func Test_Observable_SequenceEqual_DifferentLengthSequence(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sequenceShorter := testObservable(ctx, 2, 5, 12, 43, 98, 100)
	sequenceLonger := testObservable(ctx, 2, 5, 12, 43, 98, 100, 213, 512)

	resultForShorter := testObservable(ctx, 2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequenceShorter)
	Assert(ctx, t, resultForShorter, HasItem(false))

	resultForLonger := testObservable(ctx, 2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequenceLonger)
	Assert(ctx, t, resultForLonger, HasItem(false))
}

func Test_Observable_SequenceEqual_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := Empty().SequenceEqual(Empty())
	Assert(ctx, t, result, HasItem(true))
}

func Test_Observable_Send(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan Item, 10)
	testObservable(ctx, 1, 2, 3, errFoo).Send(ch)
	assert.Equal(t, Of(1), <-ch)
	assert.Equal(t, Of(2), <-ch)
	assert.Equal(t, Of(3), <-ch)
	assert.Equal(t, Error(errFoo), <-ch)
}

type message struct {
	id int
}

func Test_Observable_Serialize_Struct(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, message{3}, message{5}, message{1}, message{2}, message{4}).
		Serialize(1, func(i interface{}) int {
			return i.(message).id
		})
	Assert(ctx, t, obs, HasItems(message{1}, message{2}, message{3}, message{4}, message{5}))
}

func Test_Observable_Serialize_Duplicates(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 3, 2, 6, 4, 5).
		Serialize(1, func(i interface{}) int {
			return i.(int)
		})
	Assert(ctx, t, obs, HasItems(1, 2, 3, 4, 5, 6))
}

func Test_Observable_Serialize_Loop(t *testing.T) {
	defer goleak.VerifyNone(t)
	idx := 0
	<-Range(1, 10000).
		Serialize(0, func(i interface{}) int {
			return i.(int)
		}).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			return i, nil
		}, WithCPUPool()).
		DoOnNext(func(i interface{}) {
			v := i.(int)
			if v != idx {
				assert.FailNow(t, "not sequential", "expected=%d, got=%d", idx, v)
			}
			idx++
		})
}

func Test_Observable_Serialize_DifferentFrom(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, message{13}, message{15}, message{11}, message{12}, message{14}).
		Serialize(11, func(i interface{}) int {
			return i.(message).id
		})
	Assert(ctx, t, obs, HasItems(message{11}, message{12}, message{13}, message{14}, message{15}))
}

func Test_Observable_Serialize_ContextCanceled(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	obs := Never().Serialize(1, func(i interface{}) int {
		return i.(message).id
	}, WithContext(ctx))
	Assert(ctx, t, obs, IsEmpty(), HasNoError())
}

func Test_Observable_Serialize_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, message{3}, message{5}, message{7}, message{2}, message{4}).
		Serialize(1, func(i interface{}) int {
			return i.(message).id
		})
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Observable_Serialize_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, message{3}, message{1}, errFoo, message{2}, message{4}).
		Serialize(1, func(i interface{}) int {
			return i.(message).id
		})
	Assert(ctx, t, obs, HasItems(message{1}), HasError(errFoo))
}

func Test_Observable_Skip(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 0, 1, 2, 3, 4, 5).Skip(3)
	Assert(ctx, t, obs, HasItems(3, 4, 5))
}

func Test_Observable_Skip_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 0, 1, 2, 3, 4, 5).Skip(3, WithCPUPool())
	Assert(ctx, t, obs, HasItems(3, 4, 5))
}

func Test_Observable_SkipLast(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 0, 1, 2, 3, 4, 5).SkipLast(3)
	Assert(ctx, t, obs, HasItems(0, 1, 2))
}

func Test_Observable_SkipLast_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 0, 1, 2, 3, 4, 5).SkipLast(3, WithCPUPool())
	Assert(ctx, t, obs, HasItems(0, 1, 2))
}

func Test_Observable_SkipWhile(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).SkipWhile(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i != 3
		default:
			return true
		}
	})

	Assert(ctx, t, obs, HasItems(3, 4, 5), HasNoError())
}

func Test_Observable_SkipWhile_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).SkipWhile(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i != 3
		default:
			return true
		}
	}, WithCPUPool())

	Assert(ctx, t, obs, HasItems(3, 4, 5), HasNoError())
}

func Test_Observable_StartWithIterable(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 4, 5, 6).StartWith(testObservable(ctx, 1, 2, 3))
	Assert(ctx, t, obs, HasItems(1, 2, 3, 4, 5, 6), HasNoError())
}

func Test_Observable_StartWithIterable_Error1(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 4, 5, 6).StartWith(testObservable(ctx, 1, errFoo, 3))
	Assert(ctx, t, obs, HasItems(1), HasError(errFoo))
}

func Test_Observable_StartWithIterable_Error2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 4, errFoo, 6).StartWith(testObservable(ctx, 1, 2, 3))
	Assert(ctx, t, obs, HasItems(1, 2, 3, 4), HasError(errFoo))
}

func Test_Observable_SumFloat32_OnlyFloat32(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float32(1.0), float32(2.0), float32(3.0)).SumFloat32(),
		HasItem(float32(6.)))
}

func Test_Observable_SumFloat32_DifferentTypes(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float32(1.1), 2, int8(3), int16(1), int32(1), int64(1)).SumFloat32(),
		HasItem(float32(9.1)))
}

func Test_Observable_SumFloat32_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).SumFloat32(), HasAnError())
}

func Test_Observable_SumFloat32_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().SumFloat32(), IsEmpty())
}

func Test_Observable_SumFloat64_OnlyFloat64(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).SumFloat64(),
		HasItem(6.6))
}

func Test_Observable_SumFloat64_DifferentTypes(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, float32(1.0), 2, int8(3), 4., int16(1), int32(1), int64(1)).SumFloat64(),
		HasItem(13.))
}

func Test_Observable_SumFloat64_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, "x").SumFloat64(), HasAnError())
}

func Test_Observable_SumFloat64_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().SumFloat64(), IsEmpty())
}

func Test_Observable_SumInt64_OnlyInt64(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1, 2, 3).SumInt64(), HasItem(int64(6)))
}

func Test_Observable_SumInt64_DifferentTypes(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, int8(1), int(2), int16(3), int32(4), int64(5)).SumInt64(),
		HasItem(int64(15)))
}

func Test_Observable_SumInt64_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, testObservable(ctx, 1.1, 2.2, 3.3).SumInt64(), HasAnError())
}

func Test_Observable_SumInt64_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Assert(ctx, t, Empty().SumInt64(), IsEmpty())
}

func Test_Observable_Take(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).Take(3)
	Assert(ctx, t, obs, HasItems(1, 2, 3))
}

func Test_Observable_Take_Interval(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Interval(WithDuration(time.Nanosecond), WithContext(ctx)).Take(3)
	Assert(ctx, t, obs, CustomPredicate(func(items []interface{}) error {
		if len(items) != 3 {
			return errors.New("3 items are expected")
		}
		return nil
	}))
}

func Test_Observable_TakeLast(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).TakeLast(3)
	Assert(ctx, t, obs, HasItems(3, 4, 5))
}

func Test_Observable_TakeLast_LessThanNth(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 4, 5).TakeLast(3)
	Assert(ctx, t, obs, HasItems(4, 5))
}

func Test_Observable_TakeLast_LessThanNth2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 4, 5).TakeLast(100000)
	Assert(ctx, t, obs, HasItems(4, 5))
}

func Test_Observable_TakeUntil(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).TakeUntil(func(item interface{}) bool {
		return item == 3
	})
	Assert(ctx, t, obs, HasItems(1, 2, 3))
}

func Test_Observable_TakeWhile(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3, 4, 5).TakeWhile(func(item interface{}) bool {
		return item != 3
	})
	Assert(ctx, t, obs, HasItems(1, 2))
}

func Test_Observable_TimeInterval(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).TimeInterval()
	Assert(ctx, t, obs, CustomPredicate(func(items []interface{}) error {
		if len(items) != 3 {
			return fmt.Errorf("expected 3 items, got %d items", len(items))
		}
		return nil
	}))
}

func Test_Observable_Timestamp(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observe := testObservable(ctx, 1, 2, 3).Timestamp().Observe()
	v := (<-observe).V.(TimestampItem)
	assert.Equal(t, 1, v.V)
	v = (<-observe).V.(TimestampItem)
	assert.Equal(t, 2, v.V)
	v = (<-observe).V.(TimestampItem)
	assert.Equal(t, 3, v.V)
}

func Test_Observable_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observe := testObservable(ctx, 1, errFoo).Timestamp().Observe()
	v := (<-observe).V.(TimestampItem)
	assert.Equal(t, 1, v.V)
	assert.True(t, (<-observe).Error())
}

func Test_Observable_ToMap(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 3, 4, 5, true, false).ToMap(func(_ context.Context, i interface{}) (interface{}, error) {
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
	Assert(ctx, t, obs, HasItem(map[interface{}]interface{}{
		3: 3,
		4: 4,
		5: 5,
		0: true,
		1: false,
	}))
}

func Test_Observable_ToMapWithValueSelector(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	keySelector := func(_ context.Context, i interface{}) (interface{}, error) {
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
	valueSelector := func(_ context.Context, i interface{}) (interface{}, error) {
		switch v := i.(type) {
		case int:
			return v * 10, nil
		case bool:
			return v, nil
		default:
			return i, nil
		}
	}
	single := testObservable(ctx, 3, 4, 5, true, false).ToMapWithValueSelector(keySelector, valueSelector)
	Assert(ctx, t, single, HasItem(map[interface{}]interface{}{
		3: 30,
		4: 40,
		5: 50,
		0: true,
		1: false,
	}))
}

func Test_Observable_ToSlice(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s, err := testObservable(ctx, 1, 2, 3).ToSlice(5)
	assert.Equal(t, []interface{}{1, 2, 3}, s)
	assert.Equal(t, 5, cap(s))
	assert.NoError(t, err)
}

func Test_Observable_ToSlice_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s, err := testObservable(ctx, 1, 2, errFoo, 3).ToSlice(0)
	assert.Equal(t, []interface{}{1, 2}, s)
	assert.Equal(t, errFoo, err)
}

func Test_Observable_Unmarshal(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, []byte(`{"id":1}`), []byte(`{"id":2}`)).Unmarshal(json.Unmarshal,
		func() interface{} {
			return &testStruct{}
		})
	Assert(ctx, t, obs, HasItems(&testStruct{
		ID: 1,
	}, &testStruct{
		ID: 2,
	}))
}

func Test_Observable_Unmarshal_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, []byte(`{"id":1`), []byte(`{"id":2}`)).Unmarshal(json.Unmarshal,
		func() interface{} {
			return &testStruct{}
		})
	Assert(ctx, t, obs, HasAnError())
}

func Test_Observable_Unmarshal_Parallel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, []byte(`{"id":1}`), []byte(`{"id":2}`)).Unmarshal(json.Unmarshal,
		func() interface{} {
			return &testStruct{}
		}, WithPool(1))
	Assert(ctx, t, obs, HasItems(&testStruct{
		ID: 1,
	}, &testStruct{
		ID: 2,
	}))
}

func Test_Observable_Unmarshal_Parallel_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, []byte(`{"id":1`), []byte(`{"id":2}`)).Unmarshal(json.Unmarshal,
		func() interface{} {
			return &testStruct{}
		}, WithCPUPool())
	Assert(ctx, t, obs, HasAnError())
}

func Test_Observable_WindowWithCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observe := testObservable(ctx, 1, 2, 3, 4, 5).WindowWithCount(2).Observe()
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(1, 2))
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(3, 4))
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(5))
}

func Test_Observable_WindowWithCount_ZeroCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observe := testObservable(ctx, 1, 2, 3, 4, 5).WindowWithCount(0).Observe()
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(1, 2, 3, 4, 5))
}

func Test_Observable_WindowWithCount_ObservableError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observe := testObservable(ctx, 1, 2, errFoo, 4, 5).WindowWithCount(2).Observe()
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(1, 2))
	Assert(ctx, t, (<-observe).V.(Observable), IsEmpty(), HasError(errFoo))
}

func Test_Observable_WindowWithCount_InputError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Empty().WindowWithCount(-1)
	Assert(ctx, t, obs, HasAnError())
}

func Test_Observable_WindowWithTime(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan Item, 10)
	ch <- Of(1)
	ch <- Of(2)
	obs := FromChannel(ch)
	go func() {
		time.Sleep(30 * time.Millisecond)
		ch <- Of(3)
		close(ch)
	}()

	observe := obs.WindowWithTime(WithDuration(10*time.Millisecond), WithBufferedChannel(10)).Observe()
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(1, 2))
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(3))
}

func Test_Observable_WindowWithTimeOrCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan Item, 10)
	ch <- Of(1)
	ch <- Of(2)
	obs := FromChannel(ch)
	go func() {
		time.Sleep(30 * time.Millisecond)
		ch <- Of(3)
		close(ch)
	}()

	observe := obs.WindowWithTimeOrCount(WithDuration(10*time.Millisecond), 1, WithBufferedChannel(10)).Observe()
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(1))
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(2))
	Assert(ctx, t, (<-observe).V.(Observable), HasItems(3))
}

func Test_Observable_ZipFromObservable(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs1 := testObservable(ctx, 1, 2, 3)
	obs2 := testObservable(ctx, 10, 20, 30)
	zipper := func(_ context.Context, elem1, elem2 interface{}) (interface{}, error) {
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
	Assert(ctx, t, zip, HasItems(11, 22, 33))
}

func Test_Observable_ZipFromObservable_DifferentLength1(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs1 := testObservable(ctx, 1, 2, 3)
	obs2 := testObservable(ctx, 10, 20)
	zipper := func(_ context.Context, elem1, elem2 interface{}) (interface{}, error) {
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
	Assert(ctx, t, zip, HasItems(11, 22))
}

func Test_Observable_ZipFromObservable_DifferentLength2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs1 := testObservable(ctx, 1, 2)
	obs2 := testObservable(ctx, 10, 20, 30)
	zipper := func(_ context.Context, elem1, elem2 interface{}) (interface{}, error) {
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
	Assert(ctx, t, zip, HasItems(11, 22))
}
