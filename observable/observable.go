package observable

import (
	"sync"
	"time"

	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/observer"
	"github.com/reactivex/rxgo/subscription"
)

// Observable is a basic observable interface
type Observable interface {
	rx.Iterator

	Distinct(apply fx.KeySelectorFunc) Observable
	DistinctUntilChanged(apply fx.KeySelectorFunc) Observable
	Filter(apply fx.FilterableFunc) Observable
	First() Observable
	FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable
	Last() Observable
	Map(apply fx.MappableFunc) Observable
	Scan(apply fx.ScannableFunc) Observable
	Skip(nth uint) Observable
	SkipLast(nth uint) Observable
	Subscribe(handler rx.EventHandler, opts ...Option) <-chan subscription.Subscription
	Take(nth uint) Observable
	TakeLast(nth uint) Observable

	ElementAt(index uint) Single
}

// Single is similar to an Observable but emits only one single element or an error notification.
type Single interface {
	Observable
}

// observator is a structure handling a channel of interface{} and implementing Observable
type observator struct {
	ch                  chan interface{}
	errorOnSubscription error
}

// New creates an Observable
func New(buffer uint) Observable {
	ch := make(chan interface{}, int(buffer))
	return &observator{
		ch: ch,
	}
}

// NewFromChannel creates an Observable from a given channel
func NewFromChannel(ch chan interface{}) Observable {
	return &observator{
		ch: ch,
	}
}

// CheckHandler checks the underlying type of an EventHandler.
func CheckEventHandler(handler rx.EventHandler) observer.Observer {
	return observer.New(handler)
}

// Next returns the next item on the Observable.
func (o *observator) Next() (interface{}, error) {
	if next, ok := <-o.ch; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// Subscribe subscribes an EventHandler and returns a Subscription channel.
func (o *observator) Subscribe(handler rx.EventHandler, opts ...Option) <-chan subscription.Subscription {
	done := make(chan subscription.Subscription)
	sub := subscription.New().Subscribe()

	ob := CheckEventHandler(handler)

	// Parse options
	var observableOptions options
	for _, opt := range opts {
		opt.apply(&observableOptions)
	}

	if o.errorOnSubscription != nil {
		go func() {
			ob.OnError(o.errorOnSubscription)
			close(done)
		}()
		return done
	}

	if observableOptions.parallelism == 0 {
		go func() {
		OuterLoop:
			for item := range o.ch {
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
	} else {
		wg := sync.WaitGroup{}

		for i := 0; i < observableOptions.parallelism; i++ {
			wg.Add(1)

			go func() {
			OuterLoop:
				for item := range o.ch {
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
				wg.Done()
			}()
		}

		go func() {
			wg.Wait()
			// OnDone only gets executed if there's no error.
			if sub.Error == nil {
				ob.OnDone()
			}

			done <- sub.Unsubscribe()
		}()
	}

	return done
}

/*
func (o *observable) Unsubscribe() subscription.Subscription {
	// Stub: to be implemented
	return subscription.New()
}
*/

// Map maps a MappableFunc predicate to each item in Observable and
// returns a new Observable with applied items.
func (o *observator) Map(apply fx.MappableFunc) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			out <- apply(item)
		}
		close(out)
	}()
	return &observator{ch: out}
}

func (o *observator) ElementAt(index uint) Single {
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
	return &observator{ch: out}
}

// Take takes first n items in the original Obserable and returns
// a new Observable with the taken items.
func (o *observator) Take(nth uint) Observable {
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
	return &observator{ch: out}
}

// TakeLast takes last n items in the original Observable and returns
// a new Observable with the taken items.
func (o *observator) TakeLast(nth uint) Observable {
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
	return &observator{ch: out}
}

// Filter filters items in the original Observable and returns
// a new Observable with the filtered items.
func (o *observator) Filter(apply fx.FilterableFunc) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			if apply(item) {
				out <- item
			}
		}
		close(out)
	}()
	return &observator{ch: out}
}

// First returns new Observable which emit only first item.
func (o *observator) First() Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			out <- item
			break
		}
		close(out)
	}()
	return &observator{ch: out}
}

// Last returns a new Observable which emit only last item.
func (o *observator) Last() Observable {
	out := make(chan interface{})
	go func() {
		var last interface{}
		for item := range o.ch {
			last = item
		}
		out <- last
		close(out)
	}()
	return &observator{ch: out}
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *observator) Distinct(apply fx.KeySelectorFunc) Observable {
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
	return &observator{ch: out}
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original
// Observable and returns a new Observable.
func (o *observator) DistinctUntilChanged(apply fx.KeySelectorFunc) Observable {
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
	return &observator{ch: out}
}

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observator) Skip(nth uint) Observable {
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
	return &observator{ch: out}
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observator) SkipLast(nth uint) Observable {
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
	return &observator{ch: out}
}

// Scan applies ScannableFunc predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o *observator) Scan(apply fx.ScannableFunc) Observable {
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
	return &observator{ch: out}
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
	return &observator{ch: source}
}

// Error returns an Observable that invokes an Observer's onError method
// when the Observer subscribes to it.
func Error(err error) Observable {
	return &observator{
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
	return &observator{ch: source}
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
	return &observator{ch: source}
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
		return &observator{ch: source}
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
		return &observator{ch: source}
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
	return &observator{ch: source}
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

	return &observator{ch: source}
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

	return &observator{ch: source}
}
