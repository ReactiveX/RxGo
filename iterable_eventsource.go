package rxgo

import (
	"context"
	"sync"
)

type eventSourceIterable struct {
	sync.RWMutex
	observers []chan Item
	disposed  bool
}

func newEventSourceIterable(ctx context.Context, next <-chan Item, strategy BackpressureStrategy) Iterable {
	it := &eventSourceIterable{
		observers: make([]chan Item, 0),
	}

	go func() {
		switch strategy {
		default:
			fallthrough
		case Block:
			for {
				select {
				case <-ctx.Done():
					it.closeAllObservers()
					return
				case item := <-next:
					it.RLock()
					for _, observer := range it.observers {
						observer <- item
					}
					it.RUnlock()
				}
			}
		case Drop:
			for {
				select {
				case <-ctx.Done():
					it.closeAllObservers()
					return
				case item := <-next:
					it.RLock()
					for _, observer := range it.observers {
						select {
						default:
						case observer <- item:
						}
					}
					it.RUnlock()
				}
			}
		}
	}()

	return it
}

func (i *eventSourceIterable) closeAllObservers() {
	i.Lock()
	for _, observer := range i.observers {
		close(observer)
	}
	i.disposed = true
	i.Unlock()
}

func (i *eventSourceIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(opts...)
	var next chan Item
	if toBeBuffered, cap := option.withBuffer(); toBeBuffered {
		next = make(chan Item, cap)
	} else {
		next = make(chan Item)
	}

	i.Lock()
	if i.disposed {
		close(next)
	} else {
		i.observers = append(i.observers, next)
	}
	i.Unlock()
	return next
}
