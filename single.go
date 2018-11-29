package rx

import (
	"github.com/reactivex/rxgo/fx"
)

// Single is similar to an Observable but emits only one single element or an error notification.
type Single interface {
	Filter(apply fx.Predicate) Single
	Map(apply fx.Function) Single
	Subscribe(handler EventHandler, opts ...Option) SingleObserver
}

type single struct {
	ch chan interface{}
}

// CheckHandler checks the underlying type of an EventHandler.
func CheckSingleEventHandler(handler EventHandler) SingleObserver {
	return NewSingleObserver(handler)
}

func NewSingle() Single {
	return &single{
		ch: make(chan interface{}),
	}
}

func NewSingleFromChannel(ch chan interface{}) Single {
	return &single{
		ch: ch,
	}
}

func (s *single) Filter(apply fx.Predicate) Single {
	out := make(chan interface{})
	go func() {
		item := <-s.ch
		if apply(item) {
			out <- item
		}
		close(out)
	}()
	return &single{ch: out}
}

func (s *single) Map(apply fx.Function) Single {
	out := make(chan interface{})
	go func() {
		item := <-s.ch
		out <- apply(item)
		close(out)
	}()
	return &single{ch: out}
}

func (s *single) Subscribe(handler EventHandler, opts ...Option) SingleObserver {
	ob := CheckSingleEventHandler(handler)

	// Parse options
	var observableOptions options
	for _, opt := range opts {
		opt.apply(&observableOptions)
	}

	go func() {
		for item := range s.ch {
			switch item := item.(type) {
			case error:
				ob.OnError(item)

				// Record the error and break the loop.
				return
			default:
				ob.OnSuccess(item)
			}
		}

		return
	}()

	return ob
}
