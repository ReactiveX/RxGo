package grx

import (
        "errors"
	"sync"
        "time"
)

// Observable is a stream of events implemented by an internal channel.
type Observable struct {
        Stream    chan interface{}

        // Pointer to a default Observer, or the one subscribed to itself.
	Observer *Observer
}

// To query a channel's length, this method is not goroutine safe
// and should block or it will return incorrect result.
func (o *Observable) isCompleted() bool {
        if len(o.Stream) > 0 {
                return false
        }
        return true
}

// hasNext determines whether there is a next object in the Observable.
func (o *Observable) hasNext() bool {
	_, ok := <-o.Stream
	return ok
}

// New constructs an empty Observable with 0 or more buffer length.
// myStream := observable.New(1)
func NewObservable(buf ...int) *Observable {
        bufferLen := 0
        if len(buf) > 0 {
                bufferLen = buf[len(buf)-1]
        }
        o := &Observable{
                Stream: make(chan interface{}, bufferLen),
                //Observer: &observer.Observer{},
		Observer: new(Observer),
        }
        o.Observer.Observable = o
        return o
}

// Create creates a new Observable provided by a function that takes an Observer
// as an argument. It is recommended to be used with SubscribeWith() method.
// o := observable.Create(func(ob *observer.Observer) {
//         ob.OnNext("Hello")
//         ob.OnError(errors.New("This is an error."))
// })
func CreateObservable(fn func(*Observer)) *Observable {
        o := NewObservable()
        go fn(o.Observer)
        return o
}

func CreateFromChannel(items chan interface{}) *Observable {
        if items != nil {
                return &Observable{
                        Stream: items,
                        Observer: new(Observer),
                }
        }
        return NewObservable()
}

// Add adds an Event to the Observable and return that Observable.
// myStream = myStream.Add(10)
func (o *Observable) Add(v interface{}) *Observable {
	if _, ok := <-o.Stream; ok {
		go func() {
			o.Stream <- v
		}()
		return o
	}
	ochan := make(chan interface{})

	go func() {
		for item := range o.Stream {
			ochan <- item
		}
		o.Stream = ochan
	}()
	return o
}

// Empty creates an Observable with one last item marked as "completed".
// myStream := observable.Empty()
func Empty() *Observable {
        o := NewObservable()
        go func() {
                close(o.Stream)
        }()
        return o
}

// Interval creates an Observable emitting incremental integers infinitely
// between each give interval.
// source := observable.Interval(1 * time.Second)
func Interval(d time.Duration) *Observable {
        o := NewObservable(1)
        i := 0
        go func() {
                for {
                        o.Stream <- i
                        <-time.After(d)
                        i++
                }
        }()
        return o
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, end int) *Observable {
        o := NewObservable(end - start)
        go func() {
                for i := start; i < end; i++ {
                        o.Stream <- i
                }
        }()
        return o
}

// Just creates an observable with only one item and emit "as-is".
// source := observable.Just("https://someurl.com/api")
func Just(items ...interface{}) *Observable {
        o := NewObservable(1)
        go func() {
		for _, item := range items {
			o.Stream <- item
		}
                close(o.Stream)
        }()
        return o
}

// From creates an Observable from a slice of items and emit them in order.
func From(items []interface{}) *Observable {
        o := NewObservable(len(items))
        go func() {
                for _, item := range items {
                        o.Stream <- item
                }
                close(o.Stream)
        }()
        return o
}

// Start creates an Observable from one or more directive-like functions
// myStream := observable.Start(f1, f2, f3)
func Start(fs ...func() interface{}) *Observable {
        o := NewObservable(len(fs))
        var wg sync.WaitGroup
        for _, f := range fs {
                wg.Add(1)
                go func(fn func() interface{}) {
                        o.Stream <-fn()
                        wg.Done()
                }(f)
        }
        go func() {
                wg.Wait()
                close(o.Stream)
        }()
        return o
}

// Subscribe subscribes an Observer to the Observable and starts it.
func (o *Observable) Subscribe(ob *Observer) (*Subscription, error) {
        if o == nil {
                return nil, errors.New("Observer is not initialized.")
        }
        if o.Stream == nil {
                return nil, errors.New("Stream is not initialized.")
        }

        var wg sync.WaitGroup
        wg.Add(1)

        go func(stream chan interface{}) {
                for item := range stream {
                        switch v := item.(type) {
                        case error:
                                if ob.ErrHandler != nil {
                                        ob.ErrHandler(v)
                                }
                                return
                        case interface{}:
                                if ob.NextHandler != nil {
                                        ob.NextHandler(v)
                                }
                        }
                }
                wg.Done()
        }(o.Stream)

        go func() {
                wg.Wait()
                if ob.DoneHandler != nil {
                        ob.DoneHandler()
                }
        }()

        return &Subscription{Subscribe: time.Now()}, nil
}

// SubscribeWith subscribes handlers to the Observable and starts it.
func (o *Observable) SubscribeFunc(nxtf func(v interface{}), errf func(e error), donef func()) (*Subscription, error) {
        if o == nil {
                return nil, errors.New("Observable is not initialized.")
        }
        if o.Stream == nil {
                return nil, errors.New("Stream is not initialized.")
        }

        var wg sync.WaitGroup
        wg.Add(1)
        go func(stream chan interface{}) {
                for item := range stream {
                        switch v := item.(type) {
                        case error:
                                if errf != nil {
                                        errf(v)
                                }
                                return
                        case interface{}:
                                if nxtf != nil {
                                        nxtf(v)
                                }
                        }
                }
                wg.Done()
        }(o.Stream)

        go func() {
                wg.Wait()
                if donef != nil {
                        donef()
                }
        }()
        return &Subscription{Subscribe: time.Now()}, nil
}

// SubscribeHandler subscribes a Handler to the Observable and starts it.
func (o *Observable) SubscribeHandler(h Handler, hs ...Handler) (*Subscription, error) {
        if o == nil {
                return nil, errors.New("Observable is not initialized.")
        }
        if o.Stream == nil {
                return nil, errors.New("Stream is not initialized.")
        }

	handlers := []Handler{h}
	handlers = append(handlers, hs...)

        nc, errc, donec := make(chan NextFunc), make(chan ErrFunc), make(chan DoneFunc)
        go func() {
		for _, h := range handlers {
			switch fn := h.(type) {
			case NextFunc:
				nc <- fn
				close(nc)
			case ErrFunc:
				errc <- fn
				close(errc)
			case DoneFunc:
				donec <- fn
				close(donec)
			}
		}
        }()
	
        var wg sync.WaitGroup
        wg.Add(1)
        go func(stream chan interface{}) {
                for item := range stream {

                        // TODO: What's the point of type switching here?
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
        }(o.Stream)

        go func() {
                wg.Wait()
                if fn, ok := <-donec; ok {
                        fn()
                }
        }()
        return &Subscription{Subscribe: time.Now()}, nil
}
