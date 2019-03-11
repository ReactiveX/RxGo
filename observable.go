package rxgo

import (
	"sync"
	"time"

	"fmt"

	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/optional"
	"github.com/reactivex/rxgo/options"
)

// Observable is a basic observable interface
type Observable interface {
	Iterable
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
	SequenceEqual(obs Observable) Single
	Skip(nth uint) Observable
	SkipLast(nth uint) Observable
	SkipWhile(apply Predicate) Observable
	StartWithItems(items ...interface{}) Observable
	StartWithIterable(iterable Iterable) Observable
	StartWithObservable(observable Observable) Observable
	Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer
	SumFloat32() Single
	SumFloat64() Single
	SumInt64() Single
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
	iterable            Iterable
	errorOnSubscription error
	observableFactory   func() Observable
	onErrorReturn       ErrorFunction
	onErrorResumeNext   ErrorToObservableFunction
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
	it := observable.Iterator()
	for {
		if item, err := it.Next(); err == nil {
			switch item := item.(type) {
			case error:
				if observable.getOnErrorReturn() != nil {
					observer.OnNext(observable.getOnErrorReturn()(item))
					// Stop the subscription
					return nil
				} else if observable.getOnErrorResumeNext() != nil {
					observable = observable.getOnErrorResumeNext()(item)
					it = observable.Iterator()
				} else {
					observer.OnError(item)
					return item
				}
			default:
				observer.OnNext(item)
			}
		} else {
			break
		}
	}
	return nil
}

func (o *observable) Iterator() Iterator {
	return o.iterable.Iterator()
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
	f := func(out chan interface{}) {
		it := o.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- apply(item)
			} else {
				break
			}
		}
		close(out)
	}

	return newColdObservable(f)
}

func (o *observable) ElementAt(index uint) Single {
	f := func(out chan interface{}) {
		takeCount := 0

		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if takeCount == int(index) {
					out <- item
					close(out)
					return
				}
				takeCount += 1
			} else {
				break
			}
		}
		out <- errors.New(errors.ElementAtError)
		close(out)
	}
	return newColdSingle(f)
}

// Take takes first n items in the original Obserable and returns
// a new Observable with the taken items.
func (o *observable) Take(nth uint) Observable {
	f := func(out chan interface{}) {
		takeCount := 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if takeCount < int(nth) {
					takeCount += 1
					out <- item
					continue
				}
				break
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// TakeLast takes last n items in the original Observable and returns
// a new Observable with the taken items.
func (o *observable) TakeLast(nth uint) Observable {
	f := func(out chan interface{}) {
		buf := make([]interface{}, nth)
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if len(buf) >= int(nth) {
					buf = buf[1:]
				}
				buf = append(buf, item)
			} else {
				break
			}
		}
		for _, takenItem := range buf {
			out <- takenItem
		}
		close(out)
	}
	return newColdObservable(f)
}

// Filter filters items in the original Observable and returns
// a new Observable with the filtered items.
func (o *observable) Filter(apply Predicate) Observable {
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if apply(item) {
					out <- item
				}
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// First returns new Observable which emit only first item.
func (o *observable) First() Observable {
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
				break
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// Last returns a new Observable which emit only last item.
func (o *observable) Last() Observable {
	f := func(out chan interface{}) {
		var last interface{}
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				last = item
			} else {
				break
			}
		}
		out <- last
		close(out)
	}
	return newColdObservable(f)
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *observable) Distinct(apply Function) Observable {
	f := func(out chan interface{}) {
		keysets := make(map[interface{}]struct{})
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				key := apply(item)
				_, ok := keysets[key]
				if !ok {
					out <- item
				}
				keysets[key] = struct{}{}
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original
// Observable and returns a new Observable.
func (o *observable) DistinctUntilChanged(apply Function) Observable {
	f := func(out chan interface{}) {
		var current interface{}
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				key := apply(item)
				if current != key {
					out <- item
					current = key
				}
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observable) Skip(nth uint) Observable {
	f := func(out chan interface{}) {
		skipCount := 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if skipCount < int(nth) {
					skipCount += 1
					continue
				}
				out <- item
			} else {
				break
			}
		}
		close(out)
	}

	return newColdObservable(f)
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observable) SkipLast(nth uint) Observable {
	f := func(out chan interface{}) {
		buf := make(chan interface{}, nth)
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				select {
				case buf <- item:
				default:
					out <- (<-buf)
					buf <- item
				}
			} else {
				break
			}
		}
		close(buf)
		close(out)
	}
	return newColdObservable(f)
}

// Scan applies Function2 predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o *observable) Scan(apply Function2) Observable {
	f := func(out chan interface{}) {
		var current interface{}
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				tmp := apply(current, item)
				out <- tmp
				current = tmp
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// Compares first items of two sequences and returns true if they are equal and false if
// they are not. Besides, it returns two new sequences - input sequences without compared items.
func popAndCompareFirstItems(
	inputSequence1 []interface{},
	inputSequence2 []interface{}) (bool, []interface{}, []interface{}) {
	if len(inputSequence1) > 0 && len(inputSequence2) > 0 {
		s1, sequence1 := inputSequence1[0], inputSequence1[1:]
		s2, sequence2 := inputSequence2[0], inputSequence2[1:]
		return s1 == s2, sequence1, sequence2
	}
	return true, inputSequence1, inputSequence2
}

// SequenceEqual emits true if an Observable and the input Observable emit the same items,
// in the same order, with the same termination state. Otherwise, it emits false.
func (o *observable) SequenceEqual(obs Observable) Single {
	oChan := make(chan interface{})
	obsChan := make(chan interface{})

	go func() {
		it := obs.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				obsChan <- item
			} else {
				break
			}
		}
		close(obsChan)
	}()

	go func() {
		it := o.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				oChan <- item
			} else {
				break
			}
		}
		close(oChan)
	}()

	f := func(out chan interface{}) {
		var mainSequence []interface{}
		var obsSequence []interface{}
		areCorrect := true
		isMainChannelClosed := false
		isObsChannelClosed := false

	mainLoop:
		for {
			select {
			case item, ok := <-oChan:
				if ok {
					mainSequence = append(mainSequence, item)
					areCorrect, mainSequence, obsSequence = popAndCompareFirstItems(mainSequence, obsSequence)
				} else {
					isMainChannelClosed = true
				}
			case item, ok := <-obsChan:
				if ok {
					obsSequence = append(obsSequence, item)
					areCorrect, mainSequence, obsSequence = popAndCompareFirstItems(mainSequence, obsSequence)
				} else {
					isObsChannelClosed = true
				}
			}

			if !areCorrect || (isMainChannelClosed && isObsChannelClosed) {
				break mainLoop
			}
		}

		out <- areCorrect && len(mainSequence) == 0 && len(obsSequence) == 0
		close(out)
	}

	return newColdSingle(f)
}

func (o *observable) Reduce(apply Function2) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		var acc interface{}
		empty := true
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false
				acc = apply(acc, item)
			} else {
				break
			}
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
	f := func(out chan interface{}) {
		var count int64
		it := o.iterable.Iterator()
		for {
			if _, err := it.Next(); err == nil {
				count++
			} else {
				break
			}
		}
		out <- count
		close(out)
	}
	return newColdSingle(f)
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) FirstOrDefault(defaultValue interface{}) Single {
	f := func(out chan interface{}) {
		first := defaultValue
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				first = item
				break
			} else {
				break
			}
		}
		out <- first
		close(out)
	}
	return newColdSingle(f)
}

// Last returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) LastOrDefault(defaultValue interface{}) Single {
	f := func(out chan interface{}) {
		last := defaultValue
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				last = item
			} else {
				break
			}
		}
		out <- last
		close(out)
	}
	return newColdSingle(f)
}

// TakeWhile emits items emitted by an Observable as long as the
// specified condition is true, then skip the remainder.
func (o *observable) TakeWhile(apply Predicate) Observable {
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if apply(item) {
					out <- item
					continue
				}
				break
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
func (o *observable) SkipWhile(apply Predicate) Observable {
	f := func(out chan interface{}) {
		skip := true
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if !skip {
					out <- item
				} else {
					if !apply(item) {
						out <- item
						skip = false
					}
				}
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// ToList collects all items from an Observable and emit them as a single List.
func (o *observable) ToList() Observable {
	f := func(out chan interface{}) {
		s := make([]interface{}, 0)
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				s = append(s, item)
			} else {
				break
			}
		}
		out <- s
		close(out)
	}
	return newColdObservable(f)
}

// ToMap convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function
func (o *observable) ToMap(keySelector Function) Observable {
	f := func(out chan interface{}) {
		m := make(map[interface{}]interface{})
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				m[keySelector(item)] = item
			} else {
				break
			}
		}
		out <- m
		close(out)
	}
	return newColdObservable(f)
}

// ToMapWithValueSelector convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function and valued by another
// value function
func (o *observable) ToMapWithValueSelector(keySelector Function, valueSelector Function) Observable {
	f := func(out chan interface{}) {
		m := make(map[interface{}]interface{})
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				m[keySelector(item)] = valueSelector(item)
			} else {
				break
			}
		}
		out <- m
		close(out)
	}
	return newColdObservable(f)
}

// ZipFromObservable che emissions of multiple Observables together via a specified function
// and emit single items for each combination based on the results of this function
func (o *observable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		it2 := publisher.Iterator()
	OuterLoop:
		for {
			if item1, err := it.Next(); err == nil {
				for {
					if item2, err := it2.Next(); err == nil {
						out <- zipper(item1, item2)
						continue OuterLoop
					} else {
						break
					}
				}
				break OuterLoop
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// ForEach subscribes to the Observable and receives notifications for each element.
func (o *observable) ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
	doneFunc handlers.DoneFunc, opts ...options.Option) Observer {
	return o.Subscribe(CheckEventHandlers(nextFunc, errFunc, doneFunc), opts...)
}

// Publish returns a ConnectableObservable which waits until its connect method
// is called before it begins emitting items to those Observers that have subscribed to it.
func (o *observable) Publish() ConnectableObservable {
	return newConnectableObservableFromObservable(o)
}

func (o *observable) All(predicate Predicate) Single {
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if !predicate(item) {
					out <- false
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- true
		close(out)
	}
	return newColdSingle(f)
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
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if equal(item) {
					out <- true
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- false
		close(out)
	}
	return newColdSingle(f)
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (o *observable) DefaultIfEmpty(defaultValue interface{}) Observable {
	f := func(out chan interface{}) {
		empty := true
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false
				out <- item
			} else {
				break
			}
		}
		if empty {
			out <- defaultValue
		}
		close(out)
	}
	return newColdObservable(f)
}

// DoOnEach operator allows you to establish a callback that the resulting Observable
// will call each time it emits an item
func (o *observable) DoOnEach(onNotification Consumer) Observable {
	f := func(out chan interface{}) {
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
				onNotification(item)
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservable(f)
}

// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
// at most count times, at a particular frequency.
func (o *observable) Repeat(count int64, frequency Duration) Observable {
	if count != Indefinitely {
		if count < 0 {
			count = 0
		}
	}

	f := func(out chan interface{}) {
		persist := make([]interface{}, 0)
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
				persist = append(persist, item)
			} else {
				break
			}
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
	}
	return newColdObservable(f)
}

// AverageInt calculates the average of numbers emitted by an Observable and emits this average int.
func (o *observable) AverageInt() Single {
	f := func(out chan interface{}) {
		sum := 0
		count := 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// AverageInt8 calculates the average of numbers emitted by an Observable and emits this average int8.
func (o *observable) AverageInt8() Single {
	f := func(out chan interface{}) {
		var sum int8 = 0
		var count int8 = 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int8); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// AverageInt16 calculates the average of numbers emitted by an Observable and emits this average int16.
func (o *observable) AverageInt16() Single {
	f := func(out chan interface{}) {
		var sum int16 = 0
		var count int16 = 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int16); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// AverageInt32 calculates the average of numbers emitted by an Observable and emits this average int32.
func (o *observable) AverageInt32() Single {
	f := func(out chan interface{}) {
		var sum int32 = 0
		var count int32 = 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int32); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// AverageInt64 calculates the average of numbers emitted by an Observable and emits this average int64.
func (o *observable) AverageInt64() Single {
	f := func(out chan interface{}) {
		var sum int64 = 0
		var count int64 = 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int64); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits this average float32.
func (o *observable) AverageFloat32() Single {
	f := func(out chan interface{}) {
		var sum float32 = 0
		var count float32 = 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(float32); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// AverageFloat64 calculates the average of numbers emitted by an Observable and emits this average float64.
func (o *observable) AverageFloat64() Single {
	f := func(out chan interface{}) {
		var sum float64 = 0
		var count float64 = 0
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(float64); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
	return newColdSingle(f)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (o *observable) Max(comparator Comparator) OptionalSingle {
	out := make(chan optional.Optional)
	go func() {
		empty := true
		var max interface{} = nil
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false

				if max == nil {
					max = item
				} else {
					if comparator(max, item) == Smaller {
						max = item
					}
				}
			} else {
				break
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
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false

				if min == nil {
					min = item
				} else {
					if comparator(min, item) == Greater {
						min = item
					}
				}
			} else {
				break
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
	f := func(out chan interface{}) {
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
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
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
			} else {
				break
			}
		}
		if iCount != 0 {
			out <- buffer[:iCount]
		}

		close(out)
	}
	return newColdObservable(f)
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (o *observable) BufferWithTime(timespan, timeshift Duration) Observable {
	f := func(out chan interface{}) {
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
			it := o.iterable.Iterator()
			for {
				if item, err := it.Next(); err == nil {
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
				} else {
					break
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

	}
	return newColdObservable(f)
}

// BufferWithTimeOrCount returns an Observable that emits buffers of items it collects
// from the source Observable. The resulting Observable emits connected,
// non-overlapping buffers, each of a fixed duration specified by the timespan argument
// or a maximum size specified by the count argument (whichever is reached first).
// When the source Observable completes or encounters an error, the resulting Observable
// emits the current buffer and propagates the notification from the source Observable.
func (o *observable) BufferWithTimeOrCount(timespan Duration, count int) Observable {
	f := func(out chan interface{}) {
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
			it := o.iterable.Iterator()
			for {
				if item, err := it.Next(); err == nil {
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
				} else {
					break
				}
			}
			errCh <- nil
		}()

	}
	return newColdObservable(f)
}

// SumInt64 calculates the average of integers emitted by an Observable and emits an int64.
func (o *observable) SumInt64() Single {
	f := func(out chan interface{}) {
		var sum int64
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case int:
					sum = sum + int64(item)
				case int8:
					sum = sum + int64(item)
				case int16:
					sum = sum + int64(item)
				case int32:
					sum = sum + int64(item)
				case int64:
					sum = sum + item
				default:
					out <- errors.New(errors.IllegalInputError,
						fmt.Sprintf("expected type: int, int8, int16, int32 or int64, got %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- sum
		close(out)
	}
	return newColdSingle(f)
}

// SumFloat32 calculates the average of float32 emitted by an Observable and emits a float32.
func (o *observable) SumFloat32() Single {
	f := func(out chan interface{}) {
		var sum float32
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case int:
					sum = sum + float32(item)
				case int8:
					sum = sum + float32(item)
				case int16:
					sum = sum + float32(item)
				case int32:
					sum = sum + float32(item)
				case int64:
					sum = sum + float32(item)
				case float32:
					sum = sum + item
				default:
					out <- errors.New(errors.IllegalInputError,
						fmt.Sprintf("expected type: float32, int, int8, int16, int32 or int64, got %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- sum
		close(out)
	}
	return newColdSingle(f)
}

// SumFloat64 calculates the average of float64 emitted by an Observable and emits a float64.
func (o *observable) SumFloat64() Single {
	f := func(out chan interface{}) {
		var sum float64
		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case int:
					sum = sum + float64(item)
				case int8:
					sum = sum + float64(item)
				case int16:
					sum = sum + float64(item)
				case int32:
					sum = sum + float64(item)
				case int64:
					sum = sum + float64(item)
				case float32:
					sum = sum + float64(item)
				case float64:
					sum = sum + item
				default:
					out <- errors.New(errors.IllegalInputError,
						fmt.Sprintf("expected type: float32, float64, int, int8, int16, int32 or int64, got %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- sum
		close(out)
	}
	return newColdSingle(f)
}

// StartWithItems returns an Observable that emits the specified items before it begins to emit items emitted
// by the source Observable.
func (o *observable) StartWithItems(items ...interface{}) Observable {
	f := func(out chan interface{}) {
		for _, item := range items {
			out <- item
		}

		it := o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
	return newColdObservable(f)
}

// StartWithIterable returns an Observable that emits the items in a specified Iterable before it begins to
// emit items emitted by the source Observable.
func (o *observable) StartWithIterable(iterable Iterable) Observable {
	f := func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		it = o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
	return newColdObservable(f)
}

// StartWithObservable returns an Observable that emits the items in a specified Observable before it begins to
// emit items emitted by the source Observable.
func (o *observable) StartWithObservable(obs Observable) Observable {
	f := func(out chan interface{}) {
		it := obs.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		it = o.iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
	return newColdObservable(f)
}
