package connectable

import (
	"sync"
	"time"

	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/fx"
	"github.com/jochasinga/grx/observable"
	"github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/subscription"
)

// Connectable "is" a Basic observable with a slice to keep record of subscribed observers.
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

// From creates a Connectable from a slice of interface{}
func From(items []interface{}) Connectable {
	source := make(chan interface{}, len(items))
	go func() {
		for _, item := range items {
			source <- item
		}
		close(source)
	}()
	return Connectable{Observable: source}
}

func Empty() Connectable {
	source := make(chan interface{})
	go func() {
		close(source)
	}()
	return Connectable{Observable: source}
}

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

func Start(f fx.DirectiveFunc, fs ...fx.DirectiveFunc) Connectable {
	if len(fs) > 0 {
		fs = append([]fx.DirectiveFunc{f}, fs...)
	} else {
		fs = []fx.DirectiveFunc{f}
	}

	source := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f fx.DirectiveFunc) {
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

func (co Connectable) Subscribe(handler bases.EventHandler) Connectable {
	ob := observable.CheckEventHandler(handler)
	co.observers = append(co.observers, ob)
	return co
}

func (co Connectable) Connect() <-chan subscription.Subscription {
	done := make(chan subscription.Subscription, 1)
	sub := subscription.New().Subscribe()

	for _, ob := range co.observers {
		fin := make(chan struct{})

		go func(ob observer.Observer) {

			//c := make(chan interface{}, len(co.Observable))

			/*
				for item := range co.Observable {
					c <- item
					close(c)
				}
			*/

			for item := range co.Observable {
				switch item := item.(type) {
				case error:
					ob.OnError(item)

					// Record error
					sub.Error = item
					break
				default:
					ob.OnNext(item)
				}
			}

			fin <- struct{}{}
			close(fin)
		}(ob)

		go func() {
			<-fin
			if sub.Error == nil {
				ob.OnDone()
			}
			done <- sub.Unsubscribe()
		}()
	}
	return done
}

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
