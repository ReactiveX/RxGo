package observable

import (
	"sync"
	"time"

	//"github.com/jochasinga/grx/bang"
	"github.com/jochasinga/grx/bases"
	//"github.com/jochasinga/grx/errors"
	//"github.com/jochasinga/grx/emittable"
	"github.com/jochasinga/grx/fx"
	"github.com/jochasinga/grx/observer"
	//"github.com/jochasinga/grx/eventstream"
	"github.com/jochasinga/grx/handlers"
	//"github.com/jochasinga/grx/subject"
	"github.com/jochasinga/grx/subscription"
)

// Observable is a basic observable channel
type Observable <-chan interface{}

var DefaultObservable = make(Observable)

// New creates an Observable
func New(buffer uint) Observable {
	return make(Observable, int(buffer))
}

func checkEventHandler(handler bases.EventHandler) observer.Observer {
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
	done := make(chan subscription.Subscription, 1)
	sub := subscription.New().Subscribe()

	ob := checkEventHandler(handler)

	go func() {
	OuterLoop:
		for item := range o {
			switch item := item.(type) {
			case error:
				err := item.(error)
				ob.ErrHandler(err)

				// Record the error and return without completing.
				sub.Error = err
				break OuterLoop
			default:
				ob.OnNext(item)
			}
		}

		// This part only gets executed if there wasn't an error.
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

func From(items []interface{}) Observable {
	source := make(chan interface{}, len(items))
	go func() {
		for _, item := range items {
			source <- item
		}
		close(source)
	}()
	return Observable(source)
}

/*
// Add adds an item to the Observable and returns that Observable.
// If the Observable is done, it creates a new one and return it.
func (o *Observable) Add(e Emitter, es ...Emitter) *Observable {
	es = append([]Emitter{e}, es...)
	out := New()
	if !o.isDone() {
		go func() {
			for _, e := range es {
				o.EventStream <- e
			}
		}()
		*o = *out
		return out
	}
	go func() {
		for _, e := range es {
			out.EventStream <- e
		}
	}()
	return out
}
*/

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
func Interval(term chan struct{}, d time.Duration) Observable {
	source := make(chan interface{})
	go func() {
		i := 0
		for {
			select {
			case <-term:
				return
			case <-time.After(d):
				source <- i
			}
			i++
		}
		close(source)
	}()
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

	source := make(chan interface{}, len(fs))

	var wg sync.WaitGroup
	wg.Add(len(fs))
	for _, f := range fs {
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

/*
// Subscribe subscribes an EventHandler to the receiving Observable and starts the stream
func (o *Observable) Subscribe(handler bases.EventHandler) (bases.Subscriptor, error) {
	if err := checkObservable(o); err != nil {
		return nil, err
	}

	isObserver := false

	// Set up default handlers
	var (
		//ob    *observer.Observer
		nextf handlers.NextFunc = func(it bases.Item) {}
		errf  handlers.ErrFunc  = func(err error) {}
		donef handlers.DoneFunc = func() {}
	)

	switch handler := handler.(type) {
	case handlers.NextFunc:
		nextf = handler
	case handlers.ErrFunc:
		errf = handler
	case handlers.DoneFunc:
		donef = handler
	case *observer.Observer:
		//ob = handler
		o.observer.Sentinel = handler
		isObserver = true
	}

	if !isObserver {
		o.observer.Sentinel = observer.New(func(ob *observer.Observer) {
			ob.NextHandler = nextf
			ob.ErrHandler = errf
			ob.DoneHandler = donef
		})
	}

	// TODO: This should be asynchronous
	for emitter := range o.EventStream {
		if ob, ok := o.observer.Sentinel.(*observer.Observer); ok {
			ob.Handle(emitter)
		}
	}

	if ob, ok := o.observer.Sentinel.(*observer.Observer); ok {
		ob.DoneHandler.Handle(emittable.DefaultEmittable)
	}

	return o.subscriptor.Subscribe(), nil
}
*/
