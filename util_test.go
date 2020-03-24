package rxgo

import (
	"errors"
)

type testStruct struct {
	ID int `json:"id"`
}

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")
)

func channelValue(items ...interface{}) chan Item {
	next := make(chan Item)
	go func() {
		for _, item := range items {
			switch item := item.(type) {
			default:
				next <- Of(item)
			case error:
				next <- Error(item)
			}
		}
		close(next)
	}()
	return next
}

func testObservable(items ...interface{}) Observable {
	return FromChannel(channelValue(items...))
}
