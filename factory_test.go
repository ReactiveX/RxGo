package rxgo

import (
	"context"
	"testing"
)

func Test_FromChannel(t *testing.T) {
	next := channelValue(1, 2, 3, closeCmd)
	obs := FromChannel(next)
	assertObservable(t, context.Background(), obs, hasItems(1, 2, 3), hasNotRaisedError())
}
