package rxgo

import (
	"github.com/reactivex/rxgo/v2/handlers"
	"github.com/reactivex/rxgo/v2/optional"
	"github.com/reactivex/rxgo/v2/options"
)

// Single is similar to an Observable but emits only one single element or an error notification.
type Single interface {
	Filter(apply Predicate) OptionalSingle
	Map(apply Function) Single
	Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver
}

type OptionalSingle interface {
	Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver
}

type single struct {
	ch chan interface{}
}

type optionalSingle struct {
	ch chan optional.Optional
}

func newSingleFrom(item interface{}) Single {
	s := single{
		ch: make(chan interface{}),
	}

	go func() {
		s.ch <- item
		close(s.ch)
	}()

	return &s
}

func newOptionalSingleFrom(opt optional.Optional) OptionalSingle {
	s := optionalSingle{
		ch: make(chan optional.Optional),
	}

	go func() {
		s.ch <- opt
		close(s.ch)
	}()

	return &s
}

// CheckHandler checks the underlying type of an EventHandler.
func CheckSingleEventHandler(handler handlers.EventHandler) SingleObserver {
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

func NewOptionalSingleFromChannel(ch chan optional.Optional) OptionalSingle {
	return &optionalSingle{
		ch: ch,
	}
}

func (s *single) Filter(apply Predicate) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		item := <-s.ch
		if apply(item) {
			out <- optional.Of(item)
		} else {
			out <- optional.Empty()
		}
		close(out)
		return
	}()

	return &optionalSingle{
		ch: out,
	}
}

func (s *single) Map(apply Function) Single {
	out := make(chan interface{})
	go func() {
		item := <-s.ch
		out <- apply(item)
		close(out)
	}()
	return &single{ch: out}
}

func (s *single) Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver {
	ob := CheckSingleEventHandler(handler)

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

func (s *optionalSingle) Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver {
	ob := CheckSingleEventHandler(handler)

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
