package rxgo

import (
	"sync"
)

type channelIterable struct {
	next                   <-chan Item
	opts                   []Option
	subscribers            []chan Item
	mutex                  sync.RWMutex
	producerAlreadyCreated bool
}

func newChannelIterable(next <-chan Item, opts ...Option) Iterable {
	return &channelIterable{
		next:        next,
		subscribers: make([]chan Item, 0),
		opts:        opts,
	}
}

func (i *channelIterable) Observe(opts ...Option) <-chan Item {
	mergedOptions := append(i.opts, opts...)
	option := parseOptions(mergedOptions...)

	if !option.isConnectable() {
		return i.next
	}

	if option.isConnectOperation() {
		i.connect()
		return nil
	}

	ch := option.buildChannel()
	i.mutex.Lock()
	i.subscribers = append(i.subscribers, ch)
	i.mutex.Unlock()
	return ch
}

func (i *channelIterable) connect() {
	i.mutex.Lock()
	if !i.producerAlreadyCreated {
		go i.produce()
		i.producerAlreadyCreated = true
	}
	i.mutex.Unlock()
}

func (i *channelIterable) produce() {
	for item := range i.next {
		i.mutex.RLock()
		for _, subscriber := range i.subscribers {
			subscriber <- item
		}
		i.mutex.RUnlock()
	}

	i.mutex.RLock()
	for _, subscriber := range i.subscribers {
		close(subscriber)
	}
	i.mutex.RUnlock()
}
