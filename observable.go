package rxgo

import (
	"container/ring"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
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
	Debounce(debounceTimeValue time.Duration) Observable
	DefaultIfEmpty(defaultValue interface{}) Observable
	Distinct(apply Function) Observable
	DistinctUntilChanged(apply Function) Observable
	DoOnEach(onNotification Consumer) Observable
	ElementAt(index uint) Single
	Filter(apply Predicate) Observable
	First() Observable
	FirstOrDefault(defaultValue interface{}) Single
	FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable
	ForEach(nextFunc NextFunc, errFunc ErrFunc,
		doneFunc DoneFunc, opts ...Option) Observer
	IgnoreElements() Observable
	Last() Observable
	LastOrDefault(defaultValue interface{}) Single
	Map(apply Function, opts ...Option) Observable
	Marshal(Marshaler, ...Option) Observable
	Max(comparator Comparator) OptionalSingle
	Min(comparator Comparator) OptionalSingle
	OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable
	OnErrorReturn(resumeFunc ErrorFunction) Observable
	OnErrorReturnItem(item interface{}) Observable
	Publish() ConnectableObservable
	Reduce(apply Function2) OptionalSingle
	Repeat(count int64, frequency Duration) Observable
	Sample(obs Observable) Observable
	Scan(apply Function2) Observable
	SequenceEqual(obs Observable) Single
	Send(chan<- interface{})
	Skip(nth uint) Observable
	SkipLast(nth uint) Observable
	SkipWhile(apply Predicate) Observable
	StartWithItems(item interface{}, items ...interface{}) Observable
	StartWithIterable(iterable Iterable) Observable
	StartWithObservable(observable Observable) Observable
	// TODO Pass a context to cancel the subscription
	Subscribe(handler EventHandler, opts ...Option) Observer
	SumFloat32() Single
	SumFloat64() Single
	SumInt64() Single
	Take(nth uint) Observable
	TakeLast(nth uint) Observable
	TakeUntil(apply Predicate) Observable
	TakeWhile(apply Predicate) Observable
	Timeout(ctx context.Context) Observable
	ToChannel(opts ...Option) Channel
	ToMap(keySelector Function) Single
	ToMapWithValueSelector(keySelector, valueSelector Function) Single
	ToSlice() Single
	Unmarshal(Unmarshaler, func() interface{}, ...Option) Observable
	ZipFromObservable(publisher Observable, zipper Function2) Observable
	getCustomErrorStrategy() func(Observable, Observer, error) error
	getNextStrategy() func(Observer, interface{}) error
}

// observable is a structure handling a channel of interface{} and implementing Observable
type observable struct {
	// subscribeStrategy represents the observable strategy once an observer subscribes to the observable
	subscribeStrategy func(o *observable, ob Observer)
	// customErrorStrategy represents a custom strategy if the observable encounters an error
	customErrorStrategy func(Observable, Observer, error) error
	// nextStrategy represents the stategy for the next item
	nextStrategy func(Observer, interface{}) error
	// errorOnSubscription defines an error to be sent to the observer once it subscribes to the observable
	errorOnSubscription error
	// coldIterable is used while pulling data
	coldIterable Iterable
	// hotItemChannel is used while data are pushed from the observable
	hotItemChannel <-chan interface{}
	// hotObservers list the observer registered to the hot observable with none bp mode.
	hotObservers []Observer
	// hotObservers protects hotObservers from concurrent accesses
	hotObserversMutex sync.Mutex
	// hotSubscribers list the subscribers registered to the hot observable with none buffer bp mode.
	hotSubscribers []chan<- interface{}
}

func onErrorReturn(f ErrorFunction) func(Observable, Observer, error) error {
	return func(_ Observable, observer Observer, err error) error {
		return observer.OnNext(f(err))
	}
}

func onErrorResumeNext(f ErrorToObservableFunction) func(Observable, Observer, error) error {
	return func(observable Observable, observer Observer, err error) error {
		return iterate(f(err), observer)
	}
}

func onErrorReturnItem(item interface{}) func(Observable, Observer, error) error {
	return func(observable Observable, observer Observer, _ error) error {
		return observer.OnNext(item)
	}
}

func onNextIgnore() func(Observer, interface{}) error {
	return func(observer Observer, item interface{}) error {
		return nil
	}
}

func onNext() func(Observer, interface{}) error {
	return func(observer Observer, item interface{}) error {
		return observer.OnNext(item)
	}
}

func iterate(observable Observable, observer Observer) error {
	it := observable.Iterator(context.Background())
	for {
		if item, err := it.Next(context.Background()); err == nil {
			switch item := item.(type) {
			case error:
				errorStrategy := observable.getCustomErrorStrategy()
				if errorStrategy == nil {
					err := observer.OnError(item)
					if err != nil {
						return err
					}
					return item
				}
				err := errorStrategy(observable, observer, item)
				if err != nil {
					return err
				}
			default:
				err := observable.getNextStrategy()(observer, item)
				if err != nil {
					return err
				}
			}
		} else {
			break
		}
	}
	return nil
}

func coldSubscribe() func(o *observable, ob Observer) {
	return func(o *observable, ob Observer) {
		go func() {
			e := iterate(o, ob)
			if e == nil {
				err := ob.OnDone()
				if err != nil {
					panic(errors.Wrap(err, "error while sending done signal from observable"))
				}
			}
		}()
	}
}

func hotSubscribeStrategyNoneBackPressure() func(o *observable, ob Observer) {
	return func(o *observable, ob Observer) {
		o.hotObserversMutex.Lock()
		o.hotObservers = append(o.hotObservers, ob)
		o.hotObserversMutex.Unlock()
	}
}

func hotSubscribeStrategyBufferBackPressure(bpBuffer int) func(o *observable, ob Observer) {
	return func(o *observable, ob Observer) {
		o.hotObserversMutex.Lock()
		ch := make(chan interface{}, bpBuffer)
		go func() {
			for item := range ch {
				ob.Handle(item)
			}
		}()
		o.hotSubscribers = append(o.hotSubscribers, ch)
		o.hotObserversMutex.Unlock()
	}
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

func (o *observable) Iterator(ctx context.Context) Iterator {
	return o.coldIterable.Iterator(ctx)
}

func (o *observable) All(predicate Predicate) Single {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits this average float32.
func (o *observable) AverageFloat32() Single {
	f := func(out chan interface{}) {
		var sum float32
		var count float32
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(float32); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: float32, got: %t", item))
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

// AverageInt calculates the average of numbers emitted by an Observable and emits this average int.
func (o *observable) AverageInt() Single {
	f := func(out chan interface{}) {
		sum := 0
		count := 0
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(int); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: int, got: %t", item))
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
		var sum int8
		var count int8
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(int8); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: int8, got: %t", item))
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
		var sum float64
		var count float64
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(float64); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: float64, got: %t", item))
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
		var sum int16
		var count int16
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(int16); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: int16, got: %t", item))
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
		var sum int32
		var count int32
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(int32); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: int32, got: %t", item))
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
		var sum int64
		var count int64
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if v, ok := item.(int64); ok {
					sum += v
					count++
				} else {
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: int64, got: %t", item))
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

// BufferWithCount returns an Observable that emits buffers of items it collects
// from the source Observable.
// The resulting Observable emits buffers every skip items, each containing a slice of count items.
// When the source Observable completes or encounters an error,
// the resulting Observable emits the current buffer and propagates
// the notification from the source Observable.
func (o *observable) BufferWithCount(count, skip int) Observable {
	f := func(out chan interface{}) {
		if count <= 0 {
			out <- errors.Wrap(&IllegalInputError{}, "count must be positive")
			close(out)
			return
		}

		if skip <= 0 {
			out <- errors.Wrap(&IllegalInputError{}, "skip must be positive")
			close(out)
			return
		}

		buffer := make([]interface{}, count)
		iCount := 0
		iSkip := 0
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
						buffer = make([]interface{}, count)
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
	return newColdObservableFromFunction(f)
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (o *observable) BufferWithTime(timespan, timeshift Duration) Observable {
	f := func(out chan interface{}) {
		if timespan == nil || timespan.duration() == 0 {
			out <- errors.Wrap(&IllegalInputError{}, "timespan must no be nil")
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
			it := o.coldIterable.Iterator(context.Background())
			for {
				if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
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
			out <- errors.Wrap(&IllegalInputError{}, "timespan must not be nil")
			close(out)
			return
		}

		if count <= 0 {
			out <- errors.Wrap(&IllegalInputError{}, "count must be positive")
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
			it := o.coldIterable.Iterator(context.Background())
			for {
				if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

// Contains returns an Observable that emits a Boolean that indicates whether
// the source Observable emitted an item (the comparison is made against a predicate).
func (o *observable) Contains(equal Predicate) Single {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

func (o *observable) Count() Single {
	f := func(out chan interface{}) {
		var count int64
		it := o.coldIterable.Iterator(context.Background())
		for {
			if _, err := it.Next(context.Background()); err == nil {
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

//Debounce only emit an item from an Observable if a particular timespan has passed
func (o *observable) Debounce(debounceTimeValue time.Duration) Observable {
	f := func(out chan interface{}) {
		time.Sleep(debounceTimeValue * time.Second)
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (o *observable) DefaultIfEmpty(defaultValue interface{}) Observable {
	f := func(out chan interface{}) {
		empty := true
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *observable) Distinct(apply Function) Observable {
	f := func(out chan interface{}) {
		keysets := make(map[interface{}]struct{})
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original
// Observable and returns a new Observable.
func (o *observable) DistinctUntilChanged(apply Function) Observable {
	f := func(out chan interface{}) {
		var current interface{}
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

// DoOnEach operator allows you to establish a callback that the resulting Observable
// will call each time it emits an item
func (o *observable) DoOnEach(onNotification Consumer) Observable {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
				onNotification(item)
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

func (o *observable) ElementAt(index uint) Single {
	f := func(out chan interface{}) {
		takeCount := 0

		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if takeCount == int(index) {
					out <- item
					close(out)
					return
				}
				takeCount++
			} else {
				break
			}
		}
		out <- &IndexOutOfBoundError{}
		out <- &IndexOutOfBoundError{}
		close(out)
	}
	return newColdSingle(f)
}

// Filter filters items in the original Observable and returns
// a new Observable with the filtered items.
func (o *observable) Filter(apply Predicate) Observable {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if apply(item) {
					out <- item
				}
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) FirstOrDefault(defaultValue interface{}) Single {
	f := func(out chan interface{}) {
		first := defaultValue
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

// First returns new Observable which emit only first item.
func (o *observable) First() Observable {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
				break
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// ForEach subscribes to the Observable and receives notifications for each element.
func (o *observable) ForEach(nextFunc NextFunc, errFunc ErrFunc,
	doneFunc DoneFunc, opts ...Option) Observer {
	return o.Subscribe(NewObserver(nextFunc, errFunc, doneFunc), opts...)
}

// IgnoreElements ignores all items emitted by the source ObservableSource and only calls onComplete
// or onError.
func (o *observable) IgnoreElements() Observable {
	return &observable{
		subscribeStrategy:   coldSubscribe(),
		coldIterable:        o.coldIterable,
		errorOnSubscription: o.errorOnSubscription,
		nextStrategy:        onNextIgnore(),
	}
}

// Last returns a new Observable which emit only last item.
func (o *observable) Last() Observable {
	f := func(out chan interface{}) {
		var last interface{}
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				last = item
			} else {
				break
			}
		}
		out <- last
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// Last returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) LastOrDefault(defaultValue interface{}) Single {
	f := func(out chan interface{}) {
		last := defaultValue
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

// Map maps a Function predicate to each item in Observable and
// returns a new Observable with applied items.
func (o *observable) Map(apply Function, opts ...Option) Observable {
	f := func(out chan interface{}) {
		options := ParseOptions(opts...)

		workerPoolCapacity := options.NewWorkerPool()
		wp := options.WorkerPool()

		var ctx context.Context
		if c := options.Context(); c != nil {
			ctx = c
		} else {
			c, cancel := context.WithCancel(context.Background())
			ctx = c
			defer cancel()
		}

		if workerPoolCapacity != 0 || wp != nil {
			var workerPool *workerPool
			wg := sync.WaitGroup{}

			if wp != nil {
				workerPool = wp
			} else {
				pool := newWorkerPool(ctx, workerPoolCapacity)
				workerPool = &pool
			}

			output := make(chan interface{}, workerPoolCapacity)
			it := o.Iterator(context.Background())
			for {
				if item, err := it.Next(context.Background()); err == nil {
					workerPool.sendTask(item, apply, output, &wg)
				} else {
					break
				}
			}

			workerPool.wait(func(i interface{}) {
				out <- i
			}, output, &wg)

			close(out)
		} else {
			it := o.Iterator(context.Background())
			for {
				if item, err := it.Next(context.Background()); err == nil {
					out <- apply(item)
				} else {
					break
				}
			}
			close(out)
		}
	}

	return newColdObservableFromFunction(f)
}

func (o *observable) Marshal(marshaler Marshaler, opts ...Option) Observable {
	return o.Map(func(i interface{}) interface{} {
		b, err := marshaler(i)
		if err != nil {
			return err
		}
		return b
	}, opts...)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (o *observable) Max(comparator Comparator) OptionalSingle {
	out := make(chan Optional)
	go func() {
		empty := true
		var max interface{}
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
			out <- EmptyOptional()
		} else {
			out <- Of(max)
		}
		close(out)
	}()
	return &optionalSingle{itemChannel: out}
}

// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
func (o *observable) Min(comparator Comparator) OptionalSingle {
	out := make(chan Optional)
	go func() {
		empty := true
		var min interface{}
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
			out <- EmptyOptional()
		} else {
			out <- Of(min)
		}
		close(out)
	}()
	return &optionalSingle{itemChannel: out}
}

func (o *observable) Send(ch chan<- interface{}) {
	go func() {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				ch <- item
			} else {
				close(ch)
				return
			}
		}
	}()
}

// OnErrorResumeNext instructs an Observable to pass control to another Observable rather than invoking
// onError if it encounters an error.
func (o *observable) OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable {
	o.customErrorStrategy = onErrorResumeNext(resumeSequence)
	return o
}

// OnErrorReturn instructs an Observable to emit an item (returned by a specified function)
// rather than invoking onError if it encounters an error.
func (o *observable) OnErrorReturn(resumeFunc ErrorFunction) Observable {
	o.customErrorStrategy = onErrorReturn(resumeFunc)
	return o
}

// OnErrorReturnItem instructs on observale to emit an item if it encounters an error.
func (o *observable) OnErrorReturnItem(item interface{}) Observable {
	o.customErrorStrategy = onErrorReturnItem(item)
	return o
}

// Publish returns a ConnectableObservable which waits until its connect method
// is called before it begins emitting items to those Observers that have subscribed to it.
func (o *observable) Publish() ConnectableObservable {
	return newConnectableObservableFromObservable(o)
}

func (o *observable) Reduce(apply Function2) OptionalSingle {
	out := make(chan Optional)
	go func() {
		var acc interface{}
		empty := true
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				empty = false
				acc = apply(acc, item)
			} else {
				break
			}
		}
		if empty {
			out <- EmptyOptional()
		} else {
			out <- Of(acc)
		}
		close(out)
	}()
	return NewOptionalSingleFromChannel(out)
}

// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
// at most count times, at a particular frequency.
func (o *observable) Repeat(count int64, frequency Duration) Observable {
	if count != Infinite {
		if count < 0 {
			count = 0
		}
	}

	f := func(out chan interface{}) {
		persist := make([]interface{}, 0)
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
				persist = append(persist, item)
			} else {
				break
			}
		}
		for {
			if count != Infinite {
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
	return newColdObservableFromFunction(f)
}

// Sample returns an Observable that emits the most recent items emitted by the source
// ObservableSource whenever the input Observable emits an item.
func (o *observable) Sample(obs Observable) Observable {
	f := func(out chan interface{}) {
		mainChan := make(chan interface{})
		obsChan := make(chan interface{})
		var lastEmittedItem interface{}
		isItemWaitingToBeEmitted := false
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			it := o.coldIterable.Iterator(ctx)
			for {
				if item, err := it.Next(ctx); err == nil {
					mainChan <- item
				} else {
					break
				}
			}
			close(mainChan)
		}()

		go func() {
			it := obs.Iterator(ctx)
			for {
				if item, err := it.Next(ctx); err == nil {
					obsChan <- item
				} else {
					break
				}
			}
			close(obsChan)
		}()

		defer cancel()
		defer close(out)

		for {
			select {
			case item, ok := <-mainChan:
				if ok {
					lastEmittedItem = item
					isItemWaitingToBeEmitted = true
				} else {
					return
				}
			case _, ok := <-obsChan:
				if ok {
					if isItemWaitingToBeEmitted {
						out <- lastEmittedItem
						isItemWaitingToBeEmitted = false
					}
				} else {
					return
				}
			}
		}
	}

	return newColdObservableFromFunction(f)
}

// Scan applies Function2 predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o *observable) Scan(apply Function2) Observable {
	f := func(out chan interface{}) {
		var current interface{}
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				tmp := apply(current, item)
				out <- tmp
				current = tmp
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// SequenceEqual emits true if an Observable and the input Observable emit the same items,
// in the same order, with the same termination state. Otherwise, it emits false.
func (o *observable) SequenceEqual(obs Observable) Single {
	oChan := make(chan interface{})
	obsChan := make(chan interface{})

	go func() {
		it := obs.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				obsChan <- item
			} else {
				break
			}
		}
		close(obsChan)
	}()

	go func() {
		it := o.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observable) Skip(nth uint) Observable {
	f := func(out chan interface{}) {
		skipCount := 0
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if skipCount < int(nth) {
					skipCount++
					continue
				}
				out <- item
			} else {
				break
			}
		}
		close(out)
	}

	return newColdObservableFromFunction(f)
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
func (o *observable) SkipLast(nth uint) Observable {
	f := func(out chan interface{}) {
		buf := make(chan interface{}, nth)
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

// StartWithItems returns an Observable that emits the specified items before it begins to emit items emitted
// by the source Observable.
func (o *observable) StartWithItems(item interface{}, items ...interface{}) Observable {
	f := func(out chan interface{}) {
		out <- item
		for _, item := range items {
			out <- item
		}

		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
	return newColdObservableFromFunction(f)
}

// StartWithIterable returns an Observable that emits the items in a specified Iterable before it begins to
// emit items emitted by the source Observable.
func (o *observable) StartWithIterable(iterable Iterable) Observable {
	f := func(out chan interface{}) {
		it := iterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
			} else {
				break
			}
		}

		it = o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
	return newColdObservableFromFunction(f)
}

// StartWithObservable returns an Observable that emits the items in a specified Observable before it begins to
// emit items emitted by the source Observable.
func (o *observable) StartWithObservable(obs Observable) Observable {
	f := func(out chan interface{}) {
		it := obs.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
			} else {
				break
			}
		}

		it = o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
	return newColdObservableFromFunction(f)
}

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
func (o *observable) SkipWhile(apply Predicate) Observable {
	f := func(out chan interface{}) {
		skip := true
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

// Subscribe subscribes an EventHandler and returns a Subscription channel.
func (o *observable) Subscribe(handler EventHandler, opts ...Option) Observer {
	ob := NewObserver(handler)

	if o.errorOnSubscription != nil {
		err := ob.OnError(o.errorOnSubscription)
		if err != nil {
			panic(errors.Wrap(err, "error while sending error item from observable"))
		}
		return ob
	}

	o.subscribeStrategy(o, ob)
	return ob
}

// SumInt64 calculates the average of integers emitted by an Observable and emits an int64.
func (o *observable) SumInt64() Single {
	f := func(out chan interface{}) {
		var sum int64
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				switch item := item.(type) {
				case int:
					sum += int64(item)
				case int8:
					sum += int64(item)
				case int16:
					sum += int64(item)
				case int32:
					sum += int64(item)
				case int64:
					sum += item
				default:
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: (int|int8|int16|int32|int64), got: %t", item))
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
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				switch item := item.(type) {
				case int:
					sum += float32(item)
				case int8:
					sum += float32(item)
				case int16:
					sum += float32(item)
				case int32:
					sum += float32(item)
				case int64:
					sum += float32(item)
				case float32:
					sum += item
				default:
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: (float32|int|int8|int16|int32|int64), got: %t", item))
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
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				switch item := item.(type) {
				case int:
					sum += float64(item)
				case int8:
					sum += float64(item)
				case int16:
					sum += float64(item)
				case int32:
					sum += float64(item)
				case int64:
					sum += float64(item)
				case float32:
					sum += float64(item)
				case float64:
					sum += item
				default:
					out <- errors.Wrap(&IllegalInputError{},
						fmt.Sprintf("expected type: (float32|float64|int|int8|int16|int32|int64), got: %t", item))
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

// Take takes first n items in the original Obserable and returns
// a new Observable with the taken items.
func (o *observable) Take(nth uint) Observable {
	f := func(out chan interface{}) {
		takeCount := 0
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				if takeCount < int(nth) {
					takeCount++
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
	return newColdObservableFromFunction(f)
}

// TakeLast takes last n items in the original Observable and returns
// a new Observable with the taken items.
func (o *observable) TakeLast(nth uint) Observable {
	f := func(out chan interface{}) {
		n := int(nth)
		r := ring.New(n)
		count := 0
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				count++
				r.Value = item
				r = r.Next()
			} else {
				break
			}
		}
		if count < n {
			remaining := n - count
			if remaining <= count {
				r = r.Move(n - count)
			} else {
				r = r.Move(-count)
			}
			n = count
		}
		for i := 0; i < n; i++ {
			out <- r.Value
			r = r.Next()
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// TakeUntil returns an Observable that emits items emitted by the source Observable,
// checks the specified predicate for each item, and then completes when the condition is satisfied.
func (o *observable) TakeUntil(apply Predicate) Observable {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				out <- item
				if apply(item) {
					break
				}
			} else {
				break
			}
		}
		close(out)
	}
	return newColdObservableFromFunction(f)
}

// TakeWhile returns an Observable that emits items emitted by the source ObservableSource so long as each
// item satisfied a specified condition, and then completes as soon as this condition is not satisfied.
func (o *observable) TakeWhile(apply Predicate) Observable {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

func (o *observable) Timeout(ctx context.Context) Observable {
	f := func(out chan interface{}) {
		it := o.Iterator(context.Background())
		go func() {
			for {
				if item, err := it.Next(ctx); err == nil {
					out <- item
				} else {
					out <- err
					break
				}
			}
			close(out)
		}()
	}

	return newColdObservableFromFunction(f)
}

// ToChannel collects all items from an Observable and emit them in a channel
func (o *observable) ToChannel(opts ...Option) Channel {
	options := ParseOptions(opts...)
	var ch chan interface{}
	if options.Buffer() != 0 {
		ch = make(chan interface{}, options.Buffer())
	} else {
		ch = make(chan interface{})
	}

	go func() {
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				ch <- item
			} else {
				break
			}
		}
		close(ch)
	}()

	return ch
}

// ToSlice collects all items from an Observable and emit them as a single slice.
func (o *observable) ToSlice() Single {
	f := func(out chan interface{}) {
		s := make([]interface{}, 0)
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				s = append(s, item)
			} else {
				break
			}
		}
		out <- s
		close(out)
	}
	return newColdSingle(f)
}

// ToMap convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function
func (o *observable) ToMap(keySelector Function) Single {
	f := func(out chan interface{}) {
		m := make(map[interface{}]interface{})
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				m[keySelector(item)] = item
			} else {
				break
			}
		}
		out <- m
		close(out)
	}
	return newColdSingle(f)
}

// ToMapWithValueSelector convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function and valued by another
// value function
func (o *observable) ToMapWithValueSelector(keySelector, valueSelector Function) Single {
	f := func(out chan interface{}) {
		m := make(map[interface{}]interface{})
		it := o.coldIterable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				m[keySelector(item)] = valueSelector(item)
			} else {
				break
			}
		}
		out <- m
		close(out)
	}
	return newColdSingle(f)
}

func (o *observable) Unmarshal(unmarshaler Unmarshaler, factory func() interface{}, opts ...Option) Observable {
	return o.Map(func(i interface{}) interface{} {
		v := factory()
		err := unmarshaler(i.([]byte), v)
		if err != nil {
			return err
		}
		return v
	}, opts...)
}

// ZipFromObservable che emissions of multiple Observables together via a specified function
// and emit single items for each combination based on the results of this function
func (o *observable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	f := func(out chan interface{}) {
		it := o.coldIterable.Iterator(context.Background())
		it2 := publisher.Iterator(context.Background())
	OuterLoop:
		for {
			if item1, err := it.Next(context.Background()); err == nil {
				for {
					if item2, err := it2.Next(context.Background()); err == nil {
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
	return newColdObservableFromFunction(f)
}

func (o *observable) getCustomErrorStrategy() func(Observable, Observer, error) error {
	return o.customErrorStrategy
}

func (o *observable) getNextStrategy() func(Observer, interface{}) error {
	return o.nextStrategy
}
