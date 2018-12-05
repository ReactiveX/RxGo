package rxgo

import (
	"math"
	"sync"
	"time"

	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
)

func isClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// Creates observable from based on source function. Keep it mind to call emitter.OnDone()
// to signal sequence's end.
// Example:
// - emitting none elements
// observable.Create(emitter observer.Observer, disposed bool) { emitter.OnDone() })
// - emitting one element
// observable.Create(func(emitter observer.Observer, disposed bool) {
//		emitter.OnNext("one element")
//		emitter.OnDone()
// })
func Create(source func(emitter Observer, disposed bool)) Observable {
	out := make(chan interface{})
	emitter := NewObserver(
		handlers.NextFunc(func(el interface{}) {
			if !isClosed(out) {
				out <- el
			}
		}), handlers.ErrFunc(func(err error) {
			// decide how to deal with errors
			if !isClosed(out) {
				close(out)
			}
		}), handlers.DoneFunc(func() {
			if !isClosed(out) {
				close(out)
			}
		}),
	)

	go func() {
		source(emitter, isClosed(out))
	}()

	return NewObservableFromChannel(out)
}

// Concat emit the emissions from two or more Observables without interleaving them
func Concat(observable1 Observable, observables ...Observable) Observable {
	out := make(chan interface{})
	go func() {
		it := observable1.Iterator()
		for it.Next() {
			item := it.Value()
			out <- item
		}

		for _, obs := range observables {
			it := obs.Iterator()
			for it.Next() {
				item := it.Value()
				out <- item
			}
		}

		close(out)
	}()
	return NewObservableFromChannel(out)
}

// Defer waits until an observer subscribes to it, and then it generates an Observable.
func Defer(f func() Observable) Observable {
	return &observable{
		observableFactory: f,
	}
}

// From creates a new Observable from an Iterator.
func From(it Iterator) Observable {
	out := make(chan interface{})
	go func() {
		for it.Next() {
			item := it.Value()
			out <- item
		}
		close(out)
	}()
	return NewObservableFromChannel(out)
}

// Error returns an Observable that invokes an Observer's onError method
// when the Observer subscribes to it.
func Error(err error) Observable {
	return &observable{
		errorOnSubscription: err,
	}
}

// Empty creates an Observable with no item and terminate immediately.
func Empty() Observable {
	out := make(chan interface{})
	go func() {
		close(out)
	}()
	return NewObservableFromChannel(out)
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(term chan struct{}, interval time.Duration) Observable {
	out := make(chan interface{})
	go func(term chan struct{}) {
		i := 0
	OuterLoop:
		for {
			select {
			case <-term:
				break OuterLoop
			case <-time.After(interval):
				out <- i
			}
			i++
		}
		close(out)
	}(term)
	return NewObservableFromChannel(out)
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, count int) (Observable, error) {
	if count < 0 {
		return nil, errors.New(errors.IllegalInputError, "count must be positive")
	}
	if start+count-1 > math.MaxInt32 {
		return nil, errors.New(errors.IllegalInputError, "max value is bigger than MaxInt32")
	}

	out := make(chan interface{})
	go func() {
		i := start
		for i < count+start {
			out <- i
			i++
		}
		close(out)
	}()
	return NewObservableFromChannel(out), nil
}

// Just creates an Observable with the provided item(s).
func Just(item interface{}, items ...interface{}) Observable {
	out := make(chan interface{})
	if len(items) > 0 {
		items = append([]interface{}{item}, items...)
	} else {
		items = []interface{}{item}
	}

	go func() {
		for _, item := range items {
			out <- item
		}
		close(out)
	}()

	return NewObservableFromChannel(out)
}

// Start creates an Observable from one or more directive-like Supplier
// and emits the result of each operation asynchronously on a new Observable.
func Start(f Supplier, fs ...Supplier) Observable {
	if len(fs) > 0 {
		fs = append([]Supplier{f}, fs...)
	} else {
		fs = []Supplier{f}
	}

	out := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f Supplier) {
			out <- f()
			wg.Done()
		}(f)
	}

	// Wait in another goroutine to not block
	go func() {
		wg.Wait()
		close(out)
	}()

	return NewObservableFromChannel(out)
}

// Never create an Observable that emits no items and does not terminate
func Never() Observable {
	out := make(chan interface{})
	go func() {
		select {}
	}()
	return NewObservableFromChannel(out)
}
