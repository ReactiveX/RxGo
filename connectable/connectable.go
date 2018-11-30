// Package connectable provides a Connectable and its methods.
package connectable

import (
	"sync"
	"time"

	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/handlers"
)

// Connectable can subscribe to several EventHandlers
// before starting processing with Connect.
type Connectable interface {
	Connect() <-chan rxgo.Observer
	Do(nextf func(interface{})) Connectable
	Subscribe(handler rxgo.EventHandler, opts ...rxgo.Option) Connectable
	Map(fn fx.Function) Connectable
	Filter(fn fx.Predicate) Connectable
	Scan(apply fx.Function2) Connectable
	First() Connectable
	Last() Connectable
	Distinct(apply fx.Function) Connectable
	DistinctUntilChanged(apply fx.Function) Connectable
}

type connector struct {
	rxgo.Observable
	observers []rxgo.Observer
}

// New creates a Connectable with optional observer(s) as parameters.
func New(buffer uint, observers ...rxgo.Observer) Connectable {
	return &connector{
		Observable: rxgo.NewObservable(buffer),
		observers:  observers,
	}
}

// From creates a Connectable from an Iterator.
func From(it rxgo.Iterator) Connectable {
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
	return &connector{
		Observable: rxgo.NewObservableFromChannel(source),
	}
}

// Empty creates a Connectable with no item and terminate immediately.
func Empty() Connectable {
	source := make(chan interface{})
	go func() {
		close(source)
	}()
	return &connector{
		Observable: rxgo.NewObservableFromChannel(source),
	}
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

	return &connector{
		Observable: rxgo.NewObservableFromChannel(source),
	}
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
	return &connector{
		Observable: rxgo.NewObservableFromChannel(source),
	}
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

	return &connector{
		Observable: rxgo.NewObservableFromChannel(source),
	}
}

// Start creates a Connectable from one or more directive-like Supplier
// and emits the result of each operation asynchronously on a new Connectable.
func Start(f fx.Supplier, fs ...fx.Supplier) Connectable {
	if len(fs) > 0 {
		fs = append([]fx.Supplier{f}, fs...)
	} else {
		fs = []fx.Supplier{f}
	}

	source := make(chan interface{})

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f fx.Supplier) {
			source <- f()
			wg.Done()
		}(f)
	}

	go func() {
		wg.Wait()
		close(source)
	}()

	return &connector{
		Observable: rxgo.NewObservableFromChannel(source),
	}
}

// Subscribe subscribes an EventHandler and returns a Connectable.
func (c *connector) Subscribe(handler rxgo.EventHandler,
	opts ...rxgo.Option) Connectable {
	ob := rxgo.CheckEventHandler(handler)
	c.observers = append(c.observers, ob)
	return c
}

// Do is like Subscribe but subscribes a func(interface{}) as a NextHandler
func (c *connector) Do(nextf func(interface{})) Connectable {
	ob := rxgo.NewObserver(handlers.NextFunc(nextf))
	c.observers = append(c.observers, ob)
	return c
}

// Connect activates the Observable stream and returns a channel of Subscription channel.
func (c *connector) Connect() <-chan rxgo.Observer {
	done := make(chan rxgo.Observer, 1)
	source := []interface{}{}

	for {
		item, err := c.Observable.Next()
		if err != nil {
			break
		}
		source = append(source, item)
	}

	var wg sync.WaitGroup
	wg.Add(len(c.observers))

	for _, ob := range c.observers {
		local := make([]interface{}, len(source))
		copy(local, source)

		fin := make(chan struct{})

		var e error
		go func(ob rxgo.Observer) {
		OuterLoop:
			for _, item := range local {
				switch item := item.(type) {
				case error:
					ob.OnError(item)

					// Record error
					e = item
					break OuterLoop
				default:
					ob.OnNext(item)
				}
			}
			fin <- struct{}{}
		}(ob)

		go func(ob rxgo.Observer) {
			<-fin
			if e == nil {
				ob.OnDone()
			}

			wg.Done()
		}(ob)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

// Map maps a Function predicate to each item in Connectable and
// returns a new Connectable with applied items.
func (c *connector) Map(fn fx.Function) Connectable {
	return &connector{
		Observable: c.Observable.Map(fn),
	}
}

// Filter filters items in the original Connectable and returns
// a new Connectable with the filtered items.
func (c *connector) Filter(fn fx.Predicate) Connectable {
	return &connector{
		Observable: c.Observable.Filter(fn),
	}
}

// Scan applies Function2 predicate to each item in the original
// Connectable sequentially and emits each successive value on a new Connectable.
func (c *connector) Scan(apply fx.Function2) Connectable {
	return &connector{
		Observable: c.Observable.Scan(apply),
	}
}

// First returns new Connectable which emits only first item.
func (c *connector) First() Connectable {
	return &connector{
		Observable: c.Observable.First(),
	}
}

// Last returns a new Connectable which emits only last item.
func (c *connector) Last() Connectable {
	return &connector{
		Observable: c.Observable.Last(),
	}
}

//Distinct suppress duplicate items in the original Connectable and
//returns a new Connectable.
func (c *connector) Distinct(apply fx.Function) Connectable {
	return &connector{
		Observable: c.Observable.Distinct(apply),
	}
}

//DistinctUntilChanged suppress duplicate items in the original Connectable only
// if they are successive to one another and returns a new Connectable.
func (c *connector) DistinctUntilChanged(apply fx.Function) Connectable {
	return &connector{
		Observable: c.Observable.DistinctUntilChanged(apply),
	}
}
