package rxgo

import (
	"context"
	"sync"
)

type createIterable struct {
	next                   <-chan Item
	opts                   []Option
	subscribers            []chan Item
	mutex                  sync.RWMutex
	producerAlreadyCreated bool
}

func newCreateIterable(fs []Producer, opts ...Option) Iterable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		for _, f := range fs {
			f(ctx, next)
		}
	}()

	return &createIterable{
		opts: opts,
		next: next,
	}
}

func (i *createIterable) Observe(opts ...Option) <-chan Item {
	mergedOptions := append(i.opts, opts...)
	option := parseOptions(mergedOptions...)

	if !option.isConnectable() {
		return i.next
	}

	if option.isConnectOperation() {
		i.connect(option.buildContext())
		return nil
	}

	ch := option.buildChannel()
	i.mutex.Lock()
	i.subscribers = append(i.subscribers, ch)
	i.mutex.Unlock()
	return ch
}

func (i *createIterable) connect(ctx context.Context) {
	i.mutex.Lock()
	if !i.producerAlreadyCreated {
		go i.produce(ctx)
		i.producerAlreadyCreated = true
	}
	i.mutex.Unlock()
}

func (i *createIterable) produce(ctx context.Context) {
	defer func() {
		i.mutex.RLock()
		for _, subscriber := range i.subscribers {
			close(subscriber)
		}
		i.mutex.RUnlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-i.next:
			if !ok {
				return
			}
			i.mutex.RLock()
			for _, subscriber := range i.subscribers {
				subscriber <- item
			}
			i.mutex.RUnlock()
		}
	}
}
