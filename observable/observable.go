package observable

import (
	"sync"
	"time"

	"github.com/jochasinga/grx/bang"
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/emittable"
	"github.com/jochasinga/grx/errors"
	"github.com/jochasinga/grx/eventstream"
	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/subject"
	"github.com/jochasinga/grx/subscription"
)

// DirectiveFunc defines a func that should be passed to the observable.Start method,
// and represents a simple func that takes no arguments and return a bases.Emitter type.
type Directive func() bases.Emitter

// Observable is a stream of Emitters
type Observable struct {
	eventstream.EventStream
	subscriptor bases.Subscriptor
	notifier    *bang.Notifier
	observer    *subject.Subject
}

// DefaultObservable is a default Observable used by the constructor New.
// It is preferable to using the new keyword to create one.
var DefaultObservable = func() *Observable {
	o := &Observable{
		EventStream: eventstream.New(),
		notifier:    bang.New(),
		subscriptor: subscription.DefaultSubscription,
	}
	o.observer = subject.New(func(s *subject.Subject) {
		s.Stream = o
	})
	return o
}()

func (o *Observable) Done() {
	o.notifier.Done()
}

func (o *Observable) Unsubscribe() bases.Subscriptor {
	o.notifier.Unsubscribe()
	return o.subscriptor.Unsubscribe()
}

// New returns a new pointer to a default Observable.
func New(fs ...func(*Observable)) *Observable {
	o := DefaultObservable
	if len(fs) > 0 {
		for _, f := range fs {
			f(o)
		}
	}
	return o
}

// Create creates a new Observable provided by one or more function that takes an Observer as an argument
func Create(f func(*observer.Observer), fs ...func(*observer.Observer)) *Observable {
	o := DefaultObservable
	fs = append([]func(*observer.Observer){f}, fs...)
	go func() {
		for _, f := range fs {
			if observableRef, ok := o.observer.Sentinel.(*observer.Observer); ok {
				f(observableRef)
			}
		}
	}()
	return o
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

// Empty creates an Observable with one last item marked as "completed".
func Empty() *Observable {
	o := New()
	go func() {
		o.Done()
	}()
	return o
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(d time.Duration) *Observable {
	o := New()
	go func() {
		i := 0
		for {
			o.EventStream <- emittable.From(i)
			<-time.After(d)
			i++
		}
	}()
	return o
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, end int) *Observable {
	o := New()
	go func() {
		i := start
		for i < end {
			o.EventStream <- emittable.From(i)
			i++
		}
		o.Done()
	}()
	return o
}

// Just creates an observable with only one item and emit "as-is".
func Just(v interface{}, any ...interface{}) *Observable {
	any = append([]interface{}{v}, any...)
	o := New(func(o *Observable) {
		o.EventStream = make(eventstream.EventStream, len(any))
	})

	go func() {
		for _, val := range any {
			o.EventStream <- emittable.From(val)
		}
		close(o.EventStream)
		//o.Done()
	}()
	return o
}

// From creates an Observable from an Iterator type
func From(iter bases.Iterator) *Observable {
	o := New(func(o *Observable) {
		o.EventStream = eventstream.From(iter)
	})
	return o
}

// Start creates an Observable from one or more directive-like functions
func Start(f Directive, fs ...Directive) *Observable {
	fs = append([]Directive{f}, fs...)
	o := New(func(o *Observable) {
		o.EventStream = make(eventstream.EventStream, len(fs))
	})

	var wg sync.WaitGroup
	wg.Add(len(fs))
	for _, f := range fs {
		go func(f Directive) {
			o.EventStream <- f()
			wg.Done()
		}(f)
	}
	go func() {
		wg.Wait()
		//o.Done()
		close(o.EventStream)
	}()
	return o
}

func checkObservable(o *Observable) error {
	switch {
	case o == nil:
		return NewError(errors.NilObservableError)
	case o.EventStream == nil:
		return eventstream.NewError(errors.NilEventStreamError)
	default:
		break
	}
	return nil
}

// Subscribe subscribes an EventHandler to the receiving Observable and starts the stream
func (o *Observable) Subscribe(handler bases.EventHandler) (bases.Subscriptor, error) {
	if err := checkObservable(o); err != nil {
		return nil, err
	}

	isObserver := false

	var (
		ob    *observer.Observer
		nextf handlers.NextFunc
		errf  handlers.ErrFunc
		donef handlers.DoneFunc
	)

	switch handler := handler.(type) {
	case handlers.NextFunc:
		nextf = handler
	case handlers.ErrFunc:
		errf = handler
	case handlers.DoneFunc:
		donef = handler
	case *observer.Observer:
		ob = handler
		isObserver = true
	}

	if !isObserver {
		ob = observer.New(func(ob *observer.Observer) {
			ob.NextHandler = nextf
			ob.ErrHandler = errf
			ob.DoneHandler = donef
		})
	}

	// TODO: This should be asynchronous
	for emitter := range o.EventStream {
		ob.Handle(emitter)
	}

	ob.DoneHandler.Handle(emittable.DefaultEmittable)

	return o.subscriptor.Subscribe(), nil
}
