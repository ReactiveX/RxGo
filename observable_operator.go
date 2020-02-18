package rxgo

import (
	"container/ring"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/emirpasic/gods/trees/binaryheap"
)

// All determine whether all items emitted by an Observable meet some criteria.
func (o *ObservableImpl) All(predicate Predicate, opts ...Option) Single {
	all := true
	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if !predicate(item.V) {
			dst <- Of(false)
			all = false
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if all {
			dst <- Of(true)
		}
	}, opts...)
}

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits the average float32.
func (o *ObservableImpl) AverageFloat32(opts ...Option) Single {
	var sum float32
	var count float32

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(float32); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: float32, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// AverageFloat64 calculates the average of numbers emitted by an Observable and emits the average float64.
func (o *ObservableImpl) AverageFloat64(opts ...Option) Single {
	var sum float64
	var count float64

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(float64); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: float64, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// AverageInt calculates the average of numbers emitted by an Observable and emits the average int.
func (o *ObservableImpl) AverageInt(opts ...Option) Single {
	var sum int
	var count int

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(int); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: int, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// AverageInt8 calculates the average of numbers emitted by an Observable and emits the≤ average int8.
func (o *ObservableImpl) AverageInt8(opts ...Option) Single {
	var sum int8
	var count int8

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(int8); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: int8, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// AverageInt16 calculates the average of numbers emitted by an Observable and emits the average int16.
func (o *ObservableImpl) AverageInt16(opts ...Option) Single {
	var sum int16
	var count int16

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(int16); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: int16, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// AverageInt32 calculates the average of numbers emitted by an Observable and emits the average int32.
func (o *ObservableImpl) AverageInt32(opts ...Option) Single {
	var sum int32
	var count int32

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(int32); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: int32, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// AverageInt64 calculates the average of numbers emitted by an Observable and emits this average int64.
func (o *ObservableImpl) AverageInt64(opts ...Option) Single {
	var sum int64
	var count int64

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.V.(int64); ok {
			sum += v
			count++
		} else {
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: int64, got: %t", item)})
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if count == 0 {
			dst <- Of(0)
		} else {
			dst <- Of(sum / count)
		}
	}, opts...)
}

// BufferWithCount returns an Observable that emits buffers of items it collects
// from the source Observable.
// The resulting Observable emits buffers every skip items, each containing a slice of count items.
// When the source Observable completes or encounters an error,
// the resulting Observable emits the current buffer and propagates
// the notification from the source Observable.
func (o *ObservableImpl) BufferWithCount(count, skip int, opts ...Option) Observable {
	if count <= 0 {
		return newObservableFromError(IllegalInputError{error: "count must be positive"})
	}
	if skip <= 0 {
		return newObservableFromError(IllegalInputError{error: "skip must be positive"})
	}

	buffer := make([]interface{}, count)
	iCount := 0
	iSkip := 0

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if iCount >= count {
			// Skip
			iSkip++
		} else {
			// Add to buffer
			buffer[iCount] = item.V
			iCount++
			iSkip++
		}
		if iSkip == skip {
			// Send current buffer
			dst <- Of(buffer)
			buffer = make([]interface{}, count)
			iCount = 0
			iSkip = 0
		}
	}, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if iCount != 0 {
			dst <- Of(buffer[:iCount])
		}
		dst <- item
		iCount = 0
		operator.stop()
	}, func(_ context.Context, dst chan<- Item) {
		if iCount != 0 {
			dst <- Of(buffer[:iCount])
		}
	}, opts...)
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (o *ObservableImpl) BufferWithTime(timespan, timeshift Duration, opts ...Option) Observable {
	if timespan == nil || timespan.duration() == 0 {
		return newObservableFromError(IllegalInputError{error: "timespan must no be nil"})
	}
	if timeshift == nil {
		timeshift = WithDuration(0)
	}

	var mux sync.Mutex
	var listenMutex sync.Mutex
	buffer := make([]interface{}, 0)
	stop := false
	listen := true

	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	stopped := false

	// First goroutine in charge to check the timespan
	go func() {
		for {
			time.Sleep(timespan.duration())
			mux.Lock()
			if !stop {
				next <- Of(buffer)
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

	go func() {
		observe := o.Observe()
	loop:
		for !stopped {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.Error() {
					mux.Lock()
					if len(buffer) > 0 {
						next <- Of(buffer)
					}
					next <- i
					close(next)
					stop = true
					mux.Unlock()
					return
				}
				listenMutex.Lock()
				l := listen
				listenMutex.Unlock()

				mux.Lock()
				if l {
					buffer = append(buffer, i.V)
				}
				mux.Unlock()
			}
		}
		mux.Lock()
		if len(buffer) > 0 {
			next <- Of(buffer)
		}
		close(next)
		stop = true
		mux.Unlock()
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// BufferWithTimeOrCount returns an Observable that emits buffers of items it collects from the source
// Observable either from a given count or at a given time interval.
func (o *ObservableImpl) BufferWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable {
	if timespan == nil || timespan.duration() == 0 {
		return newObservableFromError(IllegalInputError{error: "timespan must no be nil"})
	}
	if count <= 0 {
		return newObservableFromError(IllegalInputError{error: "count must be positive"})
	}

	sendCh := make(chan []interface{})
	errCh := make(chan error)
	buffer := make([]interface{}, 0)
	var bufferMutex sync.Mutex
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	// First sender goroutine
	go func() {
		for {
			select {
			case currentBuffer := <-sendCh:
				next <- Of(currentBuffer)
			case error := <-errCh:
				bufferMutex.Lock()
				if len(buffer) > 0 {
					next <- Of(buffer)
				}
				bufferMutex.Unlock()
				if error != nil {
					next <- Error(error)
				}
				close(next)
				return
			case <-time.After(timespan.duration()):
				bufferMutex.Lock()
				b := make([]interface{}, len(buffer))
				copy(b, buffer)
				buffer = make([]interface{}, 0)
				bufferMutex.Unlock()
				next <- Of(b)
			}
		}
	}()

	go func() {
		observe := o.Observe()
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.Error() {
					errCh <- i.E
					break loop
				}
				// TODO Improve implementation without mutex (sending data over channel)
				bufferMutex.Lock()
				buffer = append(buffer, i.V)
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
		close(sendCh)
		close(errCh)
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Contains determines whether an Observable emits a particular item or not.
func (o *ObservableImpl) Contains(equal Predicate, opts ...Option) Single {
	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if equal(item.V) {
			dst <- Of(true)
			operator.stop()
			return
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		dst <- Of(false)
	}, opts...)
}

// Count counts the number of items emitted by the source Observable and emit only this value.
func (o *ObservableImpl) Count(opts ...Option) Single {
	var count int64
	return newSingleFromOperator(o, func(_ context.Context, _ Item, dst chan<- Item, _ operatorOptions) {
		count++
	}, func(_ context.Context, _ Item, dst chan<- Item, operator operatorOptions) {
		count++
		dst <- Of(count)
		operator.stop()
	}, defaultEndFuncOperator, opts...)
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (o *ObservableImpl) DefaultIfEmpty(defaultValue interface{}, opts ...Option) Observable {
	empty := true

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		empty = false
		dst <- item
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if empty {
			dst <- Of(defaultValue)
		}
	}, opts...)
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *ObservableImpl) Distinct(apply Func, opts ...Option) Observable {
	keyset := make(map[interface{}]interface{})

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		key, err := apply(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}
		_, ok := keyset[key]
		if !ok {
			dst <- item
		}
		keyset[key] = nil
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original
// Observable and returns a new Observable.
func (o *ObservableImpl) DistinctUntilChanged(apply Func, opts ...Option) Observable {
	var current interface{}

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		key, err := apply(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}
		if current != key {
			dst <- item
			current = key
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// DoOnCompleted registers a callback action that will be called once the Observable terminates.
func (o *ObservableImpl) DoOnCompleted(completedFunc CompletedFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		defer completedFunc()
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-src:
				if !ok {
					return
				}
				if i.Error() {
					return
				}
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe())
	return dispose
}

// DoOnError registers a callback action that will be called if the Observable terminates abnormally.
func (o *ObservableImpl) DoOnError(errFunc ErrFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-src:
				if !ok {
					return
				}
				if i.Error() {
					errFunc(i.E)
					return
				}
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe())
	return dispose
}

// DoOnNext registers a callback action that will be called on each item emitted by the Observable.
func (o *ObservableImpl) DoOnNext(nextFunc NextFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-src:
				if !ok {
					return
				}
				if i.Error() {
					return
				}
				nextFunc(i.V)
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe())
	return dispose
}

// ElementAt emits only item n emitted by an Observable.
func (o *ObservableImpl) ElementAt(index uint, opts ...Option) Single {
	takeCount := 0
	sent := false

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if takeCount == int(index) {
			dst <- item
			sent = true
			operator.stop()
			return
		}
		takeCount++
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !sent {
			dst <- Error(&IllegalInputError{})
		}
	}, opts...)
}

// Filter emits only those items from an Observable that pass a predicate test.
func (o *ObservableImpl) Filter(apply Predicate, opts ...Option) Observable {
	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if apply(item.V) {
			dst <- item
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// First returns new Observable which emit only first item.
func (o *ObservableImpl) First(opts ...Option) OptionalSingle {
	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
		operator.stop()
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
func (o *ObservableImpl) FirstOrDefault(defaultValue interface{}, opts ...Option) Single {
	sent := false

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
		sent = true
		operator.stop()
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !sent {
			dst <- Of(defaultValue)
		}
	}, opts...)
}

// FlatMap transforms the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.
func (o *ObservableImpl) FlatMap(apply ItemToObservable, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		observe := o.Observe()
	loop1:
		for {
			select {
			case <-ctx.Done():
				break loop1
			case i, ok := <-observe:
				if !ok {
					break loop1
				}
				observe2 := apply(i).Observe()
			loop2:
				for {
					select {
					case <-ctx.Done():
						break loop2
					case i, ok := <-observe2:
						if !ok {
							break loop2
						}
						if i.Error() {
							next <- i
							return
						}
						next <- i
					}
				}
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// ForEach subscribes to the Observable and receives notifications for each element.
func (o *ObservableImpl) ForEach(nextFunc NextFunc, errFunc ErrFunc, completedFunc CompletedFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		for {
			select {
			case <-ctx.Done():
				completedFunc()
				return
			case i, ok := <-src:
				if !ok {
					completedFunc()
					return
				}
				if i.Error() {
					errFunc(i.E)
					break
				}
				nextFunc(i.V)
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe())
	return dispose
}

// IgnoreElements ignores all items emitted by the source ObservableSource and only calls onComplete
// or onError.
func (o *ObservableImpl) IgnoreElements(opts ...Option) Observable {
	return newObservableFromOperator(o, func(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// GroupBy divides an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key.
func (o *ObservableImpl) GroupBy(length int, distribution func(Item) int, opts ...Option) Observable {
	option := parseOptions(opts...)
	ctx := option.buildContext()

	s := make([]Item, length)
	chs := make([]chan Item, length)
	for i := 0; i < length; i++ {
		ch := option.buildChannel()
		chs[i] = ch
		s[i] = Of(&ObservableImpl{
			iterable: newChannelIterable(ch),
		})
	}

	go func() {
		observe := o.Observe()
		defer func() {
			for i := 0; i < length; i++ {
				close(chs[i])
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}
				idx := distribution(item)
				if idx >= length {
					err := Error(IndexOutOfBoundError{error: fmt.Sprintf("index %d, length %d", idx, length)})
					for i := 0; i < length; i++ {
						err.SendBlocking(chs[i])
					}
					return
				}
				item.SendBlocking(chs[idx])
			}
		}
	}()

	return &ObservableImpl{
		iterable: newSliceIterable(s, opts...),
	}
}

// Last returns a new Observable which emit only last item.
func (o *ObservableImpl) Last(opts ...Option) OptionalSingle {
	var last Item
	empty := true

	return newOptionalSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		last = item
		empty = false
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !empty {
			dst <- last
		}
	}, opts...)
}

// LastOrDefault returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
func (o *ObservableImpl) LastOrDefault(defaultValue interface{}, opts ...Option) Single {
	var last Item
	empty := true

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		last = item
		empty = false
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !empty {
			dst <- last
		} else {
			dst <- Of(defaultValue)
		}
	}, opts...)
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func (o *ObservableImpl) Map(apply Func, opts ...Option) Observable {
	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		res, err := apply(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
		}
		dst <- Of(res)
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// Marshal transforms the items emitted by an Observable by applying a marshalling to each item.
func (o *ObservableImpl) Marshal(marshaler Marshaler, opts ...Option) Observable {
	return o.Map(func(i interface{}) (interface{}, error) {
		return marshaler(i)
	}, opts...)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (o *ObservableImpl) Max(comparator Comparator, opts ...Option) OptionalSingle {
	empty := true
	var max interface{}

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		empty = false

		if max == nil {
			max = item.V
		} else {
			if comparator(max, item.V) < 0 {
				max = item.V
			}
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !empty {
			dst <- Of(max)
		}
	}, opts...)
}

// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
func (o *ObservableImpl) Min(comparator Comparator, opts ...Option) OptionalSingle {
	empty := true
	var min interface{}

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		empty = false

		if min == nil {
			min = item.V
		} else {
			if comparator(min, item.V) > 0 {
				min = item.V
			}
		}
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !empty {
			dst <- Of(min)
		}
	}, opts...)
}

// Observe observes an Observable by returning its channel.
func (o *ObservableImpl) Observe() <-chan Item {
	return o.iterable.Observe()
}

// OnErrorResumeNext instructs an Observable to pass control to another Observable rather than invoking
// onError if it encounters an error.
func (o *ObservableImpl) OnErrorResumeNext(resumeSequence ErrorToObservable, opts ...Option) Observable {
	return newObservableFromOperator(o, defaultNextFuncOperator, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		operator.resetIterable(resumeSequence(item.E))
	}, defaultEndFuncOperator, opts...)
}

// OnErrorReturn instructs an Observable to emit an item (returned by a specified function)
// rather than invoking onError if it encounters an error.
func (o *ObservableImpl) OnErrorReturn(resumeFunc ErrorFunc, opts ...Option) Observable {
	return newObservableFromOperator(o, defaultNextFuncOperator, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		dst <- Of(resumeFunc(item.E))
	}, defaultEndFuncOperator, opts...)
}

// OnErrorReturnItem instructs on observale to emit an item if it encounters an error.
func (o *ObservableImpl) OnErrorReturnItem(resume interface{}, opts ...Option) Observable {
	return newObservableFromOperator(o, defaultNextFuncOperator, func(_ context.Context, _ Item, dst chan<- Item, operator operatorOptions) {
		dst <- Of(resume)
	}, defaultEndFuncOperator, opts...)
}

// Reduce applies a function to each item emitted by an Observable, sequentially, and emit the final value.
func (o *ObservableImpl) Reduce(apply Func2, opts ...Option) OptionalSingle {
	var acc interface{}
	empty := true

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		empty = false
		v, err := apply(acc, item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}
		acc = v
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		if !empty {
			dst <- Of(acc)
		}
	}, opts...)
}

// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
// at most count times, at a particular frequency.
func (o *ObservableImpl) Repeat(count int64, frequency Duration, opts ...Option) Observable {
	if count != Infinite {
		if count < 0 {
			return newObservableFromError(IllegalInputError{error: "count must be positive"})
		}
	}
	seq := make([]Item, 0)

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
		seq = append(seq, item)
	}, defaultErrorFuncOperator, func(ctx context.Context, dst chan<- Item) {
		for {
			select {
			default:
				break
			case <-ctx.Done():
				return
			}
			if count != Infinite {
				if count == 0 {
					break
				}
			}
			if frequency != nil {
				time.Sleep(frequency.duration())
			}
			for _, v := range seq {
				dst <- v
			}
			count = count - 1
		}
	}, opts...)
}

// BackOffRetry implements a backoff retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
func (o *ObservableImpl) BackOffRetry(backOffCfg backoff.BackOff, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	f := func() error {
		observe := o.Observe()
		for {
			select {
			case <-ctx.Done():
				close(next)
				return nil
			case i, ok := <-observe:
				if !ok {
					return nil
				}
				if i.Error() {
					return i.E
				}
				next <- i
			}
		}
	}
	go func() {
		if err := backoff.Retry(f, backOffCfg); err != nil {
			next <- Error(err)
			close(next)
			return
		}
		close(next)
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Retry retries if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
func (o *ObservableImpl) Retry(count int, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		observe := o.Observe()
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.Error() {
					count--
					if count < 0 {
						next <- i
						break loop
					}
					observe = o.Observe()
				} else {
					next <- i
				}
			}
		}
		close(next)
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Run creates an observer without consuming the emitted items.
func (o *ObservableImpl) Run(opts ...Option) Disposed {
	dispose := make(chan struct{})
	option := parseOptions(opts...)
	ctx := option.buildContext()

	go func() {
		defer close(dispose)
		observe := o.Observe()
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-observe:
				if !ok {
					return
				}
			}
		}
	}()

	return dispose
}

// Sample returns an Observable that emits the most recent items emitted by the source
// Iterable whenever the input Iterable emits an item.
func (o *ObservableImpl) Sample(iterable Iterable, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()
	itCh := make(chan Item)
	obsCh := make(chan Item)

	go func() {
		defer close(obsCh)
		observe := o.Observe()
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				obsCh <- i
			}
		}
	}()

	go func() {
		defer close(itCh)
		observe := iterable.Observe()
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				itCh <- i
			}
		}
	}()

	go func() {
		defer close(next)
		var lastEmittedItem Item
		isItemWaitingToBeEmitted := false

		for {
			select {
			case _, ok := <-itCh:
				if ok {
					if isItemWaitingToBeEmitted {
						next <- lastEmittedItem
						isItemWaitingToBeEmitted = false
					}
				} else {
					return
				}
			case item, ok := <-obsCh:
				if ok {
					lastEmittedItem = item
					isItemWaitingToBeEmitted = true
				} else {
					return
				}
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Scan applies Func2 predicate to each item in the original
// Observable sequentially and emits each successive value on a new Observable.
func (o *ObservableImpl) Scan(apply Func2, opts ...Option) Observable {
	var current interface{}

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		v, err := apply(current, item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}
		dst <- Of(v)
		current = v
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
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
func (o *ObservableImpl) SequenceEqual(iterable Iterable, opts ...Option) Single {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()
	itCh := make(chan Item)
	obsCh := make(chan Item)

	go func() {
		defer close(obsCh)
		observe := o.Observe()
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				obsCh <- i
			}
		}
	}()

	go func() {
		defer close(itCh)
		observe := iterable.Observe()
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				itCh <- i
			}
		}
	}()

	go func() {
		var mainSequence []interface{}
		var obsSequence []interface{}
		areCorrect := true
		isMainChannelClosed := false
		isObsChannelClosed := false

	mainLoop:
		for {
			select {
			case item, ok := <-itCh:
				if ok {
					mainSequence = append(mainSequence, item)
					areCorrect, mainSequence, obsSequence = popAndCompareFirstItems(mainSequence, obsSequence)
				} else {
					isMainChannelClosed = true
				}
			case item, ok := <-obsCh:
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

		next <- Of(areCorrect && len(mainSequence) == 0 && len(obsSequence) == 0)
		close(next)
	}()

	return &SingleImpl{
		iterable: newChannelIterable(next),
	}
}

// Send sends the items to a given channel
func (o *ObservableImpl) Send(output chan<- Item, opts ...Option) {
	go func() {
		option := parseOptions(opts...)
		ctx := option.buildContext()
		observe := o.Observe()
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.Error() {
					output <- i
					break loop
				}
				output <- i
			}
		}
		close(output)
	}()
}

// Serialize forces an Observable to make serialized calls and to be well-behaved.
func (o *ObservableImpl) Serialize(from int, identifier func(interface{}) int, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()

	ctx := option.buildContext()
	mutex := sync.Mutex{}
	minHeap := binaryheap.NewWith(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})
	minHeap.Push(from)
	counter := int64(from)
	status := make(map[int]interface{})
	notif := make(chan struct{})

	// Scatter
	go func() {
		defer close(notif)
		src := o.Observe()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-src:
				if !ok {
					return
				}
				if item.Error() {
					next <- item
					return
				}

				id := identifier(item.V)
				mutex.Lock()
				if id != from {
					minHeap.Push(id)
				}
				status[id] = item.V
				mutex.Unlock()
				notif <- struct{}{}
			}
		}
	}()

	// Gather
	go func() {
		defer close(next)

		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-notif:
				if !ok {
					return
				}

				mutex.Lock()
				for !minHeap.Empty() {
					v, _ := minHeap.Peek()
					id := v.(int)
					if atomic.LoadInt64(&counter) == int64(id) {
						if itemValue, contains := status[id]; contains {
							minHeap.Pop()
							delete(status, id)
							mutex.Unlock()
							next <- Of(itemValue)
							mutex.Lock()
							atomic.AddInt64(&counter, 1)
							continue
						}
					}
					break
				}
				mutex.Unlock()
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
func (o *ObservableImpl) Skip(nth uint, opts ...Option) Observable {
	skipCount := 0

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if skipCount < int(nth) {
			skipCount++
			return
		}
		dst <- item
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
func (o *ObservableImpl) SkipLast(nth uint, opts ...Option) Observable {
	skipCount := 0

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if skipCount >= int(nth) {
			operator.stop()
			return
		}
		skipCount++
		dst <- item
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
func (o *ObservableImpl) SkipWhile(apply Predicate, opts ...Option) Observable {
	skip := true

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if !skip {
			dst <- item
		} else {
			if !apply(item.V) {
				skip = false
				dst <- item
			}
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// StartWithIterable returns an Observable that emits the items in a specified Iterable before it begins to
// emit items emitted by the source Observable.
func (o *ObservableImpl) StartWithIterable(iterable Iterable, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		observe := iterable.Observe()
	loop1:
		for {
			select {
			case <-ctx.Done():
				break loop1
			case i, ok := <-observe:
				if !ok {
					break loop1
				}
				if i.Error() {
					next <- i
					return
				}
				next <- i
			}
		}
		observe = o.Observe()
	loop2:
		for {
			select {
			case <-ctx.Done():
				break loop2
			case i, ok := <-observe:
				if !ok {
					break loop2
				}
				if i.Error() {
					next <- i
					return
				}
				next <- i
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// SumFloat32 calculates the average of float32 emitted by an Observable and emits a float32.
func (o *ObservableImpl) SumFloat32(opts ...Option) Single {
	var sum float32

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		switch i := item.V.(type) {
		default:
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: (float32|int|int8|int16|int32|int64), got: %t", item)})
			operator.stop()
			return
		case int:
			sum += float32(i)
		case int8:
			sum += float32(i)
		case int16:
			sum += float32(i)
		case int32:
			sum += float32(i)
		case int64:
			sum += float32(i)
		case float32:
			sum += i
		}
	}, defaultErrorFuncOperator, func(ctx context.Context, dst chan<- Item) {
		dst <- Of(sum)
	}, opts...)
}

// SumFloat64 calculates the average of float64 emitted by an Observable and emits a float64.
func (o *ObservableImpl) SumFloat64(opts ...Option) Single {
	var sum float64

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		switch i := item.V.(type) {
		default:
			dst <- Error(&IllegalInputError{error: fmt.Sprintf("expected type: (float32|float64|int|int8|int16|int32|int64), got: %t", item)})
			operator.stop()
			return
		case int:
			sum += float64(i)
		case int8:
			sum += float64(i)
		case int16:
			sum += float64(i)
		case int32:
			sum += float64(i)
		case int64:
			sum += float64(i)
		case float32:
			sum += float64(i)
		case float64:
			sum += i
		}
	}, defaultErrorFuncOperator, func(ctx context.Context, dst chan<- Item) {
		dst <- Of(sum)
	}, opts...)
}

// SumInt64 calculates the average of integers emitted by an Observable and emits an int64.
func (o *ObservableImpl) SumInt64(opts ...Option) Single {
	var sum int64

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		switch i := item.V.(type) {
		default:
			dst <- Error(IllegalInputError{error: fmt.Sprintf("expected type: (int|int8|int16|int32|int64), got: %t", item)})
			operator.stop()
			return
		case int:
			sum += int64(i)
		case int8:
			sum += int64(i)
		case int16:
			sum += int64(i)
		case int32:
			sum += int64(i)
		case int64:
			sum += i
		}
	}, defaultErrorFuncOperator, func(ctx context.Context, dst chan<- Item) {
		dst <- Of(sum)
	}, opts...)
}

// Take emits only the first n items emitted by an Observable.
func (o *ObservableImpl) Take(nth uint, opts ...Option) Observable {
	takeCount := 0

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if takeCount < int(nth) {
			takeCount++
			dst <- item
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// TakeLast emits only the last n items emitted by an Observable.
func (o *ObservableImpl) TakeLast(nth uint, opts ...Option) Observable {
	n := int(nth)
	r := ring.New(n)
	count := 0

	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		count++
		r.Value = item.V
		r = r.Next()
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
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
			dst <- Of(r.Value)
			r = r.Next()
		}
	}, opts...)
}

// TakeUntil returns an Observable that emits items emitted by the source Observable,
// checks the specified predicate for each item, and then completes when the condition is satisfied.
func (o *ObservableImpl) TakeUntil(apply Predicate, opts ...Option) Observable {
	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
		if apply(item.V) {
			operator.stop()
			return
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// TakeWhile returns an Observable that emits items emitted by the source ObservableSource so long as each
// item satisfied a specified condition, and then completes as soon as this condition is not satisfied.
func (o *ObservableImpl) TakeWhile(apply Predicate, opts ...Option) Observable {
	return newObservableFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		if !apply(item.V) {
			operator.stop()
			return
		}
		dst <- item
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// ToMap convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function.
func (o *ObservableImpl) ToMap(keySelector Func, opts ...Option) Single {
	m := make(map[interface{}]interface{})

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		k, err := keySelector(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}
		m[k] = item.V
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		dst <- Of(m)
	}, opts...)
}

// ToMapWithValueSelector convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function and valued by another
// value function.
func (o *ObservableImpl) ToMapWithValueSelector(keySelector, valueSelector Func, opts ...Option) Single {
	m := make(map[interface{}]interface{})

	return newSingleFromOperator(o, func(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
		k, err := keySelector(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}

		v, err := valueSelector(item.V)
		if err != nil {
			dst <- Error(err)
			operator.stop()
			return
		}

		m[k] = v
	}, defaultErrorFuncOperator, func(_ context.Context, dst chan<- Item) {
		dst <- Of(m)
	}, opts...)
}

// ToSlice collects all items from an Observable and emit them in a slice and an optional error.
func (o *ObservableImpl) ToSlice(initialCapacity int, opts ...Option) ([]interface{}, error) {
	s := make([]interface{}, 0, initialCapacity)
	var err error
	<-newObservableFromOperator(o, func(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
		s = append(s, item.V)
	}, func(_ context.Context, item Item, _ chan<- Item, operator operatorOptions) {
		err = item.E
		operator.stop()
	}, defaultEndFuncOperator, opts...).Run()
	return s, err
}

// Unmarshal transforms the items emitted by an Observable by applying an unmarshalling to each item.
func (o *ObservableImpl) Unmarshal(unmarshaler Unmarshaler, factory func() interface{}, opts ...Option) Observable {
	return o.Map(func(i interface{}) (interface{}, error) {
		v := factory()
		err := unmarshaler(i.([]byte), v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}, opts...)
}

// ZipFromIterable merges the emissions of multiple Observables together via a specified function
// and emit single items for each combination based on the results of this function.
func (o *ObservableImpl) ZipFromIterable(iterable Iterable, zipper Func2, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		it1 := o.Observe()
		it2 := iterable.Observe()
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i1, ok := <-it1:
				if !ok {
					break loop
				}
				if i1.Error() {
					next <- i1
					return
				}
				for {
					select {
					case <-ctx.Done():
						break loop
					case i2, ok := <-it2:
						if !ok {
							break loop
						}
						if i2.Error() {
							next <- i2
							return
						}
						v, err := zipper(i1.V, i2.V)
						if err != nil {
							next <- Error(err)
							return
						}
						next <- Of(v)
						continue loop
					}
				}
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}
