package rxgo

import (
	"sync"
)

// BehaviorSubject subject which returns the last received item to new subscribers
type BehaviorSubject struct {
	Subject
	lastValue     interface{}
	lastValueLock sync.Mutex
}

// NewBehaviorSubject Creates a new behavior subject
func NewBehaviorSubject(opts ...Option) *BehaviorSubject {
	res := BehaviorSubject{
		Subject:       *NewSubject(opts...), // subscriber must be able to receive last item and new items
		lastValueLock: sync.Mutex{},
	}

	return &res
}

// Next shadows base next function to capture the last item.
func (s *BehaviorSubject) Next(value interface{}) {
	s.lastValueLock.Lock()
	defer s.lastValueLock.Unlock()

	s.lastValue = value

	s.Subject.Next(value)
}

// Subscribe shadows base subscribe function to replay the last captured item.
func (s *BehaviorSubject) Subscribe() (Subscription, Observable) {
	s.Lock()
	defer s.Unlock()

	// create buffered channel to hold last item
	sub, obs := s.createSubscription(1)
	subChan := s.subscribers[sub.GetId()]

	if s.lastValue != nil {
		s.lastValueLock.Lock()
		subChan <- Of(s.lastValue)
		s.lastValueLock.Unlock()
	}

	return sub, obs
}
