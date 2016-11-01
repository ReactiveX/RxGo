package observable

import (
	"sync"
	"time"

	"github.com/jochasinga/grx/event"
	"github.com/jochasinga/grx/observer"
)

// To query a channel's length, this method should block.
func (observable *Observable) isCompleted() bool {
	if len(observable.Stream) > 0 {
		return false
	}
	return true
}

// New constructs an empty Observable with 0 or more buffer length.
// myStream := observable.New(1)
func New(buf ...int) *Observable {
	bufferLen := 0
	if len(buf) > 0 {
		bufferLen = buf[len(buf)-1]
	}
	return &Observable{ Stream: make(chan *event.Event, bufferLen) }
}

// Add adds an Event to the Observable and return that Observable.
// myStream = myStream.Add(newEvent)
func (observable *Observable) Add(ev *event.Event) *Observable {
	go func() {
		observable.Stream <- ev
	}()
	return observable
}

// Empty creates an Observable with one last item marked as "completed".
// myStream := observable.Empty()
func Empty() *Observable {
	observable := New(1)
	go func() {
		observable.Stream <- &event.Event{ Completed: true }
		close(observable.Stream)
	}()
	return observable
}

// Interval creates an Observable emitting incremental integers infinitely between each give interval.
// myStream := observable.Interval(1 * time.Second)
func Interval(d time.Duration) *Observable {
	observable := New(1)
	i := 0
	go func() {
		for {
			observable.Stream <- &event.Event{ Value: i }
			<-time.After(d)
			i++
		}
	}()
	return observable
}

// Just creates an observable with only one item and emit "as-is".
// myStream := observable.Just("https://someurl.com/api")
func Just(item interface{}) *Observable {
	observable := New(1)
	go func() {
		observable.Stream <- &event.Event{ Value: item }
		close(observable.Stream)
	}()
	return observable
}

// From creates an observable from a slice of items and emit them in order.
func From(items []interface{}) *Observable {
	observable := New(len(items))
	go func() {
		for _, item := range items {
			observable.Stream <- &event.Event{ Value: item }
		}
		close(observable.Stream)
	}()
	return observable
}

// Start creates an Observable from one or more directive-like functions
// myStream := observable.Start(f1, f2, f3)
func Start(fs ...func() *event.Event) *Observable {
	observable := New(len(fs))
	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f func() *event.Event) {
			observable.Stream <-f()
			wg.Done()
		}(f)
	}
	go func() {
		wg.Wait()
		close(observable.Stream)
	}()
	return observable
}

// Subscribe subscribes an Observer to the Observable and starts it.
func (observable *Observable) Subscribe(observer *observer.Observer) *Observable {
	if observable.Stream == nil {
		return observable
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// Loop over the Observable's channel
	go func(stream chan *event.Event) {
		for ev := range stream {
			switch {
			case ev.Error != nil:
				observer.OnError(ev)
				break
			case ev.Value != nil:
				observer.OnNext(ev)
			}
		}
		wg.Done()
	}(observable.Stream)

	go func() {
		wg.Wait()
		observer.OnCompleted(&event.Event{ Completed: true })
		
	}()
	
	// A hack for empty, finite Observable--emit a "terminal" event to signal stream's termination.
	//observer.OnCompleted(&event.Event{ Completed: true })
	return observable
}



