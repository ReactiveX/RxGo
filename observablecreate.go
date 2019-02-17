package rxgo

import (
	"math"
	"sync"
	"time"

	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
)

// newObservableFromChannel creates an Observable from a given channel
func newObservableFromChannel(ch chan interface{}) Observable {
	return &observable{
		iterable: newIterableFromChannel(ch),
	}
}

// newColdObservable creates a cold observable
func newColdObservable(f func(chan interface{})) Observable {
	return &observable{
		iterable: newIterableFromFunc(f),
	}
}

// newObservableFromIterable creates an Observable from a given iterable
func newObservableFromIterable(it Iterable) Observable {
	return &observable{
		iterable: it,
	}
}

// newObservableFromSlice creates an Observable from a given channel
func newObservableFromSlice(s []interface{}) Observable {
	return &observable{
		iterable: newIterableFromSlice(s),
	}
}

// newObservableFromRange creates an Observable from a range.
func newObservableFromRange(start, count int) Observable {
	return &observable{
		iterable: newIterableFromRange(start, count),
	}
}

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

	return newObservableFromChannel(out)
}

// Concat emit the emissions from two or more Observables without interleaving them
func Concat(observable1 Observable, observables ...Observable) Observable {
	out := make(chan interface{})
	go func() {
		it := observable1.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		for _, obs := range observables {
			it := obs.Iterator()
			for {
				if item, err := it.Next(); err == nil {
					out <- item
				} else {
					break
				}
			}
		}

		close(out)
	}()
	return newObservableFromChannel(out)
}

func FromSlice(s []interface{}) Observable {
	return newObservableFromSlice(s)
}

func FromChannel(ch chan interface{}) Observable {
	return newObservableFromChannel(ch)
}

func FromIterable(it Iterable) Observable {
	return newObservableFromIterable(it)
}

// From creates a new Observable from an Iterator.
func From(it Iterator) Observable {
	out := make(chan interface{})
	go func() {
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}
		close(out)
	}()
	return newObservableFromChannel(out)
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
	return newObservableFromChannel(out)
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
	return newObservableFromChannel(out)
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, count int) (Observable, error) {
	if count < 0 {
		return nil, errors.New(errors.IllegalInputError, "count must be positive")
	}
	if start+count-1 > math.MaxInt32 {
		return nil, errors.New(errors.IllegalInputError, "max value is bigger than MaxInt32")
	}

	return newObservableFromRange(start, count), nil
}

// Just creates an Observable with the provided item(s).
func Just(item interface{}, items ...interface{}) Observable {
	if len(items) > 0 {
		items = append([]interface{}{item}, items...)
	} else {
		items = []interface{}{item}
	}

	return newObservableFromSlice(items)
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

	return newObservableFromChannel(out)
}

// Never create an Observable that emits no items and does not terminate
func Never() Observable {
	out := make(chan interface{})
	return newObservableFromChannel(out)
}

// Timer returns an Observable that emits the zeroed value of a float64 after a
// specified delay, and then completes.
func Timer(d Duration) Observable {
	out := make(chan interface{})
	go func() {
		if d == nil {
			time.Sleep(0)
		} else {
			time.Sleep(d.duration())
		}
		out <- 0.
		close(out)
	}()
	return newObservableFromChannel(out)
}
