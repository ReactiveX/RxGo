package rxgo

import (
	"context"
	"testing"
)

func Test_OptionalSingle_Observe(t *testing.T) {
	os := FromItem(FromValue(1)).Filter(context.Background(), func(i interface{}) bool {
		return i == 1
	})
	Assert(context.Background(), t, os, HasItem(1), HasNotRaisedError())
}
