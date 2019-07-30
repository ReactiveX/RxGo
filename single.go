package rxgo

import (
	"context"
)

// Single is similar to an Observable but emits only one single element or an error notification.
type Single interface {
	Iterable
	Filter(apply Predicate) OptionalSingle
	Map(apply Function) Single
	Subscribe(handler EventHandler, opts ...Option) SingleObserver
}

type OptionalSingle interface {
	Subscribe(handler EventHandler, opts ...Option) SingleObserver
}

type single struct {
	iterable Iterable
}

type optionalSingle struct {
	ch chan Optional
}

func newSingleFrom(item interface{}) Single {
	f := func(out chan interface{}) {
		out <- item
		close(out)
	}
	return newColdSingle(f)
}

func newOptionalSingleFrom(opt Optional) OptionalSingle {
	s := optionalSingle{
		ch: make(chan Optional),
	}

	go func() {
		s.ch <- opt
		close(s.ch)
	}()

	return &s
}

// CheckSingleEventHandler checks the underlying type of an EventHandler.
func CheckSingleEventHandler(handler EventHandler) SingleObserver {
	return NewSingleObserver(handler)
}

func newColdSingle(f func(chan interface{})) Single {
	return &single{
		iterable: newIterableFromFunc(f),
	}
}

// NewOptionalSingleFromChannel creates a new OptionalSingle from a channel input
func NewOptionalSingleFromChannel(ch chan Optional) OptionalSingle {
	return &optionalSingle{
		ch: ch,
	}
}

func (s *single) Iterator(ctx context.Context) Iterator {
	return s.iterable.Iterator(context.Background())
}

func (s *single) Filter(apply Predicate) OptionalSingle {
	out := make(chan Optional)
	go func() {
		it := s.iterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if apply(item) {
					out <- Of(item)
				} else {
					out <- EmptyOptional()
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
		it := s.iterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

func (s *single) Subscribe(handler EventHandler, opts ...Option) SingleObserver {
	ob := CheckSingleEventHandler(handler)

	go func() {
		it := s.iterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

func (s *optionalSingle) Subscribe(handler EventHandler, opts ...Option) SingleObserver {
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
