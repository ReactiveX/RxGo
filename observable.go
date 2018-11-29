package rx

import (
	"sync"
	"time"

	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/optional"
)

// Observable is a basic observable interface
type Observable interface {
	Iterator

	Distinct(apply fx.Function) Observable
	DistinctUntilChanged(apply fx.Function) Observable
	Filter(apply fx.Predicate) Observable
	First() Observable
	FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable
	Last() Observable
	Map(apply fx.Function) Observable
	Scan(apply fx.Function2) Observable
	Skip(nth uint) Observable
	SkipLast(nth uint) Observable
	Subscribe(handler EventHandler, opts ...Option) Observer
	Take(nth uint) Observable
	TakeLast(nth uint) Observable
	TakeWhile(apply fx.Predicate) Observable
	ToList() Observable

	Reduce(apply fx.Function2) OptionalSingle

	Count() Single
	ElementAt(index uint) Single
	FirstOrDefault(defaultValue interface{}) Single
	LastOrDefault(defaultValue interface{}) Single
}

// observable is a structure handling a channel of interface{} and implementing Observable
type observable struct {
	ch                  chan interface{}
	errorOnSubscription error
}

// NewObservable creates an Observable
func NewObservable(buffer uint) Observable {
	ch := make(chan interface{}, int(buffer))
	return &observable{
		ch: ch,
	}
}

// NewObservableFromChannel creates an Observable from a given channel
func NewObservableFromChannel(ch chan interface{}) Observable {
	return &observable{
		ch: ch,
	}
}

// CheckHandler checks the underlying type of an EventHandler.
func CheckEventHandler(handler EventHandler) Observer {
	return NewObserver(handler)
}

// Next returns the next item on the Observable.
func (o *observable) Next() (interface{}, error) {
	if next, ok := <-o.ch; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// Subscribe subscribes an EventHandler and returns a Subscription channel.
func (o *observable) Subscribe(handler EventHandler, opts ...Option) Observer {
	ob := CheckEventHandler(handler)

	// Parse options
	var observableOptions options
	for _, opt := range opts {
		opt.apply(&observableOptions)
	}

	if o.errorOnSubscription != nil {
		go func() {
			ob.OnError(o.errorOnSubscription)
		}()
		return ob
	}

	if observableOptions.parallelism == 0 {
		go func() {
			var e error
		OuterLoop:
			for item := range o.ch {
				switch item := item.(type) {
				case error:
					ob.OnError(item)

					// Record the error and break the loop.
					e = item
					break OuterLoop
				default:
					ob.OnNext(item)
				}
			}

			// OnDone only gets executed if there's no error.
			if e == nil {
				ob.OnDone()
			}

			return
		}()
	} else {
		wg := sync.WaitGroup{}

		var e error
		for i := 0; i < observableOptions.parallelism; i++ {
			wg.Add(1)

			go func() {
			OuterLoop:
				for item := range o.ch {
					switch item := item.(type) {
					case error:
						ob.OnError(item)

						// Record the error and break the loop.
						e = item
						break OuterLoop
					default:
						ob.OnNext(item)
					}
				}
				wg.Done()
			}()
		}

		go func() {
			wg.Wait()
			// OnDone only gets executed if there's no error.
			if e == nil {
				ob.OnDone()
			}
		}()
	}

	return ob
}

// Map maps a Function predicate to each item in Observable and
// returns a new Observable with applied items.
func (o *observable) Map(apply fx.Function) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			out <- apply(item)
		}
		close(out)
	}()
	return &observable{ch: out}
}

/*
func (o *observable) Unsubscribe() subscription.Subscription {
	// Stub: to be implemented
	return subscription.New()
}
*/

func (o *observable) ElementAt(index uint) Single {
	out := make(chan interface{})
	go func() {
		takeCount := 0
		for item := range o.ch {
			if takeCount == int(index) {
				out <- item
				close(out)
				return
			}
			takeCount += 1
		}
		out <- errors.New(errors.ElementAtError)
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// Take takes first n items in the original Obserable and returns
// a new Observable with the taken items.
func (o *observable) Take(nth uint) Observable {
	out := make(chan interface{})
	go func() {
		takeCount := 0
		for item := range o.ch {
			if takeCount < int(nth) {
				takeCount += 1
				out <- item
				continue
			}
			break
		}
		close(out)
	}()
	return &observable{ch: out}
}

// TakeLast takes last n items in the original Observable and returns
// a new Observable with the taken items.
func (o *observable) TakeLast(nth uint) Observable {
	out := make(chan interface{})
	go func() {
		buf := make([]interface{}, nth)
		for item := range o.ch {
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
	return &observable{ch: out}
}

// Filter filters items in the original Observable and returns
// a new Observable with the filtered items.
func (o *observable) Filter(apply fx.Predicate) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			if apply(item) {
				out <- item
			}
		}
		close(out)
	}()
	return &observable{ch: out}
}

// First returns new Observable which emit only first item.
func (o *observable) First() Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			out <- item
			break
		}
		close(out)
	}()
	return &observable{ch: out}
}

// Last returns a new Observable which emit only last item.
func (o *observable) Last() Observable {
	out := make(chan interface{})
	go func() {
		var last interface{}
		for item := range o.ch {
			last = item
		}
		out <- last
		close(out)
	}()
	return &observable{ch: out}
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *observable) Distinct(apply fx.Function) Observable {
	out := make(chan interface{})
	go func() {
		keysets := make(map[interface{}]struct{})
		for item := range o.ch {
			key := apply(item)
			_, ok := keysets[key]
			if !ok {
				out <- item
			}
			keysets[key] = struct{}{}
		}
		close(out)
	}()
	return &observable{ch: out}
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original
// Observable and returns a new Observable.
func (o *observable) DistinctUntilChanged(apply fx.Function) Observable {
	out := make(chan interface{})
	go func() {
		var current interface{}
		for item := range o.ch {
			key := apply(item)
			if current != key {
				out <- item
				current = key
			}
		}
		close(out)
	}()
	return &observable{ch: out}
}

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observable) Skip(nth uint) Observable {
	out := make(chan interface{})
	go func() {
		skipCount := 0
		for item := range o.ch {
			if skipCount < int(nth) {
				skipCount += 1
				continue
			}
			out <- item
		}
		close(out)
	}()
	return &observable{ch: out}
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observable) SkipLast(nth uint) Observable {
	out := make(chan interface{})
	go func() {
		buf := make(chan interface{}, nth)
		for item := range o.ch {
			select {
			case buf <- item:
			default:
				out <- (<-buf)
				buf <- item
			}
		}
		close(buf)
		close(out)
	}()
	return &observable{ch: out}
}

// Scan applies Function2 predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o *observable) Scan(apply fx.Function2) Observable {
	out := make(chan interface{})

	go func() {
		var current interface{}
		for item := range o.ch {
			tmp := apply(current, item)
			out <- tmp
			current = tmp
		}
		close(out)
	}()
	return &observable{ch: out}
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

// Repeat creates an Observable emitting a given item repeatedly
func Repeat(item interface{}, ntimes ...int) Observable {
	source := make(chan interface{})

	// this is the infinity case no ntime parameter is given
	if len(ntimes) == 0 {
		go func() {
			for {
				source <- item
			}
		}()
		return &observable{ch: source}
	}

	// this repeat the item ntime
	if len(ntimes) > 0 {
		count := ntimes[0]
		if count <= 0 {
			return Empty()
		}
		go func() {
			for i := 0; i < count; i++ {
				source <- item
			}
			close(source)
		}()
		return &observable{ch: source}
	}

	return Empty()
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
	return &observable{ch: source}
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
func Start(f fx.Supplier, fs ...fx.Supplier) Observable {
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

	// Wait in another goroutine to not block
	go func() {
		wg.Wait()
		close(source)
	}()

	return &observable{ch: source}
}

func (o *observable) Reduce(apply fx.Function2) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		var acc interface{}
		empty := true
		for item := range o.ch {
			empty = false
			acc = apply(acc, item)
		}
		if empty {
			out <- optional.Empty()
		} else {
			out <- optional.Of(acc)
		}
		close(out)
	}()
	return NewOptionalSingleFromChannel(out)
}

func (o *observable) Count() Single {
	out := make(chan interface{})
	go func() {
		var count int64
		for range o.ch {
			count++
		}
		out <- count
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) FirstOrDefault(defaultValue interface{}) Single {
	out := make(chan interface{})
	go func() {
		first := defaultValue
		for item := range o.ch {
			first = item
			break
		}
		out <- first
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// Last returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) LastOrDefault(defaultValue interface{}) Single {
	out := make(chan interface{})
	go func() {
		last := defaultValue
		for item := range o.ch {
			last = item
		}
		out <- last
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// TakeWhile emits items emitted by an Observable as long as the
// specified condition is true, then skip the remainder.
func (o *observable) TakeWhile(apply fx.Predicate) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			if apply(item) {
				out <- item
				continue
			}
			break
		}
		close(out)
	}()
	return &observable{ch: out}
}

// ToList collects all items from an Observable and emit them as a single List.
func (o *observable) ToList() Observable {
	out := make(chan interface{})
	go func() {
		s := make([]interface{}, 0)
		for item := range o.ch {
			s = append(s, item)
		}
		out <- s
		close(out)
	}()
	return &observable{ch: out}
}
