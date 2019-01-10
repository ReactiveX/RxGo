package rxgo

import (
	"math"
	"sync"
	"time"

	"github.com/reactivex/rxgo/v2/errors"
	"github.com/reactivex/rxgo/v2/handlers"
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
	emitted := make(chan interface{})
	emitter := NewObserver(
		handlers.NextFunc(func(el interface{}) {
			if !isClosed(emitted) {
				emitted <- el
			}
		}), handlers.ErrFunc(func(err error) {
			// decide how to deal with errors
			if !isClosed(emitted) {
				close(emitted)
			}
		}), handlers.DoneFunc(func() {
			if !isClosed(emitted) {
				close(emitted)
			}
		}),
	)

	go func() {
		source(emitter, isClosed(emitted))
	}()

	return &observable{
		ch: emitted,
	}
}

// Concat emit the emissions from two or more Observables without interleaving them
func Concat(observable1 Observable, observables ...Observable) Observable {
	source := make(chan interface{})
	go func() {
	OuterLoop:
		for {
			item, err := observable1.Next()
			if err != nil {
				switch err := err.(type) {
				case errors.BaseError:
					if errors.ErrorCode(err.Code()) == errors.EndOfIteratorError {
						break OuterLoop
					}
				}
			} else {
				source <- item
			}
		}

		for _, it := range observables {
		OuterLoop2:
			for {
				item, err := it.Next()
				if err != nil {
					switch err := err.(type) {
					case errors.BaseError:
						if errors.ErrorCode(err.Code()) == errors.EndOfIteratorError {
							break OuterLoop2
						}
					}
				} else {
					source <- item
				}
			}
		}

		close(source)
	}()
	return &observable{ch: source}
}

// Defer waits until an observer subscribes to it, and then it generates an Observable.
func Defer(f func() Observable) Observable {
	return &observable{
		ch:                nil,
		observableFactory: f,
	}
}

// From creates a new Observable from an Iterator.
func From(it Iterator) Observable {
	source := make(chan interface{})
	go func() {
		for {
			val, err := it.Next()
			if err != nil {
				break
			}
			source <- val
		}
		close(source)
	}()
	return &observable{ch: source}
}

// Error returns an Observable that invokes an Observer's onError method
// when the Observer subscribes to it.
func Error(err error) Observable {
	return &observable{
		ch:                  nil,
		errorOnSubscription: err,
	}
}

// Empty creates an Observable with no item and terminate immediately.
func Empty() Observable {
	source := make(chan interface{})
	go func() {
		close(source)
	}()
	return &observable{ch: source}
}

// Interval creates an Observable emitting incremental integers infinitely between
// each given time interval.
func Interval(term chan struct{}, interval time.Duration) Observable {
	source := make(chan interface{})
	go func(term chan struct{}) {
		i := 0
	OuterLoop:
		for {
			select {
			case <-term:
				break OuterLoop
			case <-time.After(interval):
				source <- i
			}
			i++
		}
		close(source)
	}(term)
	return &observable{ch: source}
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, count int) (Observable, error) {
	if count < 0 {
		return nil, errors.New(errors.IllegalInputError, "count must be positive")
	}
	if start+count-1 > math.MaxInt32 {
		return nil, errors.New(errors.IllegalInputError, "max value is bigger than MaxInt32")
	}

	source := make(chan interface{})
	go func() {
		i := start
		for i < count+start {
			source <- i
			i++
		}
		close(source)
	}()
	return &observable{ch: source}, nil
}

// Just creates an Observable with the provided item(s).
func Just(item interface{}, items ...interface{}) Observable {
	source := make(chan interface{})
	if len(items) > 0 {
		items = append([]interface{}{item}, items...)
	} else {
		items = []interface{}{item}
	}

	go func() {
		for _, item := range items {
			source <- item
		}
		close(source)
	}()

	return &observable{ch: source}
}

// Start creates an Observable from one or more directive-like Supplier
// and emits the result of each operation asynchronously on a new Observable.
func Start(f Supplier, fs ...Supplier) Observable {
	if len(fs) > 0 {
		fs = append([]Supplier{f}, fs...)
	} else {
		fs = []Supplier{f}
	}

	source := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f Supplier) {
			source <- f()
			wg.Done()
		}(f)
	}

	// Wait in another goroutine to not block
	go func() {
		wg.Wait()
		close(source)
	}()

	return &observable{ch: source}
}

// Never create an Observable that emits no items and does not terminate
func Never() Observable {
	source := make(chan interface{})
	go func() {
		select {}
	}()
	return &observable{ch: source}
}
