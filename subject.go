package rxgo

import (
	"sync"
)

// ISubject defines subject API
type ISubject interface {
	Subscribe() (Subscription, Observable)
	Unsubscribe(id int)
	Next(value interface{})
	Error(err error)
	Complete()
}

// Subject a basic subject
type Subject struct {
	sync.RWMutex
	opts             []Option
	subscribers      map[int]chan<- Item
	nextSubscriberId int
}

// NewSubject creates a new subject.  with the specified observer options.
func NewSubject(opts ...Option) *Subject {
	res := Subject{
		opts:             opts,
		subscribers:      make(map[int]chan<- Item),
		nextSubscriberId: 0,
	}

	return &res
}

// Subscribe adds a subscriber to the subject. THe function returns a subscription and a new Observable.
func (s *Subject) Subscribe() (Subscription, Observable) {
	s.Lock()
	defer s.Unlock()

	return s.createSubscription(0)
}

func (s *Subject) createSubscription(bufferSize int) (Subscription, Observable) {
	id := s.nextSubscriberId
	s.nextSubscriberId++

	subChan := make(chan Item, bufferSize)
	s.subscribers[id] = subChan

	sub := NewSubscription(id, s)
	obs := FromEventSource(subChan, s.opts...)

	return sub, obs
}

// Unsubscribe removes a subscriber identified by ID from the Subject.
func (s *Subject) Unsubscribe(id int) {
	s.Lock()
	defer s.Unlock()

	subChan, found := s.subscribers[id]
	if found {
		close(subChan)
		delete(s.subscribers, id)
	}
}

// Next sends a new value to all subscribers
func (s *Subject) Next(value interface{}) {
	s.RLock()
	defer s.RUnlock()

	s.publish(Of(value))
}

// Error calls the error function on all subscribers
func (s *Subject) Error(err error) {
	s.RLock()
	defer s.RUnlock()

	s.publish(Error(err))
}

// publish sends an item to all subscribers
func (s *Subject) publish(item Item) {
	for _, subChan := range s.subscribers {
		subChan <- item
	}
}

// Complete closes all subscribers.
func (s *Subject) Complete() {
	s.Lock()
	defer s.Unlock()

	for _, subChan := range s.subscribers {
		close(subChan)
	}
}
