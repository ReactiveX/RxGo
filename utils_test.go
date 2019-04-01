package rxgo

import (
	"errors"
	"github.com/reactivex/rxgo/handlers"
	"time"
)

const timeout = 500 * time.Millisecond
const pollingInterval = 20 * time.Millisecond

var noData = errors.New("timeout")

func get(ch chan interface{}, d time.Duration) interface{} {
	select {
	case res := <-ch:
		return res
	case <-time.After(d):
		return noData
	}
}

func next(out chan interface{}) handlers.NextFunc {
	return handlers.NextFunc(func(i interface{}) {
		out <- i
	})
}
