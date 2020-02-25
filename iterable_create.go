package rxgo

import (
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

	wg := sync.WaitGroup{}
	for _, f := range fs {
		f := f
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(ctx, next)
		}()
	}
	go func() {
		wg.Wait()
		close(next)
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
		i.connect()
		return nil
	}

	ch := option.buildChannel()
	i.mutex.Lock()
	i.subscribers = append(i.subscribers, ch)
	i.mutex.Unlock()
	return ch
}

func (i *createIterable) connect() {
	i.mutex.Lock()
	if !i.producerAlreadyCreated {
		go i.produce()
		i.producerAlreadyCreated = true
	}
	i.mutex.Unlock()
}

func (i *createIterable) produce() {
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
