package observable

import (
	"sync"
	"time"

	"github.com/jochasinga/rx"
	"github.com/jochasinga/rx/errors"
	"github.com/jochasinga/rx/fx"
	"github.com/jochasinga/rx/handlers"
	"github.com/jochasinga/rx/observer"
	"github.com/jochasinga/rx/subscription"
)

// Observable is a basic observable channel
type Observable <-chan interface{}

var DefaultObservable = make(Observable)

// New creates an Observable
func New(buffer uint) Observable {
	return make(Observable, int(buffer))
}

// CheckHandler checks the underlying type of an EventHandler.
func CheckEventHandler(handler rx.EventHandler) observer.Observer {
	ob := observer.DefaultObserver
	switch handler := handler.(type) {
	case handlers.NextFunc:
		ob.NextHandler = handler
	case handlers.ErrFunc:
		ob.ErrHandler = handler
	case handlers.DoneFunc:
		ob.DoneHandler = handler
	case observer.Observer:
		ob = handler
	}
	return ob
}

// Next returns the next item on the Observable.
func (o Observable) Next() (interface{}, error) {
	if next, ok := <-o; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// Subscribe subscribes an EventHandler and returns a Subscription channel.
func (o Observable) Subscribe(handler rx.EventHandler) <-chan subscription.Subscription {
	done := make(chan subscription.Subscription)
	sub := subscription.New().Subscribe()

	ob := CheckEventHandler(handler)

	go func() {
	OuterLoop:
		for item := range o {
			switch item := item.(type) {
			case error:
				ob.OnError(item)

				// Record the error and break the loop.
				sub.Error = item
				break OuterLoop
			default:
				ob.OnNext(item)
			}
		}

		// OnDone only gets executed if there's no error.
		if sub.Error == nil {
			ob.OnDone()
		}

		done <- sub.Unsubscribe()
		return
	}()

	return done
}

/*
func (o Observable) Unsubscribe() subscription.Subscription {
	// Stub: to be implemented
	return subscription.New()
}
*/

// Map maps a MappableFunc predicate to each item in Observable and
// returns a new Observable with applied items.
func (o Observable) Map(apply fx.MappableFunc) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o {
			out <- apply(item)
		}
		close(out)
	}()
	return Observable(out)
}

// Filter filters items in the original Observable and returns
// a new Observable with the filtered items.
func (o Observable) Filter(apply fx.FilterableFunc) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o {
			if apply(item) {
				out <- item
			}
		}
		close(out)
	}()
	return Observable(out)
}

//Distinct supress duplicate items in the original Observable and returns
// a new Observable.
func (o Observable) Distinct(apply fx.KeySelectorFunc) Observable {
	out := make(chan interface{})
	go func() {
		keysets := make(map[interface{}]struct{})
		for item := range o {
			key := apply(item)
			_, ok := keysets[key]
			if !ok {
				out <- item
			}
			keysets[key] = struct{}{}
		}
		close(out)
	}()
	return Observable(out)
}

func (o Observable) DistinctUntilChanged(apply fx.KeySelectorFunc) Observable {
	out := make(chan interface{})
	go func() {
		var current interface{}
		for item := range o {
			key := apply(item)
			if current != key {
				out <- item
				current = key
			}
		}
		close(out)
	}()
	return Observable(out)
}

// Scan applies ScannableFunc predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o Observable) Scan(apply fx.ScannableFunc) Observable {
	out := make(chan interface{})

	go func() {
		var current interface{}
		for item := range o {
			out <- apply(current, item)
			current = apply(current, item)
		}
		close(out)
	}()
	return Observable(out)
}

// From creates a new Observable from an Iterator.
func From(it rx.Iterator) Observable {
	source := make(chan interface{})
	go func() {
		for {
			val, err := it.Next()
			if err != nil {
				break
			}
			source <- val
		}
		close(source)
	}()
	return Observable(source)
}

// Empty creates an Observable with no item and terminate immediately.
func Empty() Observable {
	source := make(chan interface{})
	go func() {
		close(source)
	}()
	return Observable(source)
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(term chan struct{}, timeout time.Duration) Observable {
	source := make(chan interface{})
	go func(term chan struct{}) {
		i := 0
	OuterLoop:
		for {
			select {
			case <-term:
				break OuterLoop
			case <-time.After(timeout):
				source <- i
			}
			i++
		}
		close(source)
	}(term)
	return Observable(source)
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, end int) Observable {
	source := make(chan interface{})
	go func() {
		i := start
		for i < end {
			source <- i
			i++
		}
		close(source)
	}()
	return Observable(source)
}

// Just creates an Observable with the provided item(s).
func Just(item interface{}, items ...interface{}) Observable {
	source := make(chan interface{})
	if len(items) > 0 {
		items = append([]interface{}{item}, items...)
	} else {
		items = []interface{}{item}
	}

	go func() {
		for _, item := range items {
			source <- item
		}
		close(source)
	}()

	return Observable(source)
}

// Start creates an Observable from one or more directive-like EmittableFunc
// and emits the result of each operation asynchronously on a new Observable.
func Start(f fx.EmittableFunc, fs ...fx.EmittableFunc) Observable {
	if len(fs) > 0 {
		fs = append([]fx.EmittableFunc{f}, fs...)
	} else {
		fs = []fx.EmittableFunc{f}
	}

	source := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f fx.EmittableFunc) {
			source <- f()
			wg.Done()
		}(f)
	}

	// Wait in another goroutine to not block
	go func() {
		wg.Wait()
		close(source)
	}()

	return Observable(source)
}
