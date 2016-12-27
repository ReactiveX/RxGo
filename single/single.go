package single

import (
	"time"

	"github.com/jochasinga/grx/bang"
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/errors"
	"github.com/jochasinga/grx/eventstream"
	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/subscription"
)

// Single is a Stream that either emits an Item or an error once
type Single struct {
	eventstream.EventStream
	subscriptor bases.Subscriptor
	notifier    *bang.Notifier
}

var DefaultSingle = &Single{
	EventStream: eventstream.New(),
	notifier:    bang.New(),
	subscriptor: subscription.DefaultSubscription,
}

// Emit type-asserts its value and return an Item type or an error
func (s *Single) Emit() (bases.Item, error) {
	once := <-s.EventStream
	s.Done()
	item, err := once.Emit()
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Next always return an EndOfIteratorError after the first emission
func (s *Single) Next() (bases.Emitter, error) {
	emitter, err := s.EventStream.Next()
	if err != nil {
		return nil, NewError(errors.EndOfIteratorError)
	}
	return emitter, nil
}

func (s *Single) Done() {
	s.notifier.Done()
}

func (s *Single) Subscribe(handler bases.EventHandler) (bases.Subscriptor, error) {
	if s == nil {
		return nil, NewError(errors.NilSingleError)
	}

	ob := observer.DefaultObserver
	isObserver := false

	var (
		nextf handlers.NextFunc
		errf  handlers.ErrFunc
	)

	switch handler := handler.(type) {
	case handlers.NextFunc:
		nextf = handler
	case handlers.ErrFunc:
		errf = handler
	case *observer.Observer:
		ob = handler
		isObserver = true
	}

	if !isObserver {
		ob = observer.New(func(ob *observer.Observer) {
			ob.NextHandler = nextf
			ob.ErrHandler = errf
		})
	}

	go func() {
		for emitter := range s.EventStream {
			ob.Handle(emitter)
		}
	}()

	return bases.Subscriptor(&subscription.Subscription{
		SubscribeAt: time.Now(),
	}), nil
}

func (s *Single) Unsubscribe() bases.Subscriptor {
	s.notifier.Unsubscribe()
	return s.subscriptor.Unsubscribe()
}

// From creates a Single from an empty interface type
func New(e bases.Emitter) *Single {
	s := &Single{
		EventStream: make(eventstream.EventStream),
	}
	go func() {
		s.EventStream <- e
		close(s.EventStream)
	}()
	return s
}
