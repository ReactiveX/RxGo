package rxgo

import (
	"sync"
)

type subscriber[T any] struct {
	// to prevent DATA RACE
	mu *sync.RWMutex

	// channel to transfer data
	ch chan Notification[T]

	// channel to indentify it has stopped
	stop chan struct{}

	isStopped bool

	// determine the channel was closed
	closed bool
}

func NewSubscriber[T any](bufferCount ...uint) *subscriber[T] {
	ch := make(chan Notification[T])
	if len(bufferCount) > 0 {
		ch = make(chan Notification[T], bufferCount[0])
	}
	return &subscriber[T]{
		mu:   new(sync.RWMutex),
		ch:   ch,
		stop: make(chan struct{}),
	}
}

func (s *subscriber[T]) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isStopped {
		return
	}
	s.isStopped = true
	close(s.stop)
}

func (s *subscriber[T]) Closed() <-chan struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stop
}

func (s *subscriber[T]) ForEach() <-chan Notification[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ch
}

func (s *subscriber[T]) Send() chan<- Notification[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ch
}

// this will close the stream and stop the emission of the stream data
func (s *subscriber[T]) Unsubscribe() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closeChannel()
}

func (s *subscriber[T]) closeChannel() {
	s.closed = true
	close(s.ch)
}

type safeSubscriber[T any] struct {
	*subscriber[T]

	dst Observer[T]
}

func NewSafeSubscriber[T any](onNext OnNextFunc[T], onError OnErrorFunc, onComplete OnCompleteFunc) *safeSubscriber[T] {
	sub := &safeSubscriber[T]{
		subscriber: NewSubscriber[T](),
		dst:        NewObserver(onNext, onError, onComplete),
	}
	return sub
}

func NewObserver[T any](onNext OnNextFunc[T], onError OnErrorFunc, onComplete OnCompleteFunc) Observer[T] {
	if onNext == nil {
		onNext = func(T) {}
	}
	if onError == nil {
		onError = func(error) {}
	}
	if onComplete == nil {
		onComplete = func() {}
	}
	return &consumerObserver[T]{
		onNext:     onNext,
		onError:    onError,
		onComplete: onComplete,
	}
}

type consumerObserver[T any] struct {
	onNext     func(T)
	onError    func(error)
	onComplete func()
}

func (o *consumerObserver[T]) Next(v T) {
	o.onNext(v)
}

func (o *consumerObserver[T]) Error(err error) {
	o.onError(err)
}

func (o *consumerObserver[T]) Complete() {
	o.onComplete()
}
