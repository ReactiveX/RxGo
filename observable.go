package rxgo

import (
	"sync"
	"time"

	"fmt"

	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/optional"
	"github.com/reactivex/rxgo/options"
)

// Observable is a basic observable interface
type Observable interface {
	Iterator
	All(predicate Predicate) Single
	AverageFloat32() Single
	AverageFloat64() Single
	AverageInt() Single
	AverageInt8() Single
	AverageInt16() Single
	AverageInt32() Single
	AverageInt64() Single
	BufferWithCount(count, skip int) Observable
	BufferWithTime(timespan, timeshift Duration) Observable
	BufferWithTimeOrCount(timespan Duration, count int) Observable
	Contains(equal Predicate) Single
	Count() Single
	DefaultIfEmpty(defaultValue interface{}) Observable
	Distinct(apply Function) Observable
	DistinctUntilChanged(apply Function) Observable
	DoOnEach(onNotification Consumer) Observable
	ElementAt(index uint) Single
	Filter(apply Predicate) Observable
	First() Observable
	FirstOrDefault(defaultValue interface{}) Single
	FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable
	ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
		doneFunc handlers.DoneFunc, opts ...options.Option) Observer
	Last() Observable
	LastOrDefault(defaultValue interface{}) Single
	Map(apply Function) Observable
	Max(comparator Comparator) OptionalSingle
	Min(comparator Comparator) OptionalSingle
	OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable
	OnErrorReturn(resumeFunc ErrorFunction) Observable
	Publish() ConnectableObservable
	Reduce(apply Function2) OptionalSingle
	Repeat(count int64, frequency Duration) Observable
	Scan(apply Function2) Observable
	Skip(nth uint) Observable
	SkipLast(nth uint) Observable
	SkipWhile(apply Predicate) Observable
	StartWithItems(items ...interface{}) Observable
	StartWithIterable(iterable iterable.Iterable) Observable
	StartWithObservable(observable Observable) Observable
	Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer
	Take(nth uint) Observable
	TakeLast(nth uint) Observable
	TakeWhile(apply Predicate) Observable
	ToList() Observable
	ToMap(keySelector Function) Observable
	ToMapWithValueSelector(keySelector Function, valueSelector Function) Observable
	ZipFromObservable(publisher Observable, zipper Function2) Observable
	getOnErrorResumeNext() ErrorToObservableFunction
	getOnErrorReturn() ErrorFunction
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
				if observable.getOnErrorReturn() != nil {
					observer.OnNext(observable.getOnErrorReturn()(item))
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
		results := make([]chan error, 0)
		for i := 0; i < observableOptions.Parallelism(); i++ {
			ch := make(chan error)
			go func() {
				ch <- iterate(o, ob)
			}()
			results = append(results, ch)
		}

		go func() {
			for _, ch := range results {
				err := <-ch
				if err != nil {
					return
				}
			}

			ob.OnDone()
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

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
func (o *observable) SkipWhile(apply Predicate) Observable {
	out := make(chan interface{})
	go func() {
		skip := true
		for item := range o.ch {
			if !skip {
				out <- item
			} else {
				if !apply(item) {
					out <- item
					skip = false
				}
			}
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

func (o *observable) All(predicate Predicate) Single {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			if !predicate(item) {
				out <- false
				close(out)
				return
			}
		}
		out <- true
		close(out)
	}()
	return NewSingleFromChannel(out)
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

func (o *observable) getOnErrorReturn() ErrorFunction {
	return o.onErrorReturn
}

func (o *observable) getOnErrorResumeNext() ErrorToObservableFunction {
	return o.onErrorResumeNext
}

// Contains returns an Observable that emits a Boolean that indicates whether
// the source Observable emitted an item (the comparison is made against a predicate).
func (o *observable) Contains(equal Predicate) Single {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			if equal(item) {
				out <- true
				close(out)
				return
			}
		}
		out <- false
		close(out)
	}()
	return NewSingleFromChannel(out)
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

// DoOnEach operator allows you to establish a callback that the resulting Observable
// will call each time it emits an item
func (o *observable) DoOnEach(onNotification Consumer) Observable {
	out := make(chan interface{})
	go func() {
		for item := range o.ch {
			out <- item
			onNotification(item)
		}
		close(out)
	}()
	return &observable{ch: out}
}

// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
// at most count times, at a particular frequency.
func (o *observable) Repeat(count int64, frequency Duration) Observable {
	out := make(chan interface{})

	if count != Indefinitely {
		if count < 0 {
			count = 0
		}
	}

	go func() {
		persist := make([]interface{}, 0)
		for item := range o.ch {
			out <- item
			persist = append(persist, item)
		}

		for {
			if count != Indefinitely {
				if count == 0 {
					break
				}
			}

			if frequency != nil {
				time.Sleep(frequency.duration())
			}

			for _, v := range persist {
				out <- v
			}

			count = count - 1
		}
		close(out)
	}()
	return &observable{ch: out}
}

// AverageInt calculates the average of numbers emitted by an Observable and emits this average int.
func (o *observable) AverageInt() Single {
	out := make(chan interface{})
	go func() {
		sum := 0
		count := 0
		for item := range o.ch {
			if v, ok := item.(int); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// AverageInt8 calculates the average of numbers emitted by an Observable and emits this average int8.
func (o *observable) AverageInt8() Single {
	out := make(chan interface{})
	go func() {
		var sum int8 = 0
		var count int8 = 0
		for item := range o.ch {
			if v, ok := item.(int8); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// AverageInt16 calculates the average of numbers emitted by an Observable and emits this average int16.
func (o *observable) AverageInt16() Single {
	out := make(chan interface{})
	go func() {
		var sum int16 = 0
		var count int16 = 0
		for item := range o.ch {
			if v, ok := item.(int16); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// AverageInt32 calculates the average of numbers emitted by an Observable and emits this average int32.
func (o *observable) AverageInt32() Single {
	out := make(chan interface{})
	go func() {
		var sum int32 = 0
		var count int32 = 0
		for item := range o.ch {
			if v, ok := item.(int32); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// AverageInt64 calculates the average of numbers emitted by an Observable and emits this average int64.
func (o *observable) AverageInt64() Single {
	out := make(chan interface{})
	go func() {
		var sum int64 = 0
		var count int64 = 0
		for item := range o.ch {
			if v, ok := item.(int64); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits this average float32.
func (o *observable) AverageFloat32() Single {
	out := make(chan interface{})
	go func() {
		var sum float32 = 0
		var count float32 = 0
		for item := range o.ch {
			if v, ok := item.(float32); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// AverageFloat64 calculates the average of numbers emitted by an Observable and emits this average float64.
func (o *observable) AverageFloat64() Single {
	out := make(chan interface{})
	go func() {
		var sum float64 = 0
		var count float64 = 0
		for item := range o.ch {
			if v, ok := item.(float64); ok {
				sum = sum + v
				count = count + 1
			} else {
				out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
				close(out)
				return
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}()
	return NewSingleFromChannel(out)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (o *observable) Max(comparator Comparator) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		empty := true
		var max interface{} = nil
		for item := range o.ch {
			empty = false

			if max == nil {
				max = item
			} else {
				if comparator(max, item) == Smaller {
					max = item
				}
			}
		}
		if empty {
			out <- optional.Empty()
		} else {
			out <- optional.Of(max)
		}
		close(out)
	}()
	return &optionalSingle{ch: out}
}

// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
func (o *observable) Min(comparator Comparator) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		empty := true
		var min interface{} = nil
		for item := range o.ch {
			empty = false

			if min == nil {
				min = item
			} else {
				if comparator(min, item) == Greater {
					min = item
				}
			}
		}
		if empty {
			out <- optional.Empty()
		} else {
			out <- optional.Of(min)
		}
		close(out)
	}()
	return &optionalSingle{ch: out}
}

// BufferWithCount returns an Observable that emits buffers of items it collects
// from the source Observable.
// The resulting Observable emits buffers every skip items, each containing a slice of count items.
// When the source Observable completes or encounters an error,
// the resulting Observable emits the current buffer and propagates
// the notification from the source Observable.
func (o *observable) BufferWithCount(count, skip int) Observable {
	out := make(chan interface{})
	go func() {
		if count <= 0 {
			out <- errors.New(errors.IllegalInputError, "count must be positive")
			close(out)
			return
		}

		if skip <= 0 {
			out <- errors.New(errors.IllegalInputError, "skip must be positive")
			close(out)
			return
		}

		buffer := make([]interface{}, count, count)
		iCount := 0
		iSkip := 0
		for item := range o.ch {
			switch item := item.(type) {
			case error:
				if iCount != 0 {
					out <- buffer[:iCount]
				}
				out <- item
				close(out)
				return
			default:
				if iCount >= count { // Skip
					iSkip++
				} else { // Add to buffer
					buffer[iCount] = item
					iCount++
					iSkip++
				}

				if iSkip == skip { // Send current buffer
					out <- buffer
					buffer = make([]interface{}, count, count)
					iCount = 0
					iSkip = 0
				}
			}
		}

		if iCount != 0 {
			out <- buffer[:iCount]
		}

		close(out)
	}()
	return &observable{ch: out}
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (o *observable) BufferWithTime(timespan, timeshift Duration) Observable {
	out := make(chan interface{})
	go func() {
		if timespan == nil || timespan.duration() == 0 {
			out <- errors.New(errors.IllegalInputError, "timespan must not be nil")
			close(out)
			return
		}

		if timeshift == nil {
			timeshift = WithDuration(0)
		}

		var mux sync.Mutex
		var listenMutex sync.Mutex
		buffer := make([]interface{}, 0)
		stop := false
		listen := true

		// First goroutine in charge to check the timespan
		go func() {
			for {
				time.Sleep(timespan.duration())
				mux.Lock()
				if !stop {
					out <- buffer
					buffer = make([]interface{}, 0)
					mux.Unlock()

					if timeshift.duration() != 0 {
						listenMutex.Lock()
						listen = false
						listenMutex.Unlock()
						time.Sleep(timeshift.duration())
						listenMutex.Lock()
						listen = true
						listenMutex.Unlock()
					}
				} else {
					mux.Unlock()
					return
				}
			}
		}()

		// Second goroutine in charge to retrieve the items from the source observable
		go func() {
			for item := range o.ch {
				switch item := item.(type) {
				case error:
					mux.Lock()
					if len(buffer) > 0 {
						out <- buffer
					}
					out <- item
					close(out)
					stop = true
					mux.Unlock()
					return
				default:
					listenMutex.Lock()
					l := listen
					listenMutex.Unlock()

					mux.Lock()
					if l {
						buffer = append(buffer, item)
					}
					mux.Unlock()
				}
			}
			mux.Lock()
			if len(buffer) > 0 {
				out <- buffer
			}
			close(out)
			stop = true
			mux.Unlock()
		}()

	}()
	return &observable{ch: out}
}

// BufferWithTimeOrCount returns an Observable that emits buffers of items it collects
// from the source Observable. The resulting Observable emits connected,
// non-overlapping buffers, each of a fixed duration specified by the timespan argument
// or a maximum size specified by the count argument (whichever is reached first).
// When the source Observable completes or encounters an error, the resulting Observable
// emits the current buffer and propagates the notification from the source Observable.
func (o *observable) BufferWithTimeOrCount(timespan Duration, count int) Observable {
	out := make(chan interface{})
	go func() {
		if timespan == nil || timespan.duration() == 0 {
			out <- errors.New(errors.IllegalInputError, "timespan must not be nil")
			close(out)
			return
		}

		if count <= 0 {
			out <- errors.New(errors.IllegalInputError, "count must be positive")
			close(out)
			return
		}

		sendCh := make(chan []interface{})
		errCh := make(chan error)
		buffer := make([]interface{}, 0)
		var bufferMutex sync.Mutex

		// First sender goroutine
		go func() {
			for {
				select {
				case currentBuffer := <-sendCh:
					out <- currentBuffer
				case error := <-errCh:
					if len(buffer) > 0 {
						out <- buffer
					}
					if error != nil {
						out <- error
					}
					close(out)
					return
				case <-time.After(timespan.duration()): // Send on timer
					bufferMutex.Lock()
					b := make([]interface{}, len(buffer))
					copy(b, buffer)
					buffer = make([]interface{}, 0)
					bufferMutex.Unlock()

					out <- b
				}
			}
		}()

		// Second goroutine in charge to retrieve the items from the source observable
		go func() {
			for item := range o.ch {
				switch item := item.(type) {
				case error:
					errCh <- item
					return
				default:
					bufferMutex.Lock()
					buffer = append(buffer, item)
					if len(buffer) >= count {
						b := make([]interface{}, len(buffer))
						copy(b, buffer)
						buffer = make([]interface{}, 0)
						bufferMutex.Unlock()

						sendCh <- b
					} else {
						bufferMutex.Unlock()
					}
				}
			}
			errCh <- nil
		}()

	}()
	return &observable{ch: out}
}

// StartWithItems returns an Observable that emits the specified items before it begins to emit items emitted
// by the source Observable.
func (o *observable) StartWithItems(items ...interface{}) Observable {
	out := make(chan interface{})
	go func() {
		for _, item := range items {
			out <- item
		}

		for item := range o.ch {
			out <- item
		}

		close(out)
	}()
	return &observable{ch: out}
}

// StartWithIterable returns an Observable that emits the items in a specified Iterable before it begins to
// emit items emitted by the source Observable.
func (o *observable) StartWithIterable(iterable iterable.Iterable) Observable {
	out := make(chan interface{})
	go func() {
		for {
			item, err := iterable.Next()
			if err != nil {
				break
			}
			out <- item
		}

		for item := range o.ch {
			out <- item
		}

		close(out)
	}()
	return &observable{ch: out}
}

// StartWithObservable returns an Observable that emits the items in a specified Observable before it begins to
// emit items emitted by the source Observable.
func (o *observable) StartWithObservable(obs Observable) Observable {
	out := make(chan interface{})
	go func() {
		for {
			item, err := obs.Next()
			if err != nil {
				break
			}
			out <- item
		}

		for item := range o.ch {
			out <- item
		}

		close(out)
	}()
	return &observable{ch: out}
}
