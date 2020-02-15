package rxgo

import (
	"errors"
)

var fooErr = errors.New("foo")
var closeCmd = &struct{}{}

func channelValue(items ...interface{}) chan Item {
	next := make(chan Item)
	go func() {
		for _, item := range items {
			switch item := item.(type) {
			default:
				if item == closeCmd {
					close(next)
					return
				}
				next <- FromValue(item)
			case error:
				next <- FromError(item)
				close(next)
				return
			}
		}
	}()
	return next
}
