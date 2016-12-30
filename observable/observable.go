package observable

import (
	"sync"
	//"time"

	//"github.com/jochasinga/grx/bang"
	"github.com/jochasinga/grx/bases"
	//"github.com/jochasinga/grx/emittable"
	//"github.com/jochasinga/grx/errors"
	//"github.com/jochasinga/grx/eventstream"
	//"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observer"
	//"github.com/jochasinga/grx/subject"
	//"github.com/jochasinga/grx/subscription"
)

// DirectiveFunc defines a func that should be passed to the observable.Start method,
// and represents a simple func that takes no arguments and return a bases.Emitter type.
type Directive func() bases.Emitter

type Applicator interface {
	Apply(bases.Emitter)
}

type MappableFunc func(bases.Emitter) bases.Emitter
type CurryableFunc func(interface{}) MappableFunc

//func (fx Mappable) Apply(e bases.Emitter) {}
//type Curryable func(bases.Emitter) func(interface{}) bases.Emittable
//func (fx Curryable) Apply(e bases.Emitter) {}

// Observable is a stream of Emitters
//type Observable struct {
//	source chan bases.Emitter
//subscriptor bases.Subscriptor
//notifier    *bang.Notifier
//observer    *subject.Subject
//}

type Basic <-chan bases.Emitter

func (bs Basic) Subscribe(ob observer.Observer) <-chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		for e := range bs {
			ob.NextHandler.Handle(e)
		}
		ob.DoneHandler()
		done <- struct{}{}
	}()
	return done
}

func (bs Basic) Map(fx MappableFunc) Basic {
	out := make(chan bases.Emitter)
	go func() {
		for e := range bs {
			out <- fx(e)
		}
		close(out)
	}()
	return Basic(out)
}

type Connectable struct {
	emitters  <-chan bases.Emitter
	observers []observer.Observer
}

func (cnxt Connectable) Subscribe(ob observer.Observer) Connectable {
	cnxt.observers = append(cnxt.observers, ob)
	return cnxt
}

func (cnxt Connectable) Map(fx MappableFunc) Connectable {
	out := make(chan bases.Emitter)
	go func() {
		for e := range cnxt.emitters {
			out <- fx(e)
		}
		close(out)
	}()
	return Connectable{emitters: out}
}

func (cnxt Connectable) Connect() <-chan struct{} {
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup

	for _, ob := range cnxt.observers {
		wg.Add(1)
		go func(ob observer.Observer) {
			for e := range cnxt.emitters {
				ob.NextHandler.Handle(e)
			}
			ob.DoneHandler()
			wg.Done()
		}(ob)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}

func BasicFrom(es []bases.Emitter) Basic {
	source := make(chan bases.Emitter, len(es))
	go func() {
		for _, e := range es {
			source <- e
		}
		close(source)
	}()
	return Basic(source)
}

func ConnectableFrom(es []bases.Emitter) Connectable {
	source := make(chan bases.Emitter)
	go func() {
		for _, e := range es {
			source <- e
		}
		close(source)
	}()
	return Connectable{emitters: source}
}

// DefaultObservable is a default Observable used by the constructor New.
// It is preferable to using the new keyword to create one.
/*
var DefaultObservable = func() *Observable {
	o := &Observable{
		EventStream: make(eventstream.EventStream),
		notifier:    bang.New(),
		subscriptor: subscription.DefaultSubscription,
	}
	o.observer = subject.New(func(s *subject.Subject) {
		s.Stream = o
		s.Sentinel = o.observer
	})
	return o
}()
*/
//var DefaultObservable = Observable{}

/*
func (o *Observable) Done() {
	o.notifier.Done()
}
*/

/*
func (o *Observable) Unsubscribe() bases.Subscriptor {
	o.notifier.Unsubscribe()
	return o.subscriptor.Unsubscribe()
}
*/

/*
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
*/

/*
// Create creates a new Observable provided by one or more function that takes an Observer as an argument
func Create(f func(*observer.Observer)) *Observable {
	//fs = append([]func(*observer.Observer){f}, fs...)
	o := DefaultObservable
	f(o)


		go func() {
			for _, f := range fs {
				f(o.observer.Sentinel.(*observer.Observer))
			}
			//o.observer.Sentinel = ob
			//close(o.EventStream)
		}()

	return o
}
*/

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

/*
// Empty creates an Observable with one last item marked as "completed".
func Empty() *Observable {
	o := New()
	go func() {
		//o.Done()
		close(o.EventStream)
	}()
	return o
}
*/

/*
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
*/

/*
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
*/

/*
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
*/

/*
// From creates an Observable from an Iterator type
func From(iter bases.Iterator) *Observable {
	o := New(func(o *Observable) {
		o.EventStream = eventstream.From(iter)
	})
	return o
}
*/

/*
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
*/

/*
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
*/

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
