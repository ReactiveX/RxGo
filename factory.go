package rxgo

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
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
					if atomic.LoadUint32(&counter) == size {
						next <- FromValue(f(s...))
					}
					mutex.Unlock()
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
func FromEventSource(next <-chan Item, opts ...Option) Observable {
	option := parseOptions(opts...)

	return &observable{
		iterable: newEventSourceIterable(option.buildContext(), next, option.buildBackPressureStrategy()),
	}
}

// FromSlice creates an observable from a slice.
func FromSlice(s []Item) Single {
	return &single{
		iterable: newSliceIterable(s),
	}
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(interval Duration, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		i := 0
		for {
			select {
			case <-time.After(interval.duration()):
				next <- FromValue(i)
				i++
			case <-ctx.Done():
				close(next)
				return
			}
		}
	}()
	return &observable{
		iterable: newEventSourceIterable(ctx, next, option.buildBackPressureStrategy()),
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

// Merge combines multiple Observables into one by merging their emissions
func Merge(observables []Observable, opts ...Option) Observable {
	option := parseOptions(opts...)
	ctx := option.buildContext()
	next := option.buildChannel()
	wg := sync.WaitGroup{}

	f := func(o Observable) {
		defer wg.Done()
		observe := o.Observe()
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
					return
				}
				next <- item
			}
		}
	}

	for _, o := range observables {
		wg.Add(1)
		go f(o)
	}

	go func() {
		wg.Wait()
		close(next)
	}()
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Never creates an Observable that emits no items and does not terminate.
func Never() Observable {
	next := make(chan Item)
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, count int) Observable {
	if count < 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "count must be positive"))
	}
	if start+count-1 > math.MaxInt32 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "max value is bigger than math.MaxInt32"))
	}
	return &observable{
		iterable: newRangeIterable(start, count),
	}
}

// Scatter creates an observable from multiple functions.
func Scatter(f ...ScatterFunc) Observable {
	return &observable{
		iterable: newFuncsIterable(f...),
	}
}

// Start creates an Observable from one or more directive-like Supplier
// and emits the result of each operation asynchronously on a new Observable.
func Start(fs []Supplier, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	var wg sync.WaitGroup
	for _, f := range fs {
		f := f
		wg.Add(1)
		go func() {
			next <- f(ctx)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(next)
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Timer returns an Observable that emits an empty structure after a specified delay, and then completes.
func Timer(d Duration, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		select {
		case <-ctx.Done():
			return
		case <-time.After(d.duration()):
			next <- FromValue(struct{}{})
		}
	}()
	return &observable{
		iterable: newChannelIterable(next),
	}
}
