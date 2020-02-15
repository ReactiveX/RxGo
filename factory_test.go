package rxgo

import (
	"context"
	"testing"
)

func Test_Empty(t *testing.T) {
	obs := Empty()
	AssertObservable(context.Background(), t, obs, HasNoItems())
}

func Test_FromChannel(t *testing.T) {
	obs := testObservable(1, 2, 3)
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFunc(t *testing.T) {
	obs := FromFunc(func(ctx context.Context, next chan<- Item) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		close(next)
	})
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFunc_Close(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	obs := FromFunc(func(ctx context.Context, next chan<- Item) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		cancel()
	})
	AssertObservable(ctx, t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFunc_Dup(t *testing.T) {
	obs := FromFunc(func(ctx context.Context, next chan<- Item) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		close(next)
	})
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromFunc_Error(t *testing.T) {
	obs := FromFunc(func(ctx context.Context, next chan<- Item) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromError(errFoo)
	})
	AssertObservable(context.Background(), t, obs, HasItems(1, 2), HasRaisedError(errFoo))
}

func Test_FromItem(t *testing.T) {
	single := FromItem(FromValue(1))
	AssertSingle(context.Background(), t, single, HasItem(1), HasNotRaisedError())
	AssertSingle(context.Background(), t, single, HasItem(1), HasNotRaisedError())
}

func Test_FromItems(t *testing.T) {
	obs := FromItems(FromValue(1), FromValue(2), FromValue(3))
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
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
