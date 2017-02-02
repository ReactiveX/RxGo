package observable

import (
	"sync"
	"time"

	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/fx"
	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/subscription"
)

// Observable is a basic observable channel
type Observable <-chan interface{}

var DefaultObservable = make(Observable)

// New creates an Observable
func New(buffer uint) Observable {
	return make(Observable, int(buffer))
}

func CheckEventHandler(handler bases.EventHandler) observer.Observer {
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

// Make sure Observable implements base.Observable
//
// Subscribe returns a channel of empty struct
func (o Observable) Subscribe(handler bases.EventHandler) <-chan subscription.Subscription {
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

func (o Observable) Unsubscribe() subscription.Subscription {
	// Stub: to be implemented
	return subscription.New()
}

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

func From(it bases.Iterator) Observable {
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

// Empty creates an Observable with no item and terminate once subscribed to.
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

// Just creates an observable with the provided "as-is" item(s)
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

// Start creates an Observable from one or more directive-like functions
// and emit the result of each asynchronously.
func Start(f fx.DirectiveFunc, fs ...fx.DirectiveFunc) Observable {
	if len(fs) > 0 {
		fs = append([]fx.DirectiveFunc{f}, fs...)
	} else {
		fs = []fx.DirectiveFunc{f}
	}

	source := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f fx.DirectiveFunc) {
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
