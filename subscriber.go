package rxgo

import (
	"sync"
)

type safeSubscriber[T any] struct {
	// prevent concurrent race on unsubscribe
	mu sync.Mutex

	// signal to indicate the subscribe has ended
	dispose <-chan struct{}

	ch chan DataValuer[T]

	// determine the channel was closed
	closed bool

	dst Observer[T]
}

func NewSafeSubscriber[T any](dispose <-chan struct{}, onNext func(T), onError func(error), onComplete func()) *safeSubscriber[T] {
	sub := &safeSubscriber[T]{
		dispose: dispose,
		ch:      make(chan DataValuer[T]),
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

func (s *safeSubscriber[T]) Done() <-chan struct{} {
	return s.dispose
}

func (s *safeSubscriber[T]) Next(v T) {
	if s.closed {
		return
	}
	s.ch <- Data[T]{v: v}
}

func (s *safeSubscriber[T]) Error(err error) {
	if s.closed {
		return
	}
	s.ch <- Data[T]{err: err}
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
