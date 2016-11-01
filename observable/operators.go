package observable

import (
	"sync"
	"time"

	"github.com/jochasinga/grx/event"
	"github.com/jochasinga/grx/observer"
)

// Add adds an Event to the Observable and return that Observable.
// myStream = myStream.Add(newEvent)
func (observable *Observable) Add(ev *event.Event) *Observable {
	go func() {
		observable.Stream <- ev
	}()
	return observable
}

// Create creates an Observable from scratch by means of a function?
func Create(obsrvr *observer.Observer) *Observable {
	obsrvble := New()
	obsrvble.Observer = obsrvr
	return obsrvble
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

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, end int) *Observable {
	observable := New(end - start)
	go func() {
		for i := start; i < end; i++ {
			observable.Stream <- &event.Event{ Value: i }
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
	observable.Observer = observer
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
				observable.Observer.OnError(ev)
				break
			case ev.Value != nil:
				observable.Observer.OnNext(ev)
			}
		}
		wg.Done()
	}(observable.Stream)

	go func() {
		wg.Wait()
		observable.Observer.OnCompleted(&event.Event{ Completed: true })
	}()
	
	// A hack for empty, finite Observable--emit a "terminal" event to signal stream's termination.
	//observer.OnCompleted(&event.Event{ Completed: true })
	return observable
}



