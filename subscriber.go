package rxgo

import (
	"sync"
)

type safeSubscriber[T any] struct {
	// prevent concurrent race on unsubscribe
	mu sync.RWMutex

	// signal to indicate the subscribe has ended
	// dispose <-chan struct{}

	ch chan DataValuer[T]

	// determine the channel was closed
	closed bool

	dst Observer[T]
}

func NewSafeSubscriber[T any](onNext func(T), onError func(error), onComplete func()) *safeSubscriber[T] {
	sub := &safeSubscriber[T]{
		ch: make(chan DataValuer[T]),
		dst: &consumerObserver[T]{
			onNext:     onNext,
			onError:    onError,
			onComplete: onComplete,
		},
	}
	return sub
}

func (s *safeSubscriber[T]) Closed() bool {
	return s.closed
}

func (s *safeSubscriber[T]) ForEach() <-chan DataValuer[T] {
	return s.ch
}

func (s *safeSubscriber[T]) Next(v T) {
	if s.closed {
		return
	}
	emitData(v, s.ch)
}

func (s *safeSubscriber[T]) Error(err error) {
	if s.closed {
		return
	}
	emitError(err, s.ch)
	s.closeChannel()
}

func (s *safeSubscriber[T]) Complete() {
	s.closeChannel()
}

// this will close the stream and stop the emission of the stream data
func (s *safeSubscriber[T]) Unsubscribe() {
	s.closeChannel()
}

func (s *safeSubscriber[T]) closeChannel() {
	s.mu.Lock()
	if !s.closed {
		s.closed = true
		close(s.ch)
	}
	s.mu.Unlock()
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
