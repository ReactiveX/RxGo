package rxgo

import (
	"context"
	"sync"
)

type eventSourceIterable struct {
	sync.RWMutex
	observers []chan Item
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
	i.Unlock()
}

func (i *eventSourceIterable) Observe() <-chan Item {
	next := make(chan Item)
	i.Lock()
	i.observers = append(i.observers, next)
	i.Unlock()
	return next
}
