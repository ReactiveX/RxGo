package rxgo

import (
	"context"
	"testing"
)

func Test_Empty(t *testing.T) {
	obs := Empty()
	Assert(context.Background(), t, obs, HasNoItems())
}

func Test_FromChannel(t *testing.T) {
	obs := testObservable(1, 2, 3)
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
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

func Test_FromFuncs_Dup(t *testing.T) {
	obs := FromFuncs(func(ctx context.Context, next chan<- Item, done func()) {
		next <- FromValue(1)
		next <- FromValue(2)
		next <- FromValue(3)
		done()
	})
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
	Assert(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
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
