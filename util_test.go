package rxgo

import (
	"errors"
)

type testStruct struct {
	ID int `json:"id"`
}

var errFoo = errors.New("foo")

func channelValue(items ...interface{}) chan Item {
	next := make(chan Item)
	go func() {
		for _, item := range items {
			switch item := item.(type) {
			default:
				next <- FromValue(item)
			case error:
				next <- FromError(item)
				close(next)
				return
			}
		}
		close(next)
	}()
	return next
}

func testObservable(items ...interface{}) Observable {
	return FromChannel(channelValue(items...))
}
