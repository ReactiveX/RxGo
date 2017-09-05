package observable

import (
	"reflect"
	"sync"
	"time"

	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/observer"
	"github.com/reactivex/rxgo/subscription"
)

// Observable is a basic observable channel
type Observable <-chan interface{}

var DefaultObservable = make(Observable)

var none = new(int)

// New creates an Observable
func New(buffer uint) Observable {
	return make(Observable, int(buffer))
}

// Next returns the next item on the Observable.
func (o Observable) Next() (interface{}, error) {
	if next, ok := <-o; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// Subscribe subscribes an EventHandler and returns a Subscription channel.
func (o Observable) Subscribe(handler interface{}) <-chan subscription.Subscription {
	done := make(chan subscription.Subscription)
	sub := subscription.New().Subscribe()

	ob := observer.New(handler)

	go func() {
	OuterLoop:
		for item := range o {
			switch item := item.(type) {
			case error:
				ob.OnError(item)

				// Record the error and break the loop.
				sub.Error = item
				break OuterLoop
			default:
				ob.OnNext(item)
			}
		}

		// OnDone only gets executed if there's no error.
		if sub.Error == nil {
			ob.OnDone()
		}

		done <- sub.Unsubscribe()
		return
	}()

	return done
}

/*
func (o Observable) Unsubscribe() subscription.Subscription {
	// Stub: to be implemented
	return subscription.New()
}
*/

// Map maps a MappableFunc predicate to each item in Observable and
// returns a new Observable with applied items.
func (o Observable) Map(apply fx.MappableFunc) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o {
			out <- apply(item)
		}
		close(out)
	}()
	return Observable(out)
}

// Take takes first n items in the original Obserable and returns
// a new Observable with the taken items.
func (o Observable) Take(nth uint) Observable {
	out := make(chan interface{})
	go func() {
		takeCount := 0
		for item := range o {
			if takeCount < int(nth) {
				takeCount += 1
				out <- item
				continue
			}
			break
		}
		close(out)
	}()
	return Observable(out)
}

// TakeLast takes last n items in the original Observable and returns
// a new Observable with the taken items.
func (o Observable) TakeLast(nth uint) Observable {
	out := make(chan interface{})
	go func() {
		buf := make([]interface{}, nth)
		for item := range o {
			if len(buf) >= int(nth) {
				buf = buf[1:]
			}
			buf = append(buf, item)
		}
		for _, takenItem := range buf {
			out <- takenItem
		}
		close(out)
	}()
	return Observable(out)
}

// Filter filters items in the original Observable and returns
// a new Observable with the filtered items.
func (o Observable) Filter(apply fx.FilterableFunc) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o {
			if apply(item) {
				out <- item
			}
		}
		close(out)
	}()
	return Observable(out)
}

// First returns new Observable which emit only first item.
func (o Observable) First() Observable {
	out := make(chan interface{})
	go func() {
		for item := range o {
			out <- item
			break
		}
		close(out)
	}()
	return Observable(out)
}

// Last returns a new Observable which emit only last item.
func (o Observable) Last() Observable {
	out := make(chan interface{})
	go func() {
		var last interface{}
		for item := range o {
			last = item
		}
		out <- last
		close(out)
	}()
	return Observable(out)
}

//Distinct supress duplicate items in the original Observable and returns
// a new Observable.
func (o Observable) Distinct(apply fx.KeySelectorFunc) Observable {
	out := make(chan interface{})
	go func() {
		keysets := make(map[interface{}]struct{})
		for item := range o {
			key := apply(item)
			_, ok := keysets[key]
			if !ok {
				out <- item
			}
			keysets[key] = struct{}{}
		}
		close(out)
	}()
	return Observable(out)
}

func (o Observable) DistinctUntilChanged(apply fx.KeySelectorFunc) Observable {
	out := make(chan interface{})
	go func() {
		var current interface{}
		for item := range o {
			key := apply(item)
			if current != key {
				out <- item
				current = key
			}
		}
		close(out)
	}()
	return Observable(out)
}

// Scan applies ScannableFunc predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o Observable) Scan(apply fx.ScannableFunc) Observable {
	out := make(chan interface{})

	go func() {
		var current interface{}
		for item := range o {
			out <- apply(current, item)
			current = apply(current, item)
		}
		close(out)
	}()
	return Observable(out)
}

// From creates a new Observable from an Iterator.
func From(it rx.Iterator) Observable {
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
	return Observable(source)
}

// Empty creates an Observable with no item and terminate immediately.
func Empty() Observable {
	source := make(chan interface{})
	go func() {
		close(source)
	}()
	return Observable(source)
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
	return Observable(source)
}

// Range creates an Observable that emits a particular range of sequential integers.
func Range(start, end int) Observable {
	source := make(chan interface{})
	go func() {
		i := start
		for i < end {
			source <- i
			i++
		}
		close(source)
	}()
	return Observable(source)
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

	return Observable(source)
}

// Start creates an Observable from one or more directive-like EmittableFunc
// and emits the result of each operation asynchronously on a new Observable.
func Start(f fx.EmittableFunc, fs ...fx.EmittableFunc) Observable {
	if len(fs) > 0 {
		fs = append([]fx.EmittableFunc{f}, fs...)
	} else {
		fs = []fx.EmittableFunc{f}
	}

	source := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f fx.EmittableFunc) {
			source <- f()
			wg.Done()
		}(f)
	}

	// Wait in another goroutine to not block
	go func() {
		wg.Wait()
		close(source)
	}()

	return Observable(source)
}

// Combine multiple Observables into one by merging their emissions
func Merge(o1 Observable, o2 Observable, on ...Observable) Observable {
	out := make(chan interface{})
	go func() {
		chans := append([]Observable{o1, o2}, on...)
		count := len(chans)
		cases := make([]reflect.SelectCase, count)
		for i := range cases {
			cases[i].Dir = reflect.SelectRecv
			cases[i].Chan = reflect.ValueOf(chans[i])
		}
		for count > 0 {
			chosen, recv, recvOk := reflect.Select(cases)
			if recvOk {
				out <- recv.Interface()
			} else {
				cases[chosen].Chan = reflect.ValueOf(nil)
				count--
			}
		}
		close(out)
	}()
	return Observable(out)
}

//CombineLatest emits an item whenever any of the source Observables emits an item
func CombineLatest(o []Observable, apply fx.CombinableFunc) Observable {
	out := make(chan interface{})
	go func() {
		chans := o
		count := len(chans)
		left := len(chans)
		is := make([]interface{}, len(chans))
		for i := 0; i < len(is); i++ {
			is[i] = none
		}
		cases := make([]reflect.SelectCase, count)
		for i := range cases {
			cases[i].Dir = reflect.SelectRecv
			cases[i].Chan = reflect.ValueOf(chans[i])
		}
		for count > 0 {
			chosen, recv, recvOk := reflect.Select(cases)
			if recvOk {
				if is[chosen] == none {
					left--
				}
				is[chosen] = recv.Interface()
				if left == 0 {
					out <- apply(is)
				}
			} else {
				cases[chosen].Chan = reflect.ValueOf(nil)
				count--
			}
		}
		close(out)
	}()
	return Observable(out)
}
