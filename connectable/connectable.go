// Package connectable provides a Connectable and its methods.
package connectable

import (
	"sync"
	"time"

	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
	"github.com/reactivex/rxgo/subscription"
)

// Connectable is an Observable which can subscribe several
// EventHandlers before starting processing with Connect.
type Connectable struct {
	observable.Observable
	observers []observer.Observer
}

// New creates a Connectable with optional observer(s) as parameters.
func New(buffer uint, observers ...observer.Observer) Connectable {
	return Connectable{
		Observable: make(observable.Observable, int(buffer)),
		observers:  observers,
	}
}

// From creates a Connectable from an Iterator.
func From(it rx.Iterator) Connectable {
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
	return Connectable{Observable: source}
}

// Empty creates a Connectable with no item and terminate immediately.
func Empty() Connectable {
	source := make(chan interface{})
	go func() {
		close(source)
	}()
	return Connectable{Observable: source}
}

// Interval creates a Connectable emitting incremental integers infinitely between
// each given time interval.
func Interval(term chan struct{}, timeout time.Duration) Connectable {
	source := make(chan interface{})
	go func(term chan struct{}) {
		i := 0
	OuterLoop:
		for {
			select {
			case <-term:
				break OuterLoop
			case <-time.After(timeout):
				source <- i
			}
			i++
		}
		close(source)
	}(term)

	return Connectable{Observable: source}
}

// Range creates an Connectable that emits a particular range of sequential integers.
func Range(start, end int) Connectable {
	source := make(chan interface{})
	go func() {
		i := start
		for i < end {
			source <- i
			i++
		}
		close(source)
	}()
	return Connectable{Observable: source}
}

// Just creates an Connectable with the provided item(s).
func Just(item interface{}, items ...interface{}) Connectable {
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

	return Connectable{Observable: source}
}

// Start creates a Connectable from one or more directive-like EmittableFunc
// and emits the result of each operation asynchronously on a new Connectable.
func Start(f fx.EmittableFunc, fs ...fx.EmittableFunc) Connectable {
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

	go func() {
		wg.Wait()
		close(source)
	}()

	return Connectable{Observable: source}
}

// Subscribe subscribes an EventHandler and returns a Connectable.
func (co Connectable) Subscribe(handler rx.EventHandler) Connectable {
	ob := observable.CheckEventHandler(handler)
	co.observers = append(co.observers, ob)
	return co
}

// Do is like Subscribe but subscribes a func(interface{}) as a NextHandler
func (co Connectable) Do(nextf func(interface{})) Connectable {
	ob := observer.Observer{NextHandler: nextf}
	co.observers = append(co.observers, ob)
	return co
}

// Connect activates the Observable stream and returns a channel of Subscription channel.
func (co Connectable) Connect() <-chan (chan subscription.Subscription) {
	done := make(chan (chan subscription.Subscription), 1)
	source := []interface{}{}

	for item := range co.Observable {
		source = append(source, item)
	}

	var wg sync.WaitGroup
	wg.Add(len(co.observers))

	for _, ob := range co.observers {
		local := make([]interface{}, len(source))
		copy(local, source)

		fin := make(chan struct{})
		sub := subscription.New().Subscribe()

		go func(ob observer.Observer) {
		OuterLoop:
			for _, item := range local {
				switch item := item.(type) {
				case error:
					ob.OnError(item)

					// Record error
					sub.Error = item
					break OuterLoop
				default:
					ob.OnNext(item)
				}
			}
			fin <- struct{}{}
		}(ob)

		temp := make(chan subscription.Subscription)

		go func(ob observer.Observer) {
			<-fin
			if sub.Error == nil {
				ob.OnDone()
				sub.Unsubscribe()
			}

			go func() {
				temp <- sub
				done <- temp
			}()
			wg.Done()
		}(ob)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

// Map maps a MappableFunc predicate to each item in Connectable and
// returns a new Connectable with applied items.
func (co Connectable) Map(fn fx.MappableFunc) Connectable {
	source := make(chan interface{}, len(co.Observable))
	go func() {
		for item := range co.Observable {
			source <- fn(item)
		}
		close(source)
	}()
	return Connectable{Observable: source}
}

// Filter filters items in the original Connectable and returns
// a new Connectable with the filtered items.
func (co Connectable) Filter(fn fx.FilterableFunc) Connectable {
	source := make(chan interface{}, len(co.Observable))
	go func() {
		for item := range co.Observable {
			if fn(item) {
				source <- item
			}
		}
		close(source)
	}()
	return Connectable{Observable: source}
}

// Scan applies ScannableFunc predicate to each item in the original
// Connectable sequentially and emits each successive value on a new Connectable.
func (co Connectable) Scan(apply fx.ScannableFunc) Connectable {
	out := make(chan interface{})

	go func() {
		var current interface{}
		for item := range co.Observable {
			out <- apply(current, item)
			current = apply(current, item)
		}
		close(out)
	}()
	return Connectable{Observable: out}
}

// First returns new Connectable which emits only first item.
func (co Connectable) First() Connectable {
	out := make(chan interface{})
	go func() {
		for item := range co.Observable {
			out <- item
			break
		}
		close(out)
	}()
	return Connectable{Observable: out}
}

// Last returns a new Connectable which emits only last item.
func (co Connectable) Last() Connectable {
	out := make(chan interface{})
	go func() {
		var last interface{}
		for item := range co.Observable {
			last = item
		}
		out <- last
		close(out)
	}()
	return Connectable{Observable: out}
}

//Distinct suppress duplicate items in the original Connectable and
//returns a new Connectable.
func (co Connectable) Distinct(apply fx.KeySelectorFunc) Connectable {
	out := make(chan interface{})
	go func() {
		keysets := make(map[interface{}]struct{})
		for item := range co.Observable {
			key := apply(item)
			_, ok := keysets[key]
			if !ok {
				out <- item
			}
			keysets[key] = struct{}{}
		}
		close(out)
	}()
	return Connectable{Observable: out}
}

//DistinctUntilChanged suppress duplicate items in the original Connectable only
// if they are successive to one another and returns a new Connectable.
func (co Connectable) DistinctUntilChanged(apply fx.KeySelectorFunc) Connectable {
	out := make(chan interface{})
	go func() {
		var current interface{}
		for item := range co.Observable {
			key := apply(item)
			if current != key {
				out <- item
				current = key
			}
		}
		close(out)
	}()
	return Connectable{Observable: out}
}
