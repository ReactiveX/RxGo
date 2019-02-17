package rxgo

import (
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/optional"
	"github.com/reactivex/rxgo/options"
)

// Single is similar to an Observable but emits only one single element or an error notification.
type Single interface {
	Iterable
	Filter(apply Predicate) OptionalSingle
	Map(apply Function) Single
	Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver
}

type OptionalSingle interface {
	Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver
}

type single struct {
	iterable Iterable
}

type optionalSingle struct {
	ch chan optional.Optional
}

func newSingleFrom(item interface{}) Single {
	f := func(out chan interface{}) {
		out <- item
		close(out)
	}
	return newColdSingle(f)
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

func newColdSingle(f func(chan interface{})) Single {
	return &single{
		iterable: newIterableFromFunc(f),
	}
}

func NewOptionalSingleFromChannel(ch chan optional.Optional) OptionalSingle {
	return &optionalSingle{
		ch: ch,
	}
}

func (s *single) Iterator() Iterator {
	return s.iterable.Iterator()
}

func (s *single) Filter(apply Predicate) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		it := s.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if apply(item) {
					out <- optional.Of(item)
				} else {
					out <- optional.Empty()
				}
				close(out)
				return
			} else {
				break
			}
		}
	}()

	return &optionalSingle{
		ch: out,
	}
}

func (s *single) Map(apply Function) Single {
	f := func(out chan interface{}) {
		it := s.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- apply(item)
				close(out)
				return
			} else {
				break
			}
		}
	}
	return newColdSingle(f)
}

func (s *single) Subscribe(handler handlers.EventHandler, opts ...options.Option) SingleObserver {
	ob := CheckSingleEventHandler(handler)

	go func() {
		it := s.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case error:
					ob.OnError(item)

					// Record the error and break the loop.
					return
				default:
					ob.OnSuccess(item)
				}
			} else {
				break
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
