package rxgo

import (
	"context"
	"sync"
)

type subscription struct {
	ctx context.Context
	ch  chan Item
}

type channelIterable struct {
	next                   <-chan Item
	opts                   []Option
	nextSubscriberID       uint64
	subscribers            map[uint64]*subscription
	mutex                  sync.RWMutex
	producerAlreadyCreated bool
}

func newChannelIterable(next <-chan Item, opts ...Option) Iterable {
	return &channelIterable{
		next:        next,
		subscribers: make(map[uint64]*subscription),
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
		i.connect(option.buildContext(emptyContext))
		return nil
	}

	ch := i.createSubscription(option)
	return ch
}

func (i *channelIterable) createSubscription(option Option) chan Item {
	ch := option.buildChannel()
	sctx := option.buildContext(emptyContext)

	i.mutex.Lock()
	sid := i.nextSubscriberID
	i.nextSubscriberID++
	i.subscribers[sid] = &subscription{
		ctx: sctx,
		ch:  ch,
	}
	i.mutex.Unlock()
	return ch
}

func (i *channelIterable) connect(ctx context.Context) {
	i.mutex.Lock()
	if !i.producerAlreadyCreated {
		go i.produce(ctx)
		i.producerAlreadyCreated = true
	}
	i.mutex.Unlock()
}

func (i *channelIterable) produce(ctx context.Context) {
	defer func() {
		i.mutex.RLock()
		for _, subscriber := range i.subscribers {
			close(subscriber.ch)
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
			toBeCleaned := make([]uint64, 0)
			i.mutex.RLock()
			for sid, subscriber := range i.subscribers {
				select {
				case <-subscriber.ctx.Done():
					toBeCleaned = append(toBeCleaned, sid)
				case subscriber.ch <- item:
				}
			}
			i.mutex.RUnlock()

			i.removeSubscriptions(toBeCleaned)
		}
	}
}

func (i *channelIterable) removeSubscriptions(sids []uint64) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for _, sid := range sids {
		delete(i.subscribers, sid)
	}
}
