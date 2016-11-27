package grx

import (
	"fmt"
        "errors"
	"sync"
        "time"
)

// Observable is a stream of events implemented by an internal channel.
type Observable struct {
        C    chan interface{}
	
        // Pointer to a default Observer which subscribed to itself.
	observer *Observer
	done chan struct{}
}

func (o *Observable) isDone() bool {
	if _, ok := <-o.done; ok {
		return true
	}
	return false
}

func (o *Observable) hasNext() bool {
	if _, ok := <-o.done; ok {
		return false
	}
	return true
}

// NewObservable constructs an empty Observable with 0 or more buffer length.
// myStream := observable.New(1)
func NewObservable(buf ...int) *Observable {
        bufferLen := 0
        if len(buf) > 0 {
                bufferLen = buf[len(buf)-1]
        }
        o := &Observable{
                C: make(chan interface{}, bufferLen),
		observer: new(Observer),
		done: make(chan struct{}, 1),
		
        }
        o.observer.observable = o
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
        go fn(o.observer)
        return o
}

func CreateFromChannel(items chan interface{}) *Observable {
        if items != nil {
		o := NewObservable()
		o.C = items
		return o
        }
        return NewObservable()
}

// Add adds an item to the Observable and return that Observable.
// If the Observable is done, it creates a new one and return it.
// myStream = myStream.Add(10)
func (o *Observable) Add(v interface{}) *Observable {
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
func Empty() *Observable {
        o := NewObservable()
        go func() {
                close(o.C)
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
                        o.C <- i
                        <-time.After(d)
                        i++
                }
        }()
        return o
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, end int) *Observable {
        o := NewObservable(0)
        go func() {
                for i := start; i < end; i++ {
                        o.C <- i
                }
		close(o.C)
        }()
        return o
}

// Just creates an observable with only one item and emit "as-is".
// source := observable.Just("https://someurl.com/api")
func Just(items ...interface{}) *Observable {
        o := NewObservable(1)
        go func() {
		for _, item := range items {
			o.C <- item
		}
                close(o.C)
        }()
        return o
}

// From creates an Observable from a slice of items and emit them in order.
func From(items []interface{}) *Observable {
        o := NewObservable(len(items))
        go func() {
                for _, item := range items {
                        o.C <- item
                }
                close(o.C)
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
                        o.C <-fn()
                        wg.Done()
                }(f)
        }
        go func() {
                wg.Wait()
                close(o.C)
        }()
        return o
}

// Subscribe subscribes an Observer to the Observable and starts it.
func (o *Observable) Subscribe(ob *Observer) (*Subscription, error) {
        if o == nil {
                return nil, errors.New("Observer is not initialized.")
        }
        if o.C == nil {
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
        }(o.C)

        go func() {
                wg.Wait()
		o.done <- struct{}{}
        }()

	go func() {
		select {
		case recent := <-o.done:

			// Clone to a new o.done channel so others can read from.
			o.done = make(chan struct{}, 1)
			o.done <- recent
			if ob.DoneHandler != nil {
				ob.DoneHandler()
			}
		}
	}()

        return &Subscription{Subscribe: time.Now()}, nil
}

// SubscribeWith subscribes handlers to the Observable and starts it.
func (o *Observable) SubscribeFunc(nxtf func(v interface{}), errf func(e error), donef func()) (*Subscription, error) {
        if o == nil {
                return nil, errors.New("Observable is not initialized.")
        }
        if o.C == nil {
                return nil, errors.New("Stream is not initialized.")
        }

	go func() {
		select {
		case <-o.done:
			fmt.Println("Done! Firing DoneHandler...")
			if donef != nil {
				donef()
			}
		}
	}()

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
        }(o.C)

        go func() {
                wg.Wait()
		o.done <- struct{}{}
        }()

	go func() {
		select {
		case recent := <-o.done:

			// Clone to a new o.done channel so others can read from.
			o.done = make(chan struct{}, 1)
			o.done <- recent
			if donef != nil {
				donef()
			}
		}
	}()
	
        return &Subscription{Subscribe: time.Now()}, nil
}

// SubscribeHandler subscribes a Handler to the Observable and starts it.
func (o *Observable) SubscribeHandler(h Handler, hs ...Handler) (*Subscription, error) {
        if o == nil {
                return nil, errors.New("Observable is not initialized.")
        }
        if o.C == nil {
                return nil, errors.New("Stream is not initialized.")
        }

	handlers := []Handler{h}
	handlers = append(handlers, hs...)

        nc, errc, dc := make(chan NextFunc), make(chan ErrFunc), make(chan DoneFunc)
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
				dc <- fn
				close(dc)
			}
		}
        }()
	
        var wg sync.WaitGroup
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
        return &Subscription{Subscribe: time.Now()}, nil
}
