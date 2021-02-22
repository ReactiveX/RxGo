package rxgo

import (
	"context"
	"sync"
)

type eventSourceIterable struct {
	sync.RWMutex
	observers []chan Item
	disposed  bool
	opts      []Option
}

func newEventSourceIterable(ctx context.Context, next <-chan Item, strategy BackpressureStrategy, opts ...Option) Iterable {
	it := &eventSourceIterable{
		observers: make([]chan Item, 0),
		opts:      opts,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				it.closeAllObservers()
				return
			case item, ok := <-next:
				if !ok {
					it.closeAllObservers()
					return
				}
				it.RLock()
				switch strategy {
				default:
					fallthrough
				case Block:
					for _, observer := range it.observers {
						if !item.SendContext(ctx, observer) {
							return
						}
					}
				case Drop:
					for _, observer := range it.observers {
						select {
						default:
						case <-ctx.Done():
							return
						case observer <- item:
						}
					}
				}
				it.RUnlock()
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
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()

	i.Lock()
	if i.disposed {
		close(next)
	} else {
		i.observers = append(i.observers, next)
	}
	i.Unlock()
	return next
}
