package rxgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Empty(t *testing.T) {
	obs := Empty()
	Assert(context.Background(), t, obs, HasNoItems())
}

func Test_FromChannel(t *testing.T) {
	ch := make(chan Item)
	go func() {
		ch <- FromValue(1)
		ch <- FromValue(2)
		ch <- FromValue(3)
		close(ch)
	}()
	obs := FromChannel(ch)
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromChannel_SimpleCapacity(t *testing.T) {
	ch := FromChannel(make(chan Item, 10)).Observe(WithBufferedChannel(11))
	assert.Equal(t, 10, cap(ch))
}

func Test_FromChannel_ComposedCapacity(t *testing.T) {
	obs1 := FromChannel(make(chan Item, 10)).
		Map(func(_ interface{}) (interface{}, error) {
			return 1, nil
		}, WithBufferedChannel(11))
	assert.Equal(t, 11, cap(obs1.Observe(WithBufferedChannel(13))))

	obs2 := obs1.Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(12))
	assert.Equal(t, 12, cap(obs2.Observe(WithBufferedChannel(13))))
}

func Test_FromFuncs(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		done()
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFuncs_Multiple(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		done()
	}, func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(10)
		next <- FromValue(20)
		done()
	})
	Assert(context.Background(), t, obs, HasItemsNoParticularOrder(1, 2, 10, 20), HasNotRaisedError())
}

func Test_FromFuncs_Close(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		done()
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFuncs_SingleDup(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		done()
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFuncs_ComposedDup(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		done()
	}).Map(func(i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	}).Map(func(i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNotRaisedError())
	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNotRaisedError())
}

func Test_FromFuncs_ComposedDup_EagerObservation(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		done()
	}).Map(func(i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	}, WithEagerObservation()).Map(func(i interface{}) (_ interface{}, _ error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, obs, HasItems(3, 4, 5), HasNotRaisedError())
	// In the case of an eager observation, we already consumed the items produced by FromFuncs
	// So if we create another subscription, it will be empty
	Assert(context.Background(), t, obs, HasNoItem(), HasNotRaisedError())
}

func Test_FromFuncs_Error(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromError(errFoo)
		done()
	})
	Assert(context.Background(), t, obs, HasItems(1, 2), HasRaisedError(errFoo))
}

func Test_FromFuncs_SimpleCapacity(t *testing.T) {
	ch := FromFuncs(func(_ context.Context, _ chan<- Item, done func()) {
		done()
	}).Observe(WithBufferedChannel(5))
	assert.Equal(t, 5, cap(ch))
}

func Test_FromFuncs_ComposedCapacity(t *testing.T) {
	obs1 := FromFuncs(func(_ context.Context, _ chan<- Item, done func()) {
		done()
	}).Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(11))
	assert.Equal(t, 11, cap(obs1.Observe(WithBufferedChannel(13))))

	obs2 := obs1.Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(12))
	assert.Equal(t, 12, cap(obs2.Observe(WithBufferedChannel(13))))
}

func Test_FromItem(t *testing.T) {
	single := FromItem(FromValue(1))
	Assert(context.Background(), t, single, HasItem(1), HasNotRaisedError())
	Assert(context.Background(), t, single, HasItem(1), HasNotRaisedError())
}

func Test_FromItems(t *testing.T) {
	obs := FromItems(FromValue(1), FromValue(2), FromValue(3))
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromItems_SimpleCapacity(t *testing.T) {
	ch := FromItems(FromValue(1)).Observe(WithBufferedChannel(5))
	assert.Equal(t, 5, cap(ch))
}

func Test_FromItems_ComposedCapacity(t *testing.T) {
	obs1 := FromItems(FromValue(1)).Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(11))
	assert.Equal(t, 11, cap(obs1.Observe(WithBufferedChannel(13))))

	obs2 := obs1.Map(func(_ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(12))
	assert.Equal(t, 12, cap(obs2.Observe(WithBufferedChannel(13))))
}

// TODO Sleep
// func Test_FromEventSource(t *testing.T) {
//	next := make(chan Item, 100)
//	ctx, cancel := context.WithCancel(context.Background())
//	obs := FromEventSource(ctx, next, Block)
//	next <- FromValue(1)
//	time.Sleep(50 * time.Millisecond)
//
//	ch := obs.Observe()
//	go func() {
//		for i := range ch {
//		}
//	}()
//	next <- FromValue(2)
//	time.Sleep(50 * time.Millisecond)
//	next <- FromValue(3)
//	time.Sleep(50 * time.Millisecond)
//	cancel()
//}
