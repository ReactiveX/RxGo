package rxgo

import (
	"context"
	"sync"
	"sync/atomic"
)

// Amb takes several Observables, emit all of the items from only the first of these Observables
// to emit an item or notification.
func Amb(observables []Observable, opts ...Option) Observable {
	option := parseOptions(opts...)
	ctx := option.buildContext()
	next := option.buildChannel()
	once := sync.Once{}

	f := func(o Observable) {
		it := o.Observe()

		select {
		case <-ctx.Done():
			return
		case item, ok := <-it:
			if !ok {
				return
			}
			once.Do(func() {
				defer close(next)
				if item.IsError() {
					next <- item
					return
				}
				next <- item
				for {
					select {
					case <-ctx.Done():
						return
					case item, ok := <-it:
						if !ok {
							return
						}
						if item.IsError() {
							next <- item
							return
						}
						next <- item
					}
				}
			})
		}
	}

	for _, o := range observables {
		go f(o)
	}

	return &observable{
		iterable: newChannelIterable(next),
	}
}

// CombineLatest combines the latest item emitted by each Observable via a specified function
// and emit items based on the results of this function.
func CombineLatest(f FuncN, observables []Observable, opts ...Option) Observable {
	option := parseOptions(opts...)
	ctx := option.buildContext()
	next := option.buildChannel()

	go func() {
		size := uint32(len(observables))
		var counter uint32
		s := make([]interface{}, size)
		mutex := sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(int(size))
		errCh := make(chan struct{})

		handler := func(ctx context.Context, it Iterable, i int) {
			defer wg.Done()
			observe := it.Observe()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-observe:
					if !ok {
						return
					}
					if item.IsError() {
						next <- item
						errCh <- struct{}{}
						return
					}
					if s[i] == nil {
						atomic.AddUint32(&counter, 1)
					}
					mutex.Lock()
					s[i] = item.Value
					mutex.Unlock()
					if atomic.LoadUint32(&counter) == size {
						next <- FromValue(f(s...))
					}
				}
			}
		}

		ctx, cancel := context.WithCancel(ctx)
		for i, o := range observables {
			go handler(ctx, o, i)
		}

		go func() {
			for range errCh {
				cancel()
			}
		}()

		wg.Wait()
		close(next)
		close(errCh)
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Concat emits the emissions from two or more Observables without interleaving them.
func Concat(observables []Observable, opts ...Option) Observable {
	option := parseOptions(opts...)
	ctx := option.buildContext()
	next := option.buildChannel()

	go func() {
		defer close(next)
		for _, obs := range observables {
			observe := obs.Observe()
		loop:
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-observe:
					if !ok {
						break loop
					}
					if item.IsError() {
						next <- item
						return
					}
					next <- item
				}
			}
		}
	}()
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Empty creates an Observable with no item and terminate immediately.
func Empty() Observable {
	next := make(chan Item)
	close(next)
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// FromChannel creates a cold observable from a channel.
func FromChannel(next <-chan Item) Observable {
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// FromEventSource creates a hot observable from a channel.
func FromEventSource(ctx context.Context, next <-chan Item, strategy BackpressureStrategy) Observable {
	return &observable{
		iterable: newEventSourceIterable(ctx, next, strategy),
	}
}

// FromFuncs creates an observable from multiple functions.
func FromFuncs(f ...Scatter) Observable {
	return &observable{
		iterable: newFuncsIterable(f...),
	}
}

// Just creates an Observable with the provided items.
func Just(item Item, items ...Item) Observable {
	if len(items) > 0 {
		items = append([]Item{item}, items...)
	} else {
		items = []Item{item}
	}
	return &observable{
		iterable: newSliceIterable(items),
	}
}

// JustItem creates a single from one item.
func JustItem(item Item) Single {
	return &single{
		iterable: newSliceIterable([]Item{item}),
	}
}
