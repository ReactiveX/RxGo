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
		defer func() {
			it.closeAllObservers()
		}()

		deliver := func(item Item) (done bool) {
			it.RLock()
			defer it.RUnlock()

			switch strategy {
			default:
				fallthrough
			case Block:
				for _, observer := range it.observers {
					if !item.SendContext(ctx, observer) {
						return true
					}
				}
			case Drop:
				for _, observer := range it.observers {
					select {
					default:
					case <-ctx.Done():
						return true
					case observer <- item:
					}
				}
			}
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-next:
				if !ok {
					return
				}

				if done := deliver(item); done {
					return
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
