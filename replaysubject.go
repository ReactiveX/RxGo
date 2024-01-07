package rxgo

import (
	"container/list"
	"sync"
)

// ReplaySubject subject which replays the last received items to new subscribers
type ReplaySubject struct {
	Subject
	buffer         *list.List
	bufferLock     sync.Mutex
	maxReplayItems int
}

// NewReplaySubject creates a new replay subject
func NewReplaySubject(maxReplayItems int, opts ...Option) *ReplaySubject {
	res := ReplaySubject{
		Subject:        *NewSubject(opts...), // subscriber must be able to received current buffer and new items
		maxReplayItems: maxReplayItems,
		buffer:         list.New(),
		bufferLock:     sync.Mutex{},
	}

	return &res
}

// Next shadows base next function to capture the item history
func (s *ReplaySubject) Next(value interface{}) {
	s.bufferLock.Lock()
	defer s.bufferLock.Unlock()

	// add to buffer
	s.buffer.PushBack(value)
	// check for max length
	if s.buffer.Len() > s.maxReplayItems {
		// remove oldest item at the front
		s.buffer.Remove(s.buffer.Front())
	}

	s.Subject.Next(value)
}

// Subscribe shadows base subscribe function to replay the item history
func (s *ReplaySubject) Subscribe() (Subscription, Observable) {
	s.Lock()
	defer s.Unlock()

	// create buffered channel to hold all current replay items
	sub, obs := s.createSubscription(s.buffer.Len())
	subChan := s.subscribers[sub.GetId()]

	s.bufferLock.Lock()
	defer s.bufferLock.Unlock()

	// replay buffered items
	elem := s.buffer.Front()
	for elem != nil {
		subChan <- Of(elem.Value)
		elem = elem.Next()
	}

	return sub, obs
}
