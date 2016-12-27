package observable

import (
	"sync"
	"time"

	"github.com/jochasinga/grx/bang"
	"github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/eventstream"
	"github.com/jochasinga/grx/subject"
)

// Observable is a stream of Emitters
type Observable struct {
	eventstream.EventStream
	notifier *bang.Notifier
	observer *subject.Subject
}

// DefaultObservable is a default Observable used by the constructor New.
// It is preferable to using the new keyword to create one.
var DefaultObservable = func() *Observable {
	o := &Observable{
		EventStream: eventstream.New(),
		notifier:    bang.New(),
	}
	o.observer = subject.New(func(s *subject.Subject) {
		s.Stream = o
	})
	return o
}()

func (o *Observable) Done() {
	o.notifier.Done()
}

func (o *Observable) Unsubscribe() {
	o.notifier.Unsubscribe()
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
func Create(f func(*Observer), fs ...func(*Observer)) *Observable {
	o := New()
	fs = append([]func(*Observer){f}, fs...)
	go func() {
		for _, f := range fs {
			f(o)
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

// Interval creates an Observable emitting incremental integers infinitely
// between each given time interval.
func Interval(d time.Duration) *Observable {
	o := New()
	go func() {
		i := 0
		for {
			o.EventStream <- i
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
		for i := start; i < end; i++ {
			o.EventStream <- i
		}
		o.Done()
	}()
	return o
}

// Just creates an observable with only one item and emit "as-is".
// source := observable.Just("https://someurl.com/api")
func Just(em Emitter, ems ...Emitter) Observable {
	o := New()
	emitters := append([]Emitter{em}, ems...)

	go func() {
		for _, emitter := range emitters {
			o.EventStream <- emitter
		}
		o.Done()
	}()
	return o
}

/*
// From creates a new EventStream from an Iterator
func From(iter bases.Iterator) EventStream {
	es := make(EventStream)
	go func() {
		for {
			emitter, err := iter.Next()
			if err != nil {
				return
			}
			es <- emitter
		}
		close(es)
	}()
	return es
}
*/

// From creates an Observable from an Iterator
func From(iter bases.Iterator) *Observable {
	o := New()
	o.EventStream = o.From(iter)
	o.Done()
	return o
}

// Start creates an Observable from one or more directive-like functions
func Start(f func() Emitter, fs ...func() Emitter) *Observable {
	o := New()
	fs = append([](func() Emitter){f}, fs...)

	var wg sync.WaitGroup
	wg.Add(len(fs))
	for _, f := range fs {
		go func(f func() Emitter) {
			o.EventStream <- f()
			wg.Done()
		}(f)
	}
	go func() {
		wg.Wait()
		o.Done()
	}()
	return o
}

func processStream(o *Observable, ob *Observer) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for emitter := range o.EventStream {
			ob.Handle(emitter)
		}
		wg.Done()
	}()
	
	go func() {
		wg.Wait()
		o.Done()
	}()

	go func() {
		select {
		case <-o.unsubscribed:
			return
		case <-o.done:
			ob.DoneHandler()
			return
		}
	}()
}

func checkObservable(o *Observable) error {

	switch {
	case o == nil:
		return NewError(grx.NilObservableError)
	case o.EventStream == nil:
		return eventstream.NewError(grx.NilEventStreamError)
	case o.isDone():
		return NewError(grx.EndOfIteratorError)
	default:
		return nil
	}
	return nil
}
	
/*
// SubscribeWith subscribes handlers to the Observable and starts it.
//func (o *Observable) SubscribeFunc(nxtf func(v interface{}), errf func(e error), donef func()) (Subscriptor, error) {
func (o *BaseObservable) SubscribeFunc(nxtf func(v interface{}), errf func(e error), donef func()) (Subscriptor, error) {

	err := checkObservable(o)
	if err != nil {
		return nil, err
	}

	ob := Observer(&BaseObserver{
		NextHandler: NextFunc(nxtf),
		ErrHandler:  ErrFunc(errf),
		DoneHandler: DoneFunc(donef),
	})

	o.runStream(ob)

	return Subscriptor(&Subscription{SubscribeAt: time.Now()}), nil
}
*/

// SubscribeHandler subscribes a Handler to the Observable and starts it.
func (o *Observable) Subscribe(handle EventHandler) (Subscriptor, error) {
	err := checkObservable(o)
	if err != nil {
		return nil, err
	}

	ob := observer.New()

	var (
		nextf NextFunc
		errf  ErrFunc
		donef DoneFunc
	)

	isObserver := false

	switch h := h.(type) {
	case NextFunc:
		nextf = h
	case ErrFunc:
		errf = h
	case DoneFunc:
		donef = h
	case *BaseObserver:
		ob = h
		isObserver = true
	}

	if !isObserver {
		ob = observer.New(func(ob *Observer) {
			ob.NextHandler: nextf,
			ob.ErrHandler: errf,
			ob.DoneHandler: donef,
		})
	}

	processStream(o, ob)

	return Subscriptor(&Subscription{SubscribeAt: time.Now()}), nil
}
