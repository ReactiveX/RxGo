package rxgo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func collect(ctx context.Context, ch <-chan Item) ([]interface{}, error) {
	s := make([]interface{}, 0)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case item, ok := <-ch:
			if !ok {
				return s, nil
			}
			if item.Error() {
				s = append(s, item.E)
			} else {
				s = append(s, item.V)
			}
		}
	}
}

func Test_Amb1(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Amb([]Observable{testObservable(ctx, 1, 2, 3), Empty()})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_Amb2(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Amb([]Observable{Empty(), testObservable(ctx, 1, 2, 3), Empty(), Empty()})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_CombineLatest(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := CombineLatest(func(ii ...interface{}) interface{} {
		sum := 0
		for _, v := range ii {
			if v == nil {
				continue
			}
			sum += v.(int)
		}
		return sum
	}, []Observable{testObservable(ctx, 1, 2), testObservable(ctx, 10, 11)})
	Assert(context.Background(), t, obs, IsNotEmpty())
}

func Test_CombineLatest_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := CombineLatest(func(ii ...interface{}) interface{} {
		sum := 0
		for _, v := range ii {
			sum += v.(int)
		}
		return sum
	}, []Observable{testObservable(ctx, 1, 2), Empty()})
	Assert(context.Background(), t, obs, IsEmpty())
}

func Test_CombineLatest_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := CombineLatest(func(ii ...interface{}) interface{} {
		sum := 0
		for _, v := range ii {
			sum += v.(int)
		}
		return sum
	}, []Observable{testObservable(ctx, 1, 2), testObservable(ctx, errFoo)})
	Assert(context.Background(), t, obs, IsEmpty(), HasError(errFoo))
}

func Test_Concat_SingleObservable(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Concat([]Observable{testObservable(ctx, 1, 2, 3)})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_Concat_TwoObservables(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Concat([]Observable{testObservable(ctx, 1, 2, 3), testObservable(ctx, 4, 5, 6)})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3, 4, 5, 6))
}

func Test_Concat_MoreThanTwoObservables(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Concat([]Observable{testObservable(ctx, 1, 2, 3), testObservable(ctx, 4, 5, 6), testObservable(ctx, 7, 8, 9)})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3, 4, 5, 6, 7, 8, 9))
}

func Test_Concat_EmptyObservables(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Concat([]Observable{Empty(), Empty(), Empty()})
	Assert(context.Background(), t, obs, IsEmpty())
}

func Test_Concat_OneEmptyObservable(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Concat([]Observable{Empty(), testObservable(ctx, 1, 2, 3)})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))

	obs = Concat([]Observable{testObservable(ctx, 1, 2, 3), Empty()})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_Create(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Create([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Of(3)
	}})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
}

func Test_Create_SingleDup(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Create([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Of(3)
	}})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
	Assert(context.Background(), t, obs, IsEmpty(), HasNoError())
}

func Test_Create_ContextCancelled(t *testing.T) {
	defer goleak.VerifyNone(t)
	closed1 := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	Create([]Producer{
		func(ctx context.Context, next chan<- Item) {
			cancel()
		}, func(ctx context.Context, next chan<- Item) {
			<-ctx.Done()
			closed1 <- struct{}{}
		},
	}, WithContext(ctx)).Run()

	select {
	case <-time.Tick(time.Second):
		assert.FailNow(t, "producer not closed")
	case <-closed1:
	}
}

func Test_Defer(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Of(3)
	}})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
}

func Test_Defer_Multiple(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
	}, func(ctx context.Context, next chan<- Item) {
		next <- Of(10)
		next <- Of(20)
	}})
	Assert(context.Background(), t, obs, HasItemsNoOrder(1, 2, 10, 20), HasNoError())
}

func Test_Defer_ContextCancelled(t *testing.T) {
	defer goleak.VerifyNone(t)
	closed1 := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	Defer([]Producer{
		func(ctx context.Context, next chan<- Item) {
			cancel()
		}, func(ctx context.Context, next chan<- Item) {
			<-ctx.Done()
			closed1 <- struct{}{}
		},
	}, WithContext(ctx)).Run()

	select {
	case <-time.Tick(time.Second):
		assert.FailNow(t, "producer not closed")
	case <-closed1:
	}
}

func Test_Defer_SingleDup(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Of(3)
	}})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
}

func Test_Defer_ComposedDup(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Of(3)
	}}).Map(func(_ context.Context, i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	}).Map(func(_ context.Context, i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNoError())
	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNoError())
}

func Test_Defer_ComposedDup_EagerObservation(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Of(3)
	}}).Map(func(_ context.Context, i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	}, WithObservationStrategy(Eager)).Map(func(_ context.Context, i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNoError())
	// In the case of an eager observation, we already consumed the items produced by Defer
	// So if we create another subscription, it will be empty
	Assert(context.Background(), t, obs, IsEmpty(), HasNoError())
}

func Test_Defer_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Defer([]Producer{func(ctx context.Context, next chan<- Item) {
		next <- Of(1)
		next <- Of(2)
		next <- Error(errFoo)
	}})
	Assert(context.Background(), t, obs, HasItems(1, 2), HasError(errFoo))
}

func Test_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Empty()
	Assert(context.Background(), t, obs, IsEmpty())
}

func Test_FromChannel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	obs := FromChannel(ch)
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
}

func Test_FromChannel_SimpleCapacity(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := FromChannel(make(chan Item, 10)).Observe()
	assert.Equal(t, 10, cap(ch))
}

func Test_FromChannel_ComposedCapacity(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	obs1 := FromChannel(make(chan Item, 10)).
		Map(func(_ context.Context, _ interface{}) (interface{}, error) {
			return 1, nil
		}, WithContext(ctx), WithBufferedChannel(11))
	assert.Equal(t, 11, cap(obs1.Observe()))

	obs2 := obs1.Map(func(_ context.Context, _ interface{}) (interface{}, error) {
		return 1, nil
	}, WithContext(ctx), WithBufferedChannel(12))
	assert.Equal(t, 12, cap(obs2.Observe()))
}

func Test_FromEventSource_ObservationAfterAllSent(t *testing.T) {
	defer goleak.VerifyNone(t)
	const max = 10
	next := make(chan Item, max)
	obs := FromEventSource(next, WithBackPressureStrategy(Drop))

	go func() {
		for i := 0; i < max; i++ {
			next <- Of(i)
		}
		close(next)
	}()
	time.Sleep(50 * time.Millisecond)

	Assert(context.Background(), t, obs, CustomPredicate(func(items []interface{}) error {
		if len(items) != 0 {
			return errors.New("items should be nil")
		}
		return nil
	}))
}

func Test_FromEventSource_Drop(t *testing.T) {
	defer goleak.VerifyNone(t)
	const max = 100000
	next := make(chan Item, max)
	obs := FromEventSource(next, WithBackPressureStrategy(Drop))

	go func() {
		for i := 0; i < max; i++ {
			next <- Of(i)
		}
		close(next)
	}()

	Assert(context.Background(), t, obs, CustomPredicate(func(items []interface{}) error {
		if len(items) == max {
			return errors.New("some items should be dropped")
		}
		if len(items) == 0 {
			return errors.New("no items")
		}
		return nil
	}))
}

func Test_Interval(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	obs := Interval(WithDuration(time.Nanosecond), WithContext(ctx))
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	Assert(context.Background(), t, obs, IsNotEmpty())
}

func Test_JustItem(t *testing.T) {
	defer goleak.VerifyNone(t)
	single := JustItem(1)
	Assert(context.Background(), t, single, HasItem(1), HasNoError())
	Assert(context.Background(), t, single, HasItem(1), HasNoError())
}

func Test_Just(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Just(1, 2, 3)()
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNoError())
}

func Test_Just_CustomStructure(t *testing.T) {
	defer goleak.VerifyNone(t)
	type customer struct {
		id int
	}

	obs := Just(customer{id: 1}, customer{id: 2}, customer{id: 3})()
	Assert(context.Background(), t, obs, HasItems(customer{id: 1}, customer{id: 2}, customer{id: 3}), HasNoError())
	Assert(context.Background(), t, obs, HasItems(customer{id: 1}, customer{id: 2}, customer{id: 3}), HasNoError())
}

func Test_Just_Channel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan int, 1)
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
	}()
	obs := Just(ch)()
	Assert(context.Background(), t, obs, HasItems(1, 2, 3))
}

func Test_Just_SimpleCapacity(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := Just(1)(WithBufferedChannel(5)).Observe()
	assert.Equal(t, 5, cap(ch))
}

func Test_Just_ComposedCapacity(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs1 := Just(1)().Map(func(_ context.Context, _ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(11))
	assert.Equal(t, 11, cap(obs1.Observe()))

	obs2 := obs1.Map(func(_ context.Context, _ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(12))
	assert.Equal(t, 12, cap(obs2.Observe()))
}

func Test_Merge(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Merge([]Observable{testObservable(ctx, 1, 2), testObservable(ctx, 3, 4)})
	Assert(context.Background(), t, obs, HasItemsNoOrder(1, 2, 3, 4))
}

func Test_Merge_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := Merge([]Observable{testObservable(ctx, 1, 2), testObservable(ctx, 3, errFoo)})
	// The content is not deterministic, hence we just test if we have some items
	Assert(context.Background(), t, obs, IsNotEmpty(), HasError(errFoo))
}

func Test_Merge_Interval(t *testing.T) {
	defer goleak.VerifyNone(t)
	var obs []Observable
	ctx, cancel := context.WithCancel(context.Background())
	obs = append(obs, Interval(WithDuration(3*time.Millisecond), WithContext(ctx)).
		Take(3).
		Map(func(_ context.Context, v interface{}) (interface{}, error) {
			return 10 + v.(int), nil
		}))
	obs = append(obs, Interval(WithDuration(5*time.Millisecond), WithContext(ctx)).
		Take(3).
		Map(func(_ context.Context, v interface{}) (interface{}, error) {
			return 20 + v.(int), nil
		}))

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	Assert(ctx, t, Merge(obs), HasNoError(), HasItemsNoOrder(10, 11, 12, 20, 21, 22))
}

func Test_Range(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Range(5, 3)
	Assert(context.Background(), t, obs, HasItems(5, 6, 7, 8))
	// Test whether the observable is reproducible
	Assert(context.Background(), t, obs, HasItems(5, 6, 7, 8))
}

func Test_Range_NegativeCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Range(1, -5)
	Assert(context.Background(), t, obs, HasAnError())
}

func Test_Range_MaximumExceeded(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Range(1<<31, 1)
	Assert(context.Background(), t, obs, HasAnError())
}

func Test_Start(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Start([]Supplier{func(ctx context.Context) Item {
		return Of(1)
	}, func(ctx context.Context) Item {
		return Of(2)
	}})
	Assert(context.Background(), t, obs, HasItemsNoOrder(1, 2))
}

func Test_Thrown(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Thrown(errFoo)
	Assert(context.Background(), t, obs, HasError(errFoo))
}

func Test_Timer(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs := Timer(WithDuration(time.Nanosecond))
	select {
	case <-time.Tick(time.Second):
		assert.FailNow(t, "observable not closed")
	case <-obs.Observe():
	}
}

func Test_Timer_Empty(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	obs := Timer(WithDuration(time.Hour), WithContext(ctx))
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	select {
	case <-time.Tick(time.Second):
		assert.FailNow(t, "observable not closed")
	case <-obs.Observe():
	}
}
