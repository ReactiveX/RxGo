package grx

import (
	//"fmt"
	"sync"
	"time"
)

// BaseObservable is a stream of events implemented by an internal channel.
type BaseObservable struct {
	C            chan interface{}
	unsubscribed chan struct{}
	done         chan struct{}
	_observer    Observer
}

//====================== Utility functions ===========================//
func (o *BaseObservable) isDone() bool {
	if _, ok := <-o.done; ok {
		return true
	}
	return false
}

func (o *BaseObservable) Next() (interface{}, error) {
	if _, ok := <-o.done; !ok {
		return <-o.C, nil
	}
	return nil, EndOfIteratorError
}

func (o *BaseObservable) HasNext() bool {
	if o.isDone() {
		return true
	}
	return false
}

/*
 * These unexported methods are used by the higher-level <Observable>
 * interface to call, and should not be called directly by the type
 * <BaseObservable> itself unless unavoidable.
 *
 */
func (o *BaseObservable) getInnerObserver() Observer {
	return o._observer
}

func (o *BaseObservable) setInnerObserverTo(ob Observer) Observer {
	o._observer = ob
	return o._observer
}

func (o *BaseObservable) addItem(item interface{}) {
	o.Add(item)
	//o.C <- item
}

func (o *BaseObservable) terminate() {
	o.done <- struct{}{}
	close(o.C)
}

func (o *BaseObservable) unsubscribe() {
	o.unsubscribed <- struct{}{}
}

func (o *BaseObservable) getC() chan interface{} {
	return o.C
}

func NewBaseObservable(buf ...int) *BaseObservable {
	bufferLen := 0
	if len(buf) > 0 {
		bufferLen = buf[len(buf)-1]
	}
	o := &BaseObservable{
		C: make(chan interface{}, bufferLen),
		_observer: (Observer)(&BaseObserver{
			NextHandler: (NextFunc)(nil),
			ErrHandler:  (ErrFunc)(nil),
			DoneHandler: (DoneFunc)(nil),
		}),
		unsubscribed: make(chan struct{}, 1),
		done:         make(chan struct{}, 1),
	}
	//o._observer.observable = o
	if o._observer != nil {
		_ = o._observer.setInnerObservableTo(Observable(o))
	}
	return o
}

// Create creates a new Observable provided by a function that takes an Observer
// as an argument. It is recommended to be used with SubscribeWith() method.
// o := observable.Create(func(ob *observer.Observer) {
//         ob.OnNext("Hello")
//         ob.OnError(errors.New("This is an error."))
// })
//func CreateObservable(fn func(*Observer)) *Observable {
func CreateBaseObservable(fn func(Observer)) *BaseObservable {
	//o := NewObservable()
	o := NewBaseObservable()
	fn(o._observer)
	return o
}

//func CreateFromChannel(items chan interface{}) *Observable {
/*
func CreateBaseObservableFromChan(items chan interface{}) *BaseObservable {
	if items != nil {
		//o := NewObservable()
		o := NewBaseObservable()
		o.C = items
		return o
	}
	//return NewObservable()
	return NewBaseObservable()
}
*/

/*
func CreateBaseObservableFrom(iter Iterator) *BaseObservable {
	if iter != nil {

	}
}
*/

// Add adds an item to the Observable and return that Observable.
// If the Observable is done, it creates a new one and return it.
// myStream = myStream.Add(10)
//func (o *Observable) Add(v interface{}) *Observable {
func (o *BaseObservable) Add(v interface{}) *BaseObservable {
	if !o.isDone() {
		go func() {
			o.C <- v
		}()
		return o
	}

	// else if it's done (C is closed), create a fresh channel with copies
	// of the old channel's elements for o.Stream.
	ochan := make(chan interface{})
	go func() {
		for item := range o.C {
			ochan <- item
		}
		o.C = ochan
	}()
	return o
}

// Empty creates an Observable with one last item marked as "completed".
// myStream := observable.Empty()
//func Empty() *Observable {
func Empty() *BaseObservable {
	//o := NewObservable()
	o := NewBaseObservable()
	go func() {
		o.terminate()
		//close(o.C)
	}()
	return o
}

// Interval creates an Observable emitting incremental integers infinitely
// between each give interval.
// source := observable.Interval(1 * time.Second)
//func Interval(d time.Duration) *Observable {
func Interval(d time.Duration) *BaseObservable {
	//o := NewObservable(1)
	o := NewBaseObservable(1)
	i := 0
	go func() {
		for {
			o.C <- i
			<-time.After(d)
			i++
		}
	}()
	return o
}

// Range creates an Observable that emits a particular range of sequential integers.
//func Range(start, end int) *Observable {
func Range(start, end int) *BaseObservable {
	//o := NewObservable(0)
	o := NewBaseObservable()
	go func() {
		for i := start; i < end; i++ {
			o.C <- i
		}
		o.terminate()
		//close(o.C)
	}()
	return o
}

// Just creates an observable with only one item and emit "as-is".
// source := observable.Just("https://someurl.com/api")
//func Just(items ...interface{}) *Observable {
func Just(items ...interface{}) *BaseObservable {
	o := NewBaseObservable()
	go func() {
		for _, item := range items {
			o.C <- item
		}
		o.terminate()
	}()
	return o
}

// From creates an Observable from a slice of items and emit them in order.
//func From(items []interface{}) *Observable {
func From(items []interface{}) *BaseObservable {
	//o := NewObservable(len(items))
	o := NewBaseObservable(len(items))
	go func() {
		for _, item := range items {
			o.C <- item
		}
		//close(o.C)
		o.terminate()

	}()
	return o
}

// Start creates an Observable from one or more directive-like functions
// myStream := observable.Start(f1, f2, f3)
func Start(fs ...func() interface{}) *BaseObservable {
	o := NewBaseObservable(len(fs))
	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(fn func() interface{}) {
			o.C <- fn()
			wg.Done()
		}(f)
	}

	go func() {
		wg.Wait()
		//close(o.C)
		o.terminate()
	}()
	return o
}

func (o *BaseObservable) runStream(ob Observer) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(items chan interface{}) {
		for item := range items {
			switch next := item.(type) {
			case error:
				if ob.OnError != nil {
					ob.OnError(next)
				}
				return
			default:
				if ob.OnNext != nil {
					ob.OnNext(next)
				}
			}
		}
		wg.Done()
	}(o.C)

	go func() {
		wg.Wait()
		o.terminate()
	}()

	go func() {
		select {
		case <-o.unsubscribed:
		case <-o.done:
			if ob.OnDone != nil {
				ob.OnDone()
			}
		}
	}()
}

func checkObservable(o *BaseObservable) error {
	if o == nil {
		return NilObservableError
	}
	if o.C == nil {
		return NilObservableCError
	}
	if o.isDone() {
		return EndOfIteratorError
	}
	return nil

	/*
		switch {
		default:
			return UndefinedObservableError
		case o == nil:
			return NilObservableError
		case o.C == nil:
			return NilObservableCError
		case o.isDone():
			return EndOfIteratorError
		}

		return nil
	*/
}

/*
// Subscribe subscribes an Observer to the Observable and starts it.
//func (o *Observable) Subscribe(ob *Observer) (*Subscription, error) {
func (o *BaseObservable) Subscribe(ob Observer) (Subscriptor, error) {

	err := checkObservable(o)
	if err != nil {
		return nil, err
	}

	o.runStream(ob)

	return Subscriptor(&Subscription{
		_observable:   o,
		SubscribeAt:   time.Now(),
		UnsubscribeAt: time.Time{},
	}), nil
}

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
//func (o *BaseObservable) SubscribeHandler(h EventHandler, hs ...EventHandler) (Subscriptor, error) {
func (o *BaseObservable) Subscribe(h EventHandler) (Subscriptor, error) {

	/*
		err := checkObservable(o)
		if err != nil {
			return nil, err
		}
	*/

	ob := NewBaseObserver()

	var nextf NextFunc
	var errf ErrFunc
	var donef DoneFunc

	switch handler := h.(type) {
	case NextFunc:
		nextf = handler
	case ErrFunc:
		errf = handler
	case DoneFunc:
		donef = handler
	case *BaseObserver:
		ob = handler
	}

	if ob == nil {
		ob = &BaseObserver{
			NextHandler: nextf,
			ErrHandler:  errf,
			DoneHandler: donef,
		}
	}

	o.runStream(ob)

	/*
		handlers := append([]EventHandler{h}, hs...)

		nc, errc, dc := make(chan NextFunc), make(chan ErrFunc), make(chan DoneFunc)

		var wg sync.WaitGroup
		for _, handler := range handlers {
			wg.Add(1)
			go func() {
				switch handler := handler.(type) {
				case NextFunc:
					nc <- handler
				case ErrFunc:
					errc <- handler
				case DoneFunc:
					dc <- handler
				case *BaseObserver:
					switch {
					case handler.NextHandler != nil:
						nc <- handler.NextHandler
					case handler.ErrHandler != nil:
						errc <- handler.ErrHandler
					case handler.DoneHandler != nil:
						dc <- handler.DoneHandler
					}
				}
				wg.Done()
			}()
		}

		go func() {
			wg.Wait()
			close(nc)
			close(errc)
			close(dc)
		}()

		//var wg sync.WaitGroup
		wg.Add(1)

		go func(stream chan interface{}) {
			for item := range stream {
				switch v := item.(type) {
				case error:
					if fn, ok := <-errc; ok {
						fn(v)
					}
					return
				case interface{}:
					if fn, ok := <-nc; ok {
						fn(v)
					}
				}
			}
			wg.Done()
		}(o.C)

		go func() {
			wg.Wait()
			if fn, ok := <-dc; ok {
				o.done <- struct{}{}
				fn()
			}
		}()
	*/
	return Subscriptor(&Subscription{SubscribeAt: time.Now()}), nil
}
