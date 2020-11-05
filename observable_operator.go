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

// All determines whether all items emitted by an Observable meet some criteria.
func (o *ObservableImpl) All(predicate Predicate, opts ...Option) Single {
	return single(o, func() operator {
		return &allOperator{
			predicate: predicate,
			all:       true,
		}
	}, false, false, opts...)
}

type allOperator struct {
	predicate Predicate
	all       bool
}

func (op *allOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if !op.predicate(item.V) {
		Of(false).SendContext(ctx, dst)
		op.all = false
		operatorOptions.stop()
	}
}

func (op *allOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *allOperator) end(ctx context.Context, dst chan<- Item) {
	if op.all {
		Of(true).SendContext(ctx, dst)
	}
}

func (op *allOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if item.V == false {
		Of(false).SendContext(ctx, dst)
		op.all = false
		operatorOptions.stop()
	}
}

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits the average float32.
func (o *ObservableImpl) AverageFloat32(opts ...Option) Single {
	return single(o, func() operator {
		return &averageFloat32Operator{}
	}, false, false, opts...)
}

type averageFloat32Operator struct {
	sum   float32
	count float32
}

func (op *averageFloat32Operator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: float or int, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int:
		op.sum += float32(v)
		op.count++
	case float32:
		op.sum += v
		op.count++
	case float64:
		op.sum += float32(v)
		op.count++
	}
}

func (op *averageFloat32Operator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageFloat32Operator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageFloat32Operator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageFloat32Operator)
	op.sum += v.sum
	op.count += v.count
}

// AverageFloat64 calculates the average of numbers emitted by an Observable and emits the average float64.
func (o *ObservableImpl) AverageFloat64(opts ...Option) Single {
	return single(o, func() operator {
		return &averageFloat64Operator{}
	}, false, false, opts...)
}

type averageFloat64Operator struct {
	sum   float64
	count float64
}

func (op *averageFloat64Operator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: float or int, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int:
		op.sum += float64(v)
		op.count++
	case float32:
		op.sum += float64(v)
		op.count++
	case float64:
		op.sum += v
		op.count++
	}
}

func (op *averageFloat64Operator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageFloat64Operator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageFloat64Operator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageFloat64Operator)
	op.sum += v.sum
	op.count += v.count
}

// AverageInt calculates the average of numbers emitted by an Observable and emits the average int.
func (o *ObservableImpl) AverageInt(opts ...Option) Single {
	return single(o, func() operator {
		return &averageIntOperator{}
	}, false, false, opts...)
}

type averageIntOperator struct {
	sum   int
	count int
}

func (op *averageIntOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: int, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int:
		op.sum += v
		op.count++
	}
}

func (op *averageIntOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageIntOperator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageIntOperator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageIntOperator)
	op.sum += v.sum
	op.count += v.count
}

// AverageInt8 calculates the average of numbers emitted by an Observable and emits the≤ average int8.
func (o *ObservableImpl) AverageInt8(opts ...Option) Single {
	return single(o, func() operator {
		return &averageInt8Operator{}
	}, false, false, opts...)
}

type averageInt8Operator struct {
	sum   int8
	count int8
}

func (op *averageInt8Operator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: int8, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int8:
		op.sum += v
		op.count++
	}
}

func (op *averageInt8Operator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageInt8Operator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageInt8Operator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageInt8Operator)
	op.sum += v.sum
	op.count += v.count
}

// AverageInt16 calculates the average of numbers emitted by an Observable and emits the average int16.
func (o *ObservableImpl) AverageInt16(opts ...Option) Single {
	return single(o, func() operator {
		return &averageInt16Operator{}
	}, false, false, opts...)
}

type averageInt16Operator struct {
	sum   int16
	count int16
}

func (op *averageInt16Operator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: int16, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int16:
		op.sum += v
		op.count++
	}
}

func (op *averageInt16Operator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageInt16Operator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageInt16Operator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageInt16Operator)
	op.sum += v.sum
	op.count += v.count
}

// AverageInt32 calculates the average of numbers emitted by an Observable and emits the average int32.
func (o *ObservableImpl) AverageInt32(opts ...Option) Single {
	return single(o, func() operator {
		return &averageInt32Operator{}
	}, false, false, opts...)
}

type averageInt32Operator struct {
	sum   int32
	count int32
}

func (op *averageInt32Operator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: int32, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int32:
		op.sum += v
		op.count++
	}
}

func (op *averageInt32Operator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageInt32Operator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageInt32Operator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageInt32Operator)
	op.sum += v.sum
	op.count += v.count
}

// AverageInt64 calculates the average of numbers emitted by an Observable and emits this average int64.
func (o *ObservableImpl) AverageInt64(opts ...Option) Single {
	return single(o, func() operator {
		return &averageInt64Operator{}
	}, false, false, opts...)
}

type averageInt64Operator struct {
	sum   int64
	count int64
}

func (op *averageInt64Operator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	switch v := item.V.(type) {
	default:
		Error(IllegalInputError{error: fmt.Sprintf("expected type: int64, got: %t", item)}).SendContext(ctx, dst)
		operatorOptions.stop()
	case int64:
		op.sum += v
		op.count++
	}
}

func (op *averageInt64Operator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *averageInt64Operator) end(ctx context.Context, dst chan<- Item) {
	if op.count == 0 {
		Of(0).SendContext(ctx, dst)
	} else {
		Of(op.sum/op.count).SendContext(ctx, dst)
	}
}

func (op *averageInt64Operator) gatherNext(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	v := item.V.(*averageInt64Operator)
	op.sum += v.sum
	op.count += v.count
}

// BackOffRetry implements a backoff retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
// Cannot be run in parallel.
func (o *ObservableImpl) BackOffRetry(backOffCfg backoff.BackOff, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	f := func() error {
		observe := o.Observe(opts...)
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
				i.SendContext(ctx, next)
			}
		}
	}
	go func() {
		if err := backoff.Retry(f, backOffCfg); err != nil {
			Error(err).SendContext(ctx, next)
			close(next)
			return
		}
		close(next)
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// BufferWithCount returns an Observable that emits buffers of items it collects
// from the source Observable.
// The resulting Observable emits buffers every skip items, each containing a slice of count items.
// When the source Observable completes or encounters an error,
// the resulting Observable emits the current buffer and propagates
// the notification from the source Observable.
func (o *ObservableImpl) BufferWithCount(count int, opts ...Option) Observable {
	if count <= 0 {
		return Thrown(IllegalInputError{error: "count must be positive"})
	}

	return observable(o, func() operator {
		return &bufferWithCountOperator{
			count:  count,
			buffer: make([]interface{}, count),
		}
	}, true, false, opts...)
}

type bufferWithCountOperator struct {
	count  int
	iCount int
	buffer []interface{}
}

func (op *bufferWithCountOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	op.buffer[op.iCount] = item.V
	op.iCount++
	if op.iCount == op.count {
		Of(op.buffer).SendContext(ctx, dst)
		op.iCount = 0
		op.buffer = make([]interface{}, op.count)
	}
}

func (op *bufferWithCountOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *bufferWithCountOperator) end(ctx context.Context, dst chan<- Item) {
	if op.iCount != 0 {
		Of(op.buffer[:op.iCount]).SendContext(ctx, dst)
	}
}

func (op *bufferWithCountOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (o *ObservableImpl) BufferWithTime(timespan Duration, opts ...Option) Observable {
	if timespan == nil {
		return Thrown(IllegalInputError{error: "timespan must no be nil"})
	}

	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		observe := o.Observe(opts...)
		buffer := make([]interface{}, 0)
		stop := make(chan struct{})
		mutex := sync.Mutex{}

		checkBuffer := func() {
			mutex.Lock()
			if len(buffer) != 0 {
				if !Of(buffer).SendContext(ctx, next) {
					mutex.Unlock()
					return
				}
				buffer = make([]interface{}, 0)
			}
			mutex.Unlock()
		}

		go func() {
			defer close(next)
			duration := timespan.duration()
			for {
				select {
				case <-stop:
					checkBuffer()
					return
				case <-ctx.Done():
					return
				case <-time.After(duration):
					checkBuffer()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				close(stop)
				return
			case item, ok := <-observe:
				if !ok {
					close(stop)
					return
				}
				if item.Error() {
					item.SendContext(ctx, next)
					if option.getErrorStrategy() == StopOnError {
						close(stop)
						return
					}
				} else {
					mutex.Lock()
					buffer = append(buffer, item.V)
					mutex.Unlock()
				}
			}
		}
	}

	return customObservableOperator(f, opts...)
}

// BufferWithTimeOrCount returns an Observable that emits buffers of items it collects from the source
// Observable either from a given count or at a given time interval.
func (o *ObservableImpl) BufferWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable {
	if timespan == nil {
		return Thrown(IllegalInputError{error: "timespan must no be nil"})
	}
	if count <= 0 {
		return Thrown(IllegalInputError{error: "count must be positive"})
	}

	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		observe := o.Observe(opts...)
		buffer := make([]interface{}, 0)
		stop := make(chan struct{})
		send := make(chan struct{})
		mutex := sync.Mutex{}

		checkBuffer := func() {
			mutex.Lock()
			if len(buffer) != 0 {
				if !Of(buffer).SendContext(ctx, next) {
					mutex.Unlock()
					return
				}
				buffer = make([]interface{}, 0)
			}
			mutex.Unlock()
		}

		go func() {
			defer close(next)
			duration := timespan.duration()
			for {
				select {
				case <-send:
					checkBuffer()
				case <-stop:
					checkBuffer()
					return
				case <-ctx.Done():
					return
				case <-time.After(duration):
					checkBuffer()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					close(stop)
					close(send)
					return
				}
				if item.Error() {
					item.SendContext(ctx, next)
					if option.getErrorStrategy() == StopOnError {
						close(stop)
						close(send)
						return
					}
				} else {
					mutex.Lock()
					buffer = append(buffer, item.V)
					if len(buffer) == count {
						mutex.Unlock()
						send <- struct{}{}
					} else {
						mutex.Unlock()
					}
				}
			}
		}
	}

	return customObservableOperator(f, opts...)
}

// Connect instructs a connectable Observable to begin emitting items to its subscribers.
func (o *ObservableImpl) Connect(ctx context.Context) (context.Context, Disposable) {
	ctx, cancel := context.WithCancel(ctx)
	o.Observe(WithContext(ctx), connect())
	return ctx, Disposable(cancel)
}

// Contains determines whether an Observable emits a particular item or not.
func (o *ObservableImpl) Contains(equal Predicate, opts ...Option) Single {
	return single(o, func() operator {
		return &containsOperator{
			equal:    equal,
			contains: false,
		}
	}, false, false, opts...)
}

type containsOperator struct {
	equal    Predicate
	contains bool
}

func (op *containsOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if op.equal(item.V) {
		Of(true).SendContext(ctx, dst)
		op.contains = true
		operatorOptions.stop()
	}
}

func (op *containsOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *containsOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.contains {
		Of(false).SendContext(ctx, dst)
	}
}

func (op *containsOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if item.V == true {
		Of(true).SendContext(ctx, dst)
		operatorOptions.stop()
		op.contains = true
	}
}

// Count counts the number of items emitted by the source Observable and emit only this value.
func (o *ObservableImpl) Count(opts ...Option) Single {
	return single(o, func() operator {
		return &countOperator{}
	}, true, false, opts...)
}

type countOperator struct {
	count int64
}

func (op *countOperator) next(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
	op.count++
}

func (op *countOperator) err(_ context.Context, _ Item, _ chan<- Item, operatorOptions operatorOptions) {
	operatorOptions.stop()
}

func (op *countOperator) end(ctx context.Context, dst chan<- Item) {
	Of(op.count).SendContext(ctx, dst)
}

func (op *countOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Debounce only emits an item from an Observable if a particular timespan has passed without it emitting another item.
func (o *ObservableImpl) Debounce(timespan Duration, opts ...Option) Observable {
	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		defer close(next)
		observe := o.Observe(opts...)
		var latest interface{}

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}
				if item.Error() {
					if !item.SendContext(ctx, next) {
						return
					}
					if option.getErrorStrategy() == StopOnError {
						return
					}
				} else {
					latest = item.V
				}
			case <-time.After(timespan.duration()):
				if latest != nil {
					if !Of(latest).SendContext(ctx, next) {
						return
					}
					latest = nil
				}
			}
		}
	}

	return customObservableOperator(f, opts...)
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (o *ObservableImpl) DefaultIfEmpty(defaultValue interface{}, opts ...Option) Observable {
	return observable(o, func() operator {
		return &defaultIfEmptyOperator{
			defaultValue: defaultValue,
			empty:        true,
		}
	}, true, false, opts...)
}

type defaultIfEmptyOperator struct {
	defaultValue interface{}
	empty        bool
}

func (op *defaultIfEmptyOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	op.empty = false
	item.SendContext(ctx, dst)
}

func (op *defaultIfEmptyOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *defaultIfEmptyOperator) end(ctx context.Context, dst chan<- Item) {
	if op.empty {
		Of(op.defaultValue).SendContext(ctx, dst)
	}
}

func (op *defaultIfEmptyOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *ObservableImpl) Distinct(apply Func, opts ...Option) Observable {
	return observable(o, func() operator {
		return &distinctOperator{
			apply:  apply,
			keyset: make(map[interface{}]interface{}),
		}
	}, false, false, opts...)
}

type distinctOperator struct {
	apply  Func
	keyset map[interface{}]interface{}
}

func (op *distinctOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	key, err := op.apply(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}
	_, ok := op.keyset[key]
	if !ok {
		item.SendContext(ctx, dst)
	}
	op.keyset[key] = nil
}

func (op *distinctOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *distinctOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *distinctOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	switch item.V.(type) {
	case *distinctOperator:
		return
	}

	if _, contains := op.keyset[item.V]; !contains {
		Of(item.V).SendContext(ctx, dst)
		op.keyset[item.V] = nil
	}
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original Observable.
// Cannot be run in parallel.
func (o *ObservableImpl) DistinctUntilChanged(apply Func, opts ...Option) Observable {
	return observable(o, func() operator {
		return &distinctUntilChangedOperator{
			apply: apply,
		}
	}, true, false, opts...)
}

type distinctUntilChangedOperator struct {
	apply   Func
	current interface{}
}

func (op *distinctUntilChangedOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	key, err := op.apply(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}
	if op.current != key {
		item.SendContext(ctx, dst)
		op.current = key
	}
}

func (op *distinctUntilChangedOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *distinctUntilChangedOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *distinctUntilChangedOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
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
	go handler(ctx, o.Observe(opts...))
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
	go handler(ctx, o.Observe(opts...))
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
	go handler(ctx, o.Observe(opts...))
	return dispose
}

// ElementAt emits only item n emitted by an Observable.
// Cannot be run in parallel.
func (o *ObservableImpl) ElementAt(index uint, opts ...Option) Single {
	return single(o, func() operator {
		return &elementAtOperator{
			index: index,
		}
	}, true, false, opts...)
}

type elementAtOperator struct {
	index     uint
	takeCount int
	sent      bool
}

func (op *elementAtOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if op.takeCount == int(op.index) {
		item.SendContext(ctx, dst)
		op.sent = true
		operatorOptions.stop()
		return
	}
	op.takeCount++
}

func (op *elementAtOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *elementAtOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.sent {
		Error(&IllegalInputError{}).SendContext(ctx, dst)
	}
}

func (op *elementAtOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Error returns the eventual Observable error.
// This method is blocking.
func (o *ObservableImpl) Error(opts ...Option) error {
	option := parseOptions(opts...)
	ctx := option.buildContext()
	observe := o.iterable.Observe(opts...)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-observe:
			if !ok {
				return nil
			}
			if item.Error() {
				return item.E
			}
		}
	}
}

// Errors returns an eventual list of Observable errors.
// This method is blocking
func (o *ObservableImpl) Errors(opts ...Option) []error {
	option := parseOptions(opts...)
	ctx := option.buildContext()
	observe := o.iterable.Observe(opts...)
	errs := make([]error, 0)

	for {
		select {
		case <-ctx.Done():
			return []error{ctx.Err()}
		case item, ok := <-observe:
			if !ok {
				return errs
			}
			if item.Error() {
				errs = append(errs, item.E)
			}
		}
	}
}

// Filter emits only those items from an Observable that pass a predicate test.
func (o *ObservableImpl) Filter(apply Predicate, opts ...Option) Observable {
	return observable(o, func() operator {
		return &filterOperator{apply: apply}
	}, false, true, opts...)
}

type filterOperator struct {
	apply Predicate
}

func (op *filterOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	if op.apply(item.V) {
		item.SendContext(ctx, dst)
	}
}

func (op *filterOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *filterOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *filterOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Find emits the first item passing a predicate then complete.
func (o *ObservableImpl) Find(find Predicate, opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &findOperator{
			find: find,
		}
	}, true, true, opts...)
}

type findOperator struct {
	find Predicate
}

func (op *findOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if op.find(item.V) {
		item.SendContext(ctx, dst)
		operatorOptions.stop()
	}
}

func (op *findOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *findOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *findOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// First returns new Observable which emit only first item.
// Cannot be run in parallel.
func (o *ObservableImpl) First(opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &firstOperator{}
	}, true, false, opts...)
}

type firstOperator struct{}

func (op *firstOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	item.SendContext(ctx, dst)
	operatorOptions.stop()
}

func (op *firstOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *firstOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *firstOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
// Cannot be run in parallel.
func (o *ObservableImpl) FirstOrDefault(defaultValue interface{}, opts ...Option) Single {
	return single(o, func() operator {
		return &firstOrDefaultOperator{
			defaultValue: defaultValue,
		}
	}, true, false, opts...)
}

type firstOrDefaultOperator struct {
	defaultValue interface{}
	sent         bool
}

func (op *firstOrDefaultOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	item.SendContext(ctx, dst)
	op.sent = true
	operatorOptions.stop()
}

func (op *firstOrDefaultOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *firstOrDefaultOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.sent {
		Of(op.defaultValue).SendContext(ctx, dst)
	}
}

func (op *firstOrDefaultOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// FlatMap transforms the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.
func (o *ObservableImpl) FlatMap(apply ItemToObservable, opts ...Option) Observable {
	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		defer close(next)
		observe := o.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}
				observe2 := apply(item).Observe(opts...)
			loop2:
				for {
					select {
					case <-ctx.Done():
						return
					case item, ok := <-observe2:
						if !ok {
							break loop2
						}
						if item.Error() {
							item.SendContext(ctx, next)
							if option.getErrorStrategy() == StopOnError {
								return
							}
						} else {
							if !item.SendContext(ctx, next) {
								return
							}
						}
					}
				}
			}
		}
	}

	return customObservableOperator(f, opts...)
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
	go handler(ctx, o.Observe(opts...))
	return dispose
}

// IgnoreElements ignores all items emitted by the source ObservableSource except for the errors.
// Cannot be run in parallel.
func (o *ObservableImpl) IgnoreElements(opts ...Option) Observable {
	return observable(o, func() operator {
		return &ignoreElementsOperator{}
	}, true, false, opts...)
}

type ignoreElementsOperator struct{}

func (op *ignoreElementsOperator) next(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

func (op *ignoreElementsOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *ignoreElementsOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *ignoreElementsOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Returns absolute value for int64
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

// Join combines items emitted by two Observables whenever an item from one Observable is emitted during
// a time window defined according to an item emitted by the other Observable.
// The time is extracted using a timeExtractor function.
func (o *ObservableImpl) Join(joiner Func2, right Observable, timeExtractor func(interface{}) time.Time, window Duration, opts ...Option) Observable {
	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		defer close(next)
		windowDuration := int64(window.duration())
		rBuf := make([]Item, 0)

		lObserve := o.Observe()
		rObserve := right.Observe()
	lLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case lItem, ok := <-lObserve:
				if lItem.V == nil && !ok {
					return
				}
				if lItem.Error() {
					lItem.SendContext(ctx, next)
					if option.getErrorStrategy() == StopOnError {
						return
					}
					continue
				}
				lTime := timeExtractor(lItem.V).UnixNano()
				cutPoint := 0
				for i, rItem := range rBuf {
					rTime := timeExtractor(rItem.V).UnixNano()
					if abs(lTime-rTime) <= windowDuration {
						i, err := joiner(ctx, lItem.V, rItem.V)
						if err != nil {
							Error(err).SendContext(ctx, next)
							if option.getErrorStrategy() == StopOnError {
								return
							}
							continue
						}
						Of(i).SendContext(ctx, next)
					}
					if lTime > rTime+windowDuration {
						cutPoint = i + 1
					}
				}

				rBuf = rBuf[cutPoint:]

				for {
					select {
					case <-ctx.Done():
						return
					case rItem, ok := <-rObserve:
						if rItem.V == nil && !ok {
							continue lLoop
						}
						if rItem.Error() {
							rItem.SendContext(ctx, next)
							if option.getErrorStrategy() == StopOnError {
								return
							}
							continue
						}

						rBuf = append(rBuf, rItem)
						rTime := timeExtractor(rItem.V).UnixNano()
						if abs(lTime-rTime) <= windowDuration {
							i, err := joiner(ctx, lItem.V, rItem.V)
							if err != nil {
								Error(err).SendContext(ctx, next)
								if option.getErrorStrategy() == StopOnError {
									return
								}
								continue
							}
							Of(i).SendContext(ctx, next)

							continue
						}
						continue lLoop
					}
				}
			}
		}
	}

	return customObservableOperator(f, opts...)
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
		observe := o.Observe(opts...)
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
						err.SendContext(ctx, chs[i])
					}
					return
				}
				item.SendContext(ctx, chs[idx])
			}
		}
	}()

	return &ObservableImpl{
		iterable: newSliceIterable(s, opts...),
	}
}

// GroupedObservable is the observable type emitted by the GroupByDynamic operator.
type GroupedObservable struct {
	Observable
	// Key is the distribution key
	Key string
}

// GroupByDynamic divides an Observable into a dynamic set of Observables that each emit GroupedObservable from the original Observable, organized by key.
func (o *ObservableImpl) GroupByDynamic(distribution func(Item) string, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()
	chs := make(map[string]chan Item)

	go func() {
		observe := o.Observe(opts...)
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				idx := distribution(i)
				ch, contains := chs[idx]
				if !contains {
					ch = option.buildChannel()
					chs[idx] = ch
					Of(GroupedObservable{
						Observable: &ObservableImpl{
							iterable: newChannelIterable(ch),
						},
						Key: idx,
					}).SendContext(ctx, next)
				}
				i.SendContext(ctx, ch)
			}
		}
		for _, ch := range chs {
			close(ch)
		}
		close(next)
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Last returns a new Observable which emit only last item.
// Cannot be run in parallel.
func (o *ObservableImpl) Last(opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &lastOperator{
			empty: true,
		}
	}, true, false, opts...)
}

type lastOperator struct {
	last  Item
	empty bool
}

func (op *lastOperator) next(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	op.last = item
	op.empty = false
}

func (op *lastOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *lastOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.empty {
		op.last.SendContext(ctx, dst)
	}
}

func (op *lastOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// LastOrDefault returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
// Cannot be run in parallel.
func (o *ObservableImpl) LastOrDefault(defaultValue interface{}, opts ...Option) Single {
	return single(o, func() operator {
		return &lastOrDefaultOperator{
			defaultValue: defaultValue,
			empty:        true,
		}
	}, true, false, opts...)
}

type lastOrDefaultOperator struct {
	defaultValue interface{}
	last         Item
	empty        bool
}

func (op *lastOrDefaultOperator) next(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	op.last = item
	op.empty = false
}

func (op *lastOrDefaultOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *lastOrDefaultOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.empty {
		op.last.SendContext(ctx, dst)
	} else {
		Of(op.defaultValue).SendContext(ctx, dst)
	}
}

func (op *lastOrDefaultOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func (o *ObservableImpl) Map(apply Func, opts ...Option) Observable {
	return observable(o, func() operator {
		return &mapOperator{apply: apply}
	}, false, true, opts...)
}

type mapOperator struct {
	apply Func
}

func (op *mapOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	res, err := op.apply(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}
	Of(res).SendContext(ctx, dst)
}

func (op *mapOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *mapOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *mapOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	switch item.V.(type) {
	case *mapOperator:
		return
	}
	item.SendContext(ctx, dst)
}

// Marshal transforms the items emitted by an Observable by applying a marshalling to each item.
func (o *ObservableImpl) Marshal(marshaller Marshaller, opts ...Option) Observable {
	return o.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return marshaller(i)
	}, opts...)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (o *ObservableImpl) Max(comparator Comparator, opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &maxOperator{
			comparator: comparator,
			empty:      true,
		}
	}, false, false, opts...)
}

type maxOperator struct {
	comparator Comparator
	empty      bool
	max        interface{}
}

func (op *maxOperator) next(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	op.empty = false

	if op.max == nil {
		op.max = item.V
	} else {
		if op.comparator(op.max, item.V) < 0 {
			op.max = item.V
		}
	}
}

func (op *maxOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *maxOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.empty {
		Of(op.max).SendContext(ctx, dst)
	}
}

func (op *maxOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	op.next(ctx, Of(item.V.(*maxOperator).max), dst, operatorOptions)
}

// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
func (o *ObservableImpl) Min(comparator Comparator, opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &minOperator{
			comparator: comparator,
			empty:      true,
		}
	}, false, false, opts...)
}

type minOperator struct {
	comparator Comparator
	empty      bool
	max        interface{}
}

func (op *minOperator) next(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	op.empty = false

	if op.max == nil {
		op.max = item.V
	} else {
		if op.comparator(op.max, item.V) > 0 {
			op.max = item.V
		}
	}
}

func (op *minOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *minOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.empty {
		Of(op.max).SendContext(ctx, dst)
	}
}

func (op *minOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	op.next(ctx, Of(item.V.(*minOperator).max), dst, operatorOptions)
}

// Observe observes an Observable by returning its channel.
func (o *ObservableImpl) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}

// OnErrorResumeNext instructs an Observable to pass control to another Observable rather than invoking
// onError if it encounters an error.
func (o *ObservableImpl) OnErrorResumeNext(resumeSequence ErrorToObservable, opts ...Option) Observable {
	return observable(o, func() operator {
		return &onErrorResumeNextOperator{resumeSequence: resumeSequence}
	}, true, false, opts...)
}

type onErrorResumeNextOperator struct {
	resumeSequence ErrorToObservable
}

func (op *onErrorResumeNextOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	item.SendContext(ctx, dst)
}

func (op *onErrorResumeNextOperator) err(_ context.Context, item Item, _ chan<- Item, operatorOptions operatorOptions) {
	operatorOptions.resetIterable(op.resumeSequence(item.E))
}

func (op *onErrorResumeNextOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *onErrorResumeNextOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// OnErrorReturn instructs an Observable to emit an item (returned by a specified function)
// rather than invoking onError if it encounters an error.
func (o *ObservableImpl) OnErrorReturn(resumeFunc ErrorFunc, opts ...Option) Observable {
	return observable(o, func() operator {
		return &onErrorReturnOperator{resumeFunc: resumeFunc}
	}, true, false, opts...)
}

type onErrorReturnOperator struct {
	resumeFunc ErrorFunc
}

func (op *onErrorReturnOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	item.SendContext(ctx, dst)
}

func (op *onErrorReturnOperator) err(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	Of(op.resumeFunc(item.E)).SendContext(ctx, dst)
}

func (op *onErrorReturnOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *onErrorReturnOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// OnErrorReturnItem instructs on Observable to emit an item if it encounters an error.
func (o *ObservableImpl) OnErrorReturnItem(resume interface{}, opts ...Option) Observable {
	return observable(o, func() operator {
		return &onErrorReturnItemOperator{resume: resume}
	}, true, false, opts...)
}

type onErrorReturnItemOperator struct {
	resume interface{}
}

func (op *onErrorReturnItemOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	item.SendContext(ctx, dst)
}

func (op *onErrorReturnItemOperator) err(ctx context.Context, _ Item, dst chan<- Item, _ operatorOptions) {
	Of(op.resume).SendContext(ctx, dst)
}

func (op *onErrorReturnItemOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *onErrorReturnItemOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Reduce applies a function to each item emitted by an Observable, sequentially, and emit the final value.
func (o *ObservableImpl) Reduce(apply Func2, opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &reduceOperator{
			apply: apply,
			empty: true,
		}
	}, false, false, opts...)
}

type reduceOperator struct {
	apply Func2
	acc   interface{}
	empty bool
}

func (op *reduceOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	op.empty = false
	v, err := op.apply(ctx, op.acc, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		op.empty = true
		return
	}
	op.acc = v
}

func (op *reduceOperator) err(_ context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	dst <- item
	operatorOptions.stop()
	op.empty = true
}

func (op *reduceOperator) end(ctx context.Context, dst chan<- Item) {
	if !op.empty {
		Of(op.acc).SendContext(ctx, dst)
	}
}

func (op *reduceOperator) gatherNext(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	op.next(ctx, Of(item.V.(*reduceOperator).acc), dst, operatorOptions)
}

// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
// at most count times, at a particular frequency.
// Cannot run in parallel.
func (o *ObservableImpl) Repeat(count int64, frequency Duration, opts ...Option) Observable {
	if count != Infinite {
		if count < 0 {
			return Thrown(IllegalInputError{error: "count must be positive"})
		}
	}

	return observable(o, func() operator {
		return &repeatOperator{
			count:     count,
			frequency: frequency,
			seq:       make([]Item, 0),
		}
	}, true, false, opts...)
}

type repeatOperator struct {
	count     int64
	frequency Duration
	seq       []Item
}

func (op *repeatOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	item.SendContext(ctx, dst)
	op.seq = append(op.seq, item)
}

func (op *repeatOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *repeatOperator) end(ctx context.Context, dst chan<- Item) {
	for {
		select {
		default:
		case <-ctx.Done():
			return
		}
		if op.count != Infinite {
			if op.count == 0 {
				break
			}
		}
		if op.frequency != nil {
			time.Sleep(op.frequency.duration())
		}
		for _, v := range op.seq {
			v.SendContext(ctx, dst)
		}
		op.count = op.count - 1
	}
}

func (op *repeatOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Retry retries if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
// Cannot be run in parallel.
func (o *ObservableImpl) Retry(count int, shouldRetry func(error) bool, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		observe := o.Observe(opts...)
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
					if count < 0 || !shouldRetry(i.E) {
						i.SendContext(ctx, next)
						break loop
					}
					observe = o.Observe(opts...)
				} else {
					i.SendContext(ctx, next)
				}
			}
		}
		close(next)
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Run creates an Observer without consuming the emitted items.
func (o *ObservableImpl) Run(opts ...Option) Disposed {
	dispose := make(chan struct{})
	option := parseOptions(opts...)
	ctx := option.buildContext()

	go func() {
		defer close(dispose)
		observe := o.Observe(opts...)
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
		observe := o.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				i.SendContext(ctx, obsCh)
			}
		}
	}()

	go func() {
		defer close(itCh)
		observe := iterable.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				i.SendContext(ctx, itCh)
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

// Scan apply a Func2 to each item emitted by an Observable, sequentially, and emit each successive value.
// Cannot be run in parallel.
func (o *ObservableImpl) Scan(apply Func2, opts ...Option) Observable {
	return observable(o, func() operator {
		return &scanOperator{
			apply: apply,
		}
	}, true, false, opts...)
}

type scanOperator struct {
	apply   Func2
	current interface{}
}

func (op *scanOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	v, err := op.apply(ctx, op.current, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}
	Of(v).SendContext(ctx, dst)
	op.current = v
}

func (op *scanOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *scanOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *scanOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
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

// Send sends the items to a given channel.
func (o *ObservableImpl) Send(output chan<- Item, opts ...Option) {
	go func() {
		option := parseOptions(opts...)
		ctx := option.buildContext()
		observe := o.Observe(opts...)
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
				i.SendContext(ctx, output)
			}
		}
		close(output)
	}()
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
		observe := o.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				i.SendContext(ctx, obsCh)
			}
		}
	}()

	go func() {
		defer close(itCh)
		observe := iterable.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-observe:
				if !ok {
					return
				}
				i.SendContext(ctx, itCh)
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

		Of(areCorrect && len(mainSequence) == 0 && len(obsSequence) == 0).SendContext(ctx, next)
		close(next)
	}()

	return &SingleImpl{
		iterable: newChannelIterable(next),
	}
}

// Serialize forces an Observable to make serialized calls and to be well-behaved.
func (o *ObservableImpl) Serialize(from int, identifier func(interface{}) int, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()

	ctx := option.buildContext()
	minHeap := binaryheap.NewWith(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})
	counter := int64(from)
	items := make(map[int]interface{})

	go func() {
		src := o.Observe(opts...)
		defer close(next)

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
				minHeap.Push(id)
				items[id] = item.V

				for !minHeap.Empty() {
					v, _ := minHeap.Peek()
					id := v.(int)
					if atomic.LoadInt64(&counter) == int64(id) {
						if itemValue, contains := items[id]; contains {
							minHeap.Pop()
							delete(items, id)
							Of(itemValue).SendContext(ctx, next)
							counter++
							continue
						}
					}
					break
				}
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
// Cannot be run in parallel.
func (o *ObservableImpl) Skip(nth uint, opts ...Option) Observable {
	return observable(o, func() operator {
		return &skipOperator{
			nth: nth,
		}
	}, true, false, opts...)
}

type skipOperator struct {
	nth       uint
	skipCount int
}

func (op *skipOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	if op.skipCount < int(op.nth) {
		op.skipCount++
		return
	}
	item.SendContext(ctx, dst)
}

func (op *skipOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *skipOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *skipOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
// Cannot be run in parallel.
func (o *ObservableImpl) SkipLast(nth uint, opts ...Option) Observable {
	return observable(o, func() operator {
		return &skipLastOperator{
			nth: nth,
		}
	}, true, false, opts...)
}

type skipLastOperator struct {
	nth       uint
	skipCount int
}

func (op *skipLastOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if op.skipCount >= int(op.nth) {
		operatorOptions.stop()
		return
	}
	op.skipCount++
	item.SendContext(ctx, dst)
}

func (op *skipLastOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *skipLastOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *skipLastOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
// Cannot be run in parallel.
func (o *ObservableImpl) SkipWhile(apply Predicate, opts ...Option) Observable {
	return observable(o, func() operator {
		return &skipWhileOperator{
			apply: apply,
			skip:  true,
		}
	}, true, false, opts...)
}

type skipWhileOperator struct {
	apply Predicate
	skip  bool
}

func (op *skipWhileOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	if !op.skip {
		item.SendContext(ctx, dst)
	} else {
		if !op.apply(item.V) {
			op.skip = false
			item.SendContext(ctx, dst)
		}
	}
}

func (op *skipWhileOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *skipWhileOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *skipWhileOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// StartWith emits a specified Iterable before beginning to emit the items from the source Observable.
func (o *ObservableImpl) StartWith(iterable Iterable, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		observe := iterable.Observe(opts...)
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
				i.SendContext(ctx, next)
			}
		}
		observe = o.Observe(opts...)
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
					i.SendContext(ctx, next)
					return
				}
				i.SendContext(ctx, next)
			}
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}

// SumFloat32 calculates the average of float32 emitted by an Observable and emits a float32.
func (o *ObservableImpl) SumFloat32(opts ...Option) OptionalSingle {
	return o.Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		if acc == nil {
			acc = float32(0)
		}
		sum := acc.(float32)
		switch i := elem.(type) {
		default:
			return nil, IllegalInputError{error: fmt.Sprintf("expected type: (float32|int|int8|int16|int32|int64), got: %t", elem)}
		case int:
			return sum + float32(i), nil
		case int8:
			return sum + float32(i), nil
		case int16:
			return sum + float32(i), nil
		case int32:
			return sum + float32(i), nil
		case int64:
			return sum + float32(i), nil
		case float32:
			return sum + i, nil
		}
	}, opts...)
}

// SumFloat64 calculates the average of float64 emitted by an Observable and emits a float64.
func (o *ObservableImpl) SumFloat64(opts ...Option) OptionalSingle {
	return o.Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		if acc == nil {
			acc = float64(0)
		}
		sum := acc.(float64)
		switch i := elem.(type) {
		default:
			return nil, IllegalInputError{error: fmt.Sprintf("expected type: (float32|float64|int|int8|int16|int32|int64), got: %t", elem)}
		case int:
			return sum + float64(i), nil
		case int8:
			return sum + float64(i), nil
		case int16:
			return sum + float64(i), nil
		case int32:
			return sum + float64(i), nil
		case int64:
			return sum + float64(i), nil
		case float32:
			return sum + float64(i), nil
		case float64:
			return sum + i, nil
		}
	}, opts...)
}

// SumInt64 calculates the average of integers emitted by an Observable and emits an int64.
func (o *ObservableImpl) SumInt64(opts ...Option) OptionalSingle {
	return o.Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
		if acc == nil {
			acc = int64(0)
		}
		sum := acc.(int64)
		switch i := elem.(type) {
		default:
			return nil, IllegalInputError{error: fmt.Sprintf("expected type: (int|int8|int16|int32|int64), got: %t", elem)}
		case int:
			return sum + int64(i), nil
		case int8:
			return sum + int64(i), nil
		case int16:
			return sum + int64(i), nil
		case int32:
			return sum + int64(i), nil
		case int64:
			return sum + i, nil
		}
	}, opts...)
}

// Take emits only the first n items emitted by an Observable.
// Cannot be run in parallel.
func (o *ObservableImpl) Take(nth uint, opts ...Option) Observable {
	return observable(o, func() operator {
		return &takeOperator{
			nth: nth,
		}
	}, true, false, opts...)
}

type takeOperator struct {
	nth       uint
	takeCount int
}

func (op *takeOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if op.takeCount >= int(op.nth) {
		operatorOptions.stop()
		return
	}

	op.takeCount++
	item.SendContext(ctx, dst)
}

func (op *takeOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *takeOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *takeOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// TakeLast emits only the last n items emitted by an Observable.
// Cannot be run in parallel.
func (o *ObservableImpl) TakeLast(nth uint, opts ...Option) Observable {
	return observable(o, func() operator {
		n := int(nth)
		return &takeLast{
			n: n,
			r: ring.New(n),
		}
	}, true, false, opts...)
}

type takeLast struct {
	n     int
	r     *ring.Ring
	count int
}

func (op *takeLast) next(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	op.count++
	op.r.Value = item.V
	op.r = op.r.Next()
}

func (op *takeLast) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *takeLast) end(ctx context.Context, dst chan<- Item) {
	if op.count < op.n {
		remaining := op.n - op.count
		if remaining <= op.count {
			op.r = op.r.Move(op.n - op.count)
		} else {
			op.r = op.r.Move(-op.count)
		}
		op.n = op.count
	}
	for i := 0; i < op.n; i++ {
		Of(op.r.Value).SendContext(ctx, dst)
		op.r = op.r.Next()
	}
}

func (op *takeLast) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// TakeUntil returns an Observable that emits items emitted by the source Observable,
// checks the specified predicate for each item, and then completes when the condition is satisfied.
// Cannot be run in parallel.
func (o *ObservableImpl) TakeUntil(apply Predicate, opts ...Option) Observable {
	return observable(o, func() operator {
		return &takeUntilOperator{
			apply: apply,
		}
	}, true, false, opts...)
}

type takeUntilOperator struct {
	apply Predicate
}

func (op *takeUntilOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	item.SendContext(ctx, dst)
	if op.apply(item.V) {
		operatorOptions.stop()
		return
	}
}

func (op *takeUntilOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *takeUntilOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *takeUntilOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// TakeWhile returns an Observable that emits items emitted by the source ObservableSource so long as each
// item satisfied a specified condition, and then completes as soon as this condition is not satisfied.
// Cannot be run in parallel.
func (o *ObservableImpl) TakeWhile(apply Predicate, opts ...Option) Observable {
	return observable(o, func() operator {
		return &takeWhileOperator{
			apply: apply,
		}
	}, true, false, opts...)
}

type takeWhileOperator struct {
	apply Predicate
}

func (op *takeWhileOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if !op.apply(item.V) {
		operatorOptions.stop()
		return
	}
	item.SendContext(ctx, dst)
}

func (op *takeWhileOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *takeWhileOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *takeWhileOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// TimeInterval converts an Observable that emits items into one that emits indications of the amount of time elapsed between those emissions.
func (o *ObservableImpl) TimeInterval(opts ...Option) Observable {
	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		defer close(next)
		observe := o.Observe(opts...)
		latest := time.Now().UTC()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}
				if item.Error() {
					if !item.SendContext(ctx, next) {
						return
					}
					if option.getErrorStrategy() == StopOnError {
						return
					}
				} else {
					now := time.Now().UTC()
					if !Of(now.Sub(latest)).SendContext(ctx, next) {
						return
					}
					latest = now
				}
			}
		}
	}

	return customObservableOperator(f, opts...)
}

// Timestamp attaches a timestamp to each item emitted by an Observable indicating when it was emitted.
func (o *ObservableImpl) Timestamp(opts ...Option) Observable {
	return observable(o, func() operator {
		return &timestampOperator{}
	}, true, false, opts...)
}

type timestampOperator struct {
}

func (op *timestampOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	Of(TimestampItem{
		Timestamp: time.Now().UTC(),
		V:         item.V,
	}).SendContext(ctx, dst)
}

func (op *timestampOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *timestampOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *timestampOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// ToMap convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function.
// Cannot be run in parallel.
func (o *ObservableImpl) ToMap(keySelector Func, opts ...Option) Single {
	return single(o, func() operator {
		return &toMapOperator{
			keySelector: keySelector,
			m:           make(map[interface{}]interface{}),
		}
	}, true, false, opts...)
}

type toMapOperator struct {
	keySelector Func
	m           map[interface{}]interface{}
}

func (op *toMapOperator) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	k, err := op.keySelector(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}
	op.m[k] = item.V
}

func (op *toMapOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *toMapOperator) end(ctx context.Context, dst chan<- Item) {
	Of(op.m).SendContext(ctx, dst)
}

func (op *toMapOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// ToMapWithValueSelector convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function and valued by another
// value function.
// Cannot be run in parallel.
func (o *ObservableImpl) ToMapWithValueSelector(keySelector, valueSelector Func, opts ...Option) Single {
	return single(o, func() operator {
		return &toMapWithValueSelector{
			keySelector:   keySelector,
			valueSelector: valueSelector,
			m:             make(map[interface{}]interface{}),
		}
	}, true, false, opts...)
}

type toMapWithValueSelector struct {
	keySelector, valueSelector Func
	m                          map[interface{}]interface{}
}

func (op *toMapWithValueSelector) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	k, err := op.keySelector(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}

	v, err := op.valueSelector(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}

	op.m[k] = v
}

func (op *toMapWithValueSelector) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *toMapWithValueSelector) end(ctx context.Context, dst chan<- Item) {
	Of(op.m).SendContext(ctx, dst)
}

func (op *toMapWithValueSelector) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// ToSlice collects all items from an Observable and emit them in a slice and an optional error.
// Cannot be run in parallel.
func (o *ObservableImpl) ToSlice(initialCapacity int, opts ...Option) ([]interface{}, error) {
	op := &toSliceOperator{
		s: make([]interface{}, 0, initialCapacity),
	}
	<-observable(o, func() operator {
		return op
	}, true, false, opts...).Run()
	return op.s, op.observableErr
}

type toSliceOperator struct {
	s             []interface{}
	observableErr error
}

func (op *toSliceOperator) next(_ context.Context, item Item, _ chan<- Item, _ operatorOptions) {
	op.s = append(op.s, item.V)
}

func (op *toSliceOperator) err(_ context.Context, item Item, _ chan<- Item, operatorOptions operatorOptions) {
	op.observableErr = item.E
	operatorOptions.stop()
}

func (op *toSliceOperator) end(_ context.Context, _ chan<- Item) {
}

func (op *toSliceOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Unmarshal transforms the items emitted by an Observable by applying an unmarshalling to each item.
func (o *ObservableImpl) Unmarshal(unmarshaller Unmarshaller, factory func() interface{}, opts ...Option) Observable {
	return o.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		v := factory()
		err := unmarshaller(i.([]byte), v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}, opts...)
}

// WindowWithCount periodically subdivides items from an Observable into Observable windows of a given size and emit these windows
// rather than emitting the items one at a time.
func (o *ObservableImpl) WindowWithCount(count int, opts ...Option) Observable {
	if count < 0 {
		return Thrown(IllegalInputError{error: "count must be positive or nil"})
	}

	option := parseOptions(opts...)
	return observable(o, func() operator {
		return &windowWithCountOperator{
			count:  count,
			option: option,
		}
	}, true, false, opts...)
}

type windowWithCountOperator struct {
	count          int
	iCount         int
	currentChannel chan Item
	option         Option
}

func (op *windowWithCountOperator) pre(ctx context.Context, dst chan<- Item) {
	if op.currentChannel == nil {
		ch := op.option.buildChannel()
		op.currentChannel = ch
		Of(FromChannel(ch)).SendContext(ctx, dst)
	}
}

func (op *windowWithCountOperator) post(ctx context.Context, dst chan<- Item) {
	if op.iCount == op.count {
		op.iCount = 0
		close(op.currentChannel)
		ch := op.option.buildChannel()
		op.currentChannel = ch
		Of(FromChannel(ch)).SendContext(ctx, dst)
	}
}

func (op *windowWithCountOperator) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	op.pre(ctx, dst)
	op.currentChannel <- item
	op.iCount++
	op.post(ctx, dst)
}

func (op *windowWithCountOperator) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	op.pre(ctx, dst)
	op.currentChannel <- item
	op.iCount++
	op.post(ctx, dst)
	operatorOptions.stop()
}

func (op *windowWithCountOperator) end(_ context.Context, _ chan<- Item) {
	if op.currentChannel != nil {
		close(op.currentChannel)
	}
}

func (op *windowWithCountOperator) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// WindowWithTime periodically subdivides items from an Observable into Observables based on timed windows
// and emit them rather than emitting the items one at a time.
func (o *ObservableImpl) WindowWithTime(timespan Duration, opts ...Option) Observable {
	if timespan == nil {
		return Thrown(IllegalInputError{error: "timespan must no be nil"})
	}

	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		observe := o.Observe(opts...)
		ch := option.buildChannel()
		done := make(chan struct{})
		empty := true
		mutex := sync.Mutex{}
		if !Of(FromChannel(ch)).SendContext(ctx, next) {
			return
		}

		go func() {
			defer func() {
				mutex.Lock()
				close(ch)
				mutex.Unlock()
			}()
			defer close(next)
			for {
				select {
				case <-ctx.Done():
					return
				case <-done:
					return
				case <-time.After(timespan.duration()):
					mutex.Lock()
					if empty {
						mutex.Unlock()
						continue
					}
					close(ch)
					empty = true
					ch = option.buildChannel()
					if !Of(FromChannel(ch)).SendContext(ctx, next) {
						close(done)
						return
					}
					mutex.Unlock()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case item, ok := <-observe:
				if !ok {
					close(done)
					return
				}
				if item.Error() {
					mutex.Lock()
					if !item.SendContext(ctx, ch) {
						mutex.Unlock()
						close(done)
						return
					}
					mutex.Unlock()
					if option.getErrorStrategy() == StopOnError {
						close(done)
						return
					}
				}
				mutex.Lock()
				if !item.SendContext(ctx, ch) {
					mutex.Unlock()
					return
				}
				empty = false
				mutex.Unlock()
			}
		}
	}

	return customObservableOperator(f, opts...)
}

// WindowWithTimeOrCount periodically subdivides items from an Observable into Observables based on timed windows or a specific size
// and emit them rather than emitting the items one at a time.
func (o *ObservableImpl) WindowWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable {
	if timespan == nil {
		return Thrown(IllegalInputError{error: "timespan must no be nil"})
	}
	if count < 0 {
		return Thrown(IllegalInputError{error: "count must be positive or nil"})
	}

	f := func(ctx context.Context, next chan Item, option Option, opts ...Option) {
		observe := o.Observe(opts...)
		ch := option.buildChannel()
		done := make(chan struct{})
		mutex := sync.Mutex{}
		iCount := 0
		if !Of(FromChannel(ch)).SendContext(ctx, next) {
			return
		}

		go func() {
			defer func() {
				mutex.Lock()
				close(ch)
				mutex.Unlock()
			}()
			defer close(next)
			for {
				select {
				case <-ctx.Done():
					return
				case <-done:
					return
				case <-time.After(timespan.duration()):
					mutex.Lock()
					if iCount == 0 {
						mutex.Unlock()
						continue
					}
					close(ch)
					iCount = 0
					ch = option.buildChannel()
					if !Of(FromChannel(ch)).SendContext(ctx, next) {
						close(done)
						return
					}
					mutex.Unlock()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case item, ok := <-observe:
				if !ok {
					close(done)
					return
				}
				if item.Error() {
					mutex.Lock()
					if !item.SendContext(ctx, ch) {
						mutex.Unlock()
						close(done)
						return
					}
					mutex.Unlock()
					if option.getErrorStrategy() == StopOnError {
						close(done)
						return
					}
				}
				mutex.Lock()
				if !item.SendContext(ctx, ch) {
					mutex.Unlock()
					return
				}
				iCount++
				if iCount == count {
					close(ch)
					iCount = 0
					ch = option.buildChannel()
					if !Of(FromChannel(ch)).SendContext(ctx, next) {
						mutex.Unlock()
						close(done)
						return
					}
				}
				mutex.Unlock()
			}
		}
	}

	return customObservableOperator(f, opts...)
}

// ZipFromIterable merges the emissions of an Iterable via a specified function
// and emit single items for each combination based on the results of this function.
func (o *ObservableImpl) ZipFromIterable(iterable Iterable, zipper Func2, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		defer close(next)
		it1 := o.Observe(opts...)
		it2 := iterable.Observe(opts...)
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
					i1.SendContext(ctx, next)
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
							i2.SendContext(ctx, next)
							return
						}
						v, err := zipper(ctx, i1.V, i2.V)
						if err != nil {
							Error(err).SendContext(ctx, next)
							return
						}
						Of(v).SendContext(ctx, next)
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
