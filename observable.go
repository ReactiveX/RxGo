package rxgo

import (
	"sync"

	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/optional"
	"github.com/reactivex/rxgo/options"
)

// Observable is a basic observable interface
type Observable interface {
	Iterator
	Distinct(apply Function) Observable
	DistinctUntilChanged(apply Function) Observable
	Filter(apply Predicate) Observable
	First() Observable
	FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable
	Last() Observable
	Map(apply Function) Observable
	Scan(apply Function2) Observable
	Skip(nth uint) Observable
	SkipLast(nth uint) Observable
	Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer
	Take(nth uint) Observable
	TakeLast(nth uint) Observable
	TakeWhile(apply Predicate) Observable
	ToList() Observable
	ToMap(keySelector Function) Observable
	ToMapWithValueSelector(keySelector Function, valueSelector Function) Observable
	ZipFromObservable(publisher Observable, zipper Function2) Observable
	ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
		doneFunc handlers.DoneFunc, opts ...options.Option) Observer
	OnErrorReturn(resumeFunc ErrorFunction) Observable
	OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable
	Reduce(apply Function2) OptionalSingle
	Publish() ConnectableObservable
	Count() Single
	ElementAt(index uint) Single
	FirstOrDefault(defaultValue interface{}) Single
	LastOrDefault(defaultValue interface{}) Single
	getOnonErrorReturn() ErrorFunction
	getOnErrorResumeNext() ErrorToObservableFunction
	DefaultIfEmpty(defaultValue interface{}) Observable
}

// observable is a structure handling a channel of interface{} and implementing Observable
type observable struct {
	ch                  chan interface{}
	errorOnSubscription error
	observableFactory   func() Observable
	onErrorReturn       ErrorFunction
	onErrorResumeNext   ErrorToObservableFunction
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
func CheckEventHandler(handler handlers.EventHandler) Observer {
	return NewObserver(handler)
}

// CheckHandler checks the underlying type of an EventHandler.
func CheckEventHandlers(handler ...handlers.EventHandler) Observer {
	return NewObserver(handler...)
}

func iterate(observable Observable, observer Observer) error {
	for {
		item, err := observable.Next()
		if err != nil {
			switch err := err.(type) {
			case errors.BaseError:
				if errors.ErrorCode(err.Code()) == errors.EndOfIteratorError {
					return nil
				}
			}
		} else {
			switch item := item.(type) {
			case error:
				if observable.getOnonErrorReturn() != nil {
					observer.OnNext(observable.getOnonErrorReturn()(item))
					// Stop the subscription
					return nil
				} else if observable.getOnErrorResumeNext() != nil {
					observable = observable.getOnErrorResumeNext()(item)
				} else {
					observer.OnError(item)
					return item
				}
			default:
				observer.OnNext(item)
			}

		}
	}

	return nil
}

// Next returns the next item on the Observable.
func (o *observable) Next() (interface{}, error) {
	if next, ok := <-o.ch; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// Subscribe subscribes an EventHandler and returns a Subscription channel.
func (o *observable) Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer {
	ob := CheckEventHandler(handler)

	observableOptions := options.ParseOptions(opts...)

	if o.errorOnSubscription != nil {
		go func() {
			ob.OnError(o.errorOnSubscription)
		}()
		return ob
	}

	if observableOptions.Parallelism() == 0 {
		go func() {
			e := iterate(o, ob)
			if e == nil {
				ob.OnDone()
			}
		}()
	} else {
		wg := sync.WaitGroup{}

		var e error
		for i := 0; i < observableOptions.Parallelism(); i++ {
			wg.Add(1)

			go func() {
				e = iterate(o, ob)
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
func (o *observable) Map(apply Function) Observable {
	out := make(chan interface{})

	var it Observable = o
	if o.observableFactory != nil {
		it = o.observableFactory()
	}

	go func() {
		for {
			item, err := it.Next()
			if err != nil {
				break
			}
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
func (o *observable) Filter(apply Predicate) Observable {
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
func (o *observable) Distinct(apply Function) Observable {
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
func (o *observable) DistinctUntilChanged(apply Function) Observable {
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
func (o *observable) Scan(apply Function2) Observable {
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

func (o *observable) Reduce(apply Function2) OptionalSingle {
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
func (o *observable) TakeWhile(apply Predicate) Observable {
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

// ToMap convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function
func (o *observable) ToMap(keySelector Function) Observable {
	out := make(chan interface{})
	go func() {
		m := make(map[interface{}]interface{})
		for item := range o.ch {
			m[keySelector(item)] = item
		}
		out <- m
		close(out)
	}()
	return &observable{ch: out}
}

// ToMapWithValueSelector convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function and valued by another
// value function
func (o *observable) ToMapWithValueSelector(keySelector Function, valueSelector Function) Observable {
	out := make(chan interface{})
	go func() {
		m := make(map[interface{}]interface{})
		for item := range o.ch {
			m[keySelector(item)] = valueSelector(item)
		}
		out <- m
		close(out)
	}()
	return &observable{ch: out}
}

// ZipFromObservable che emissions of multiple Observables together via a specified function
// and emit single items for each combination based on the results of this function
func (o *observable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	out := make(chan interface{})
	go func() {
	OuterLoop:
		for item1 := range o.ch {
			for {
				item2, err := publisher.Next()
				if err != nil {
					break
				}
				out <- zipper(item1, item2)
				continue OuterLoop
			}
			break OuterLoop
		}
		close(out)
	}()
	return &observable{ch: out}
}

// ForEach subscribes to the Observable and receives notifications for each element.
func (o *observable) ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
	doneFunc handlers.DoneFunc, opts ...options.Option) Observer {
	return o.Subscribe(CheckEventHandlers(nextFunc, errFunc, doneFunc), opts...)
}

// Publish returns a ConnectableObservable which waits until its connect method
// is called before it begins emitting items to those Observers that have subscribed to it.
func (o *observable) Publish() ConnectableObservable {
	return NewConnectableObservable(o)
}

// OnErrorReturn instructs an Observable to emit an item (returned by a specified function)
// rather than invoking onError if it encounters an error.
func (o *observable) OnErrorReturn(resumeFunc ErrorFunction) Observable {
	o.onErrorReturn = resumeFunc
	o.onErrorResumeNext = nil
	return o
}

// OnErrorResumeNext Instructs an Observable to pass control to another Observable rather than invoking
// onError if it encounters an error.
func (o *observable) OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable {
	o.onErrorResumeNext = resumeSequence
	o.onErrorReturn = nil
	return o
}

func (o *observable) getOnonErrorReturn() ErrorFunction {
	return o.onErrorReturn
}

func (o *observable) getOnErrorResumeNext() ErrorToObservableFunction {
	return o.onErrorResumeNext
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (o *observable) DefaultIfEmpty(defaultValue interface{}) Observable {
	out := make(chan interface{})
	go func() {
		empty := true
		for item := range o.ch {
			empty = false
			out <- item
		}
		if empty {
			out <- defaultValue
		}
		close(out)
	}()
	return &observable{ch: out}
}
