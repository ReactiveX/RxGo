package rxgo

import (
	"errors"
	"github.com/reactivex/rxgo/handlers"
	"time"
)

const timeout = 500 * time.Millisecond
const pollingInterval = 20 * time.Millisecond

var noData = errors.New("timeout")
var doneSignal = "done"

func pollItem(ch chan interface{}, d time.Duration) interface{} {
	select {
	case res := <-ch:
		return res
	case <-time.After(d):
		return noData
	}
}

func pollItems(ch chan interface{}, d time.Duration) []interface{} {
	s := make([]interface{}, 0)
	for {
		select {
		case res := <-ch:
			s = append(s, res)
		case <-time.After(d):
			return s
		}
	}
}

func nextHandler(out chan interface{}) handlers.NextFunc {
	return handlers.NextFunc(func(i interface{}) {
		out <- i
	})
}

func doneHandler(out chan interface{}) handlers.DoneFunc {
	return handlers.DoneFunc(func() {
		out <- doneSignal
	})
}

func errorHandler(out chan interface{}) handlers.ErrFunc {
	return handlers.ErrFunc(func(err error) {
		out <- err
	})
}
