package rxgo

import (
	"context"
)

// Single is similar to an Observable but emits only one single element or an error notification.
type Single interface {
	Iterable
	Filter(apply Predicate) OptionalSingle
	Map(apply Function) Single
	Subscribe(handler EventHandler, opts ...Option) Observer
}

type OptionalSingle interface {
	Subscribe(handler EventHandler, opts ...Option) Observer
}

type single struct {
	iterable Iterable
}

type optionalSingle struct {
	itemChannel chan Optional
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
		itemChannel: make(chan Optional),
	}

	go func() {
		s.itemChannel <- opt
		close(s.itemChannel)
	}()

	return &s
}

// CheckSingleEventHandler checks the underlying type of an EventHandler.
func CheckSingleEventHandler(handler EventHandler) Observer {
	return NewObserver(handler)
}

func newColdSingle(f func(chan interface{})) Single {
	return &single{
		iterable: newIterableFromFunc(f),
	}
}

// NewOptionalSingleFromChannel creates a new OptionalSingle from a channel input
func NewOptionalSingleFromChannel(ch chan Optional) OptionalSingle {
	return &optionalSingle{
		itemChannel: ch,
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
		itemChannel: out,
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

func (s *single) Subscribe(handler EventHandler, opts ...Option) Observer {
	ob := CheckSingleEventHandler(handler)

	go func() {
		it := s.iterable.Iterator(context.Background())
		if item, err := it.Next(context.Background()); err == nil {
			switch item := item.(type) {
			case error:
				ob.OnError(item)
			default:
				ob.OnNext(item)
				ob.Dispose()
			}
		} else {
			ob.OnDone()
		}
	}()

	return ob
}

func (s *optionalSingle) Subscribe(handler EventHandler, opts ...Option) Observer {
	ob := CheckSingleEventHandler(handler)

	// TODO Improve
	go func() {
		for item := range s.itemChannel {
			switch item := item.(type) {
			case error:
				ob.OnError(item)
				return
			default:
				ob.OnNext(item)
				ob.Dispose()
			}
		}
	}()

	return ob
}
