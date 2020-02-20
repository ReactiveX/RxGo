package rxgo

import (
	"context"
	"testing"
)

func Test_Single_Filter_True(t *testing.T) {
	os := JustItem(1).Filter(func(i interface{}) bool {
		return i == 1
	})
	Assert(context.Background(), t, os, HasItem(1), HasNotRaisedError())
}

func Test_Single_Filter_False(t *testing.T) {
	os := JustItem(1).Filter(func(i interface{}) bool {
		return i == 0
	})
	Assert(context.Background(), t, os, HasNoItem(), HasNotRaisedError())
}

func Test_Single_Map(t *testing.T) {
	single := JustItem(1).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, single, HasItem(2), HasNotRaisedError())
}
