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
	next := channelValue(1, 2, 3, closeCmd)
	obs := FromChannel(next)
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}

func Test_FromItem(t *testing.T) {
	single := FromItem(FromValue(1))
	AssertSingle(context.Background(), t, single, HasValue(1), HasNotRaisedError())
	AssertSingle(context.Background(), t, single, HasValue(1), HasNotRaisedError())
}

func Test_FromItems(t *testing.T) {
	obs := FromItems(FromValue(1), FromValue(2), FromValue(3))
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
	AssertObservable(context.Background(), t, obs, HasItems(1, 2, 3), HasNotRaisedError())
}
