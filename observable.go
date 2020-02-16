package rxgo

import (
	"container/ring"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tevino/abool"

	"github.com/pkg/errors"
)

// Observable is the basic observable interface.
type Observable interface {
	Iterable
	All(predicate Predicate, opts ...Option) Single
	AverageFloat32(opts ...Option) Single
	AverageFloat64(opts ...Option) Single
	AverageInt(opts ...Option) Single
	AverageInt8(opts ...Option) Single
	AverageInt16(opts ...Option) Single
	AverageInt32(opts ...Option) Single
	AverageInt64(opts ...Option) Single
	BufferWithCount(count, skip int, opts ...Option) Observable
	BufferWithTime(timespan, timeshift Duration, opts ...Option) Observable
	Contains(equal Predicate, opts ...Option) Single
	Count(opts ...Option) Single
	Filter(apply Predicate, opts ...Option) Observable
	ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc)
	// TODO With pool
	Map(apply Func, opts ...Option) Observable
	Marshal(marshaler Marshaler, opts ...Option) Observable
	// TODO Add backoff retry
	Retry(count int, opts ...Option) Observable
	SkipWhile(apply Predicate, opts ...Option) Observable
	Take(nth uint, opts ...Option) Observable
	TakeLast(nth uint, opts ...Option) Observable
	ToSlice(opts ...Option) Single
	// TODO Throttling
	Unmarshal(unmarshaler Unmarshaler, factory func() interface{}, opts ...Option) Observable
}

type observable struct {
	iterable Iterable
}

func newObservableFromHandler(ctx context.Context, source Observable, handler Iterator) Observable {
	next := make(chan Item)

	go handler(ctx, source.Observe(), next)

	return &observable{
		iterable: newChannelIterable(next),
	}
}

func defaultErrorFuncOperator(item Item, dst chan<- Item, stop func()) {
	dst <- item
	stop()
}

func defaultEndFuncOperator(_ chan<- Item) {}

func operator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) Iterable {
	next, ctx, option := buildOptionValues(opts...)

	if option.withEagerObservation() {
		if withPool, pool := option.withPool(); withPool {
			parallel(ctx, pool, next, iterable, nextFunc, errFunc, endFunc)
		} else {
			seq(ctx, next, iterable, nextFunc, errFunc, endFunc)
		}

		return newChannelIterable(next)
	}

	return &observable{
		iterable: newColdIterable(func() <-chan Item {
			next, ctx, option := buildOptionValues(opts...)
			if withPool, pool := option.withPool(); withPool {
				parallel(ctx, pool, next, iterable, nextFunc, errFunc, endFunc)
			} else {
				seq(ctx, next, iterable, nextFunc, errFunc, endFunc)
			}
			return next
		}),
	}
}

func seq(ctx context.Context, next chan Item, iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler) {
	go func() {
		stopped := false
		stop := func() {
			stopped = true
		}

		observe := iterable.Observe()
	loop:
		for !stopped {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.IsError() {
					errFunc(i, next, stop)
				} else {
					nextFunc(i, next, stop)
				}
			}
		}
		endFunc(next)
		close(next)
	}()
}

func parallel(ctx context.Context, pool int, next chan Item, iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler) {
	stopped := abool.New()
	stop := func() {
		stopped.Set()
	}

	observe := iterable.Observe()

	wg := sync.WaitGroup{}
	wg.Add(pool)
	for i := 0; i < pool; i++ {
		go func() {
			for !stopped.IsSet() {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case i, ok := <-observe:
					if !ok {
						wg.Done()
						return
					}
					if i.IsError() {
						errFunc(i, next, stop)
					} else {
						nextFunc(i, next, stop)
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		endFunc(next)
		close(next)
	}()
}

func newObservableFromOperator(iterable Iterable, nextFunc, errFunc ItemHandler, endFunc EndHandler, opts ...Option) Observable {
	return &observable{
		iterable: operator(iterable, nextFunc, errFunc, endFunc, opts...),
	}
}

func newObservableFromError(err error) Observable {
	next := make(chan Item, 1)
	next <- FromError(err)
	close(next)
	return &observable{
		iterable: newChannelIterable(next),
	}
}

func (o *observable) All(predicate Predicate, opts ...Option) Single {
	all := true
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if !predicate(item.Value) {
			dst <- FromValue(false)
			all = false
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if all {
			dst <- FromValue(true)
		}
	}, opts...)
}

func (o *observable) AverageFloat32(opts ...Option) Single {
	var sum float32
	var count float32

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(float32); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: float32, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) AverageFloat64(opts ...Option) Single {
	var sum float64
	var count float64

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(float64); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: float64, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) AverageInt(opts ...Option) Single {
	var sum int
	var count int

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) AverageInt8(opts ...Option) Single {
	var sum int8
	var count int8

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int8); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int8, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) AverageInt16(opts ...Option) Single {
	var sum int16
	var count int16

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int16); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int16, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) AverageInt32(opts ...Option) Single {
	var sum int32
	var count int32

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int32); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int32, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) AverageInt64(opts ...Option) Single {
	var sum int64
	var count int64

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int64); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int64, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

func (o *observable) BufferWithCount(count, skip int, opts ...Option) Observable {
	if count <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "count must be positive"))
	}
	if skip <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "skip must be positive"))
	}

	buffer := make([]interface{}, count)
	iCount := 0
	iSkip := 0

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if iCount >= count {
			// Skip
			iSkip++
		} else {
			// Add to buffer
			buffer[iCount] = item.Value
			iCount++
			iSkip++
		}
		if iSkip == skip {
			// Send current buffer
			dst <- FromValue(buffer)
			buffer = make([]interface{}, count)
			iCount = 0
			iSkip = 0
		}
	}, func(item Item, dst chan<- Item, stop func()) {
		if iCount != 0 {
			dst <- FromValue(buffer[:iCount])
		}
		dst <- item
		iCount = 0
		stop()
	}, func(dst chan<- Item) {
		if iCount != 0 {
			dst <- FromValue(buffer[:iCount])
		}
	}, opts...)
}

func (o *observable) BufferWithTime(timespan, timeshift Duration, opts ...Option) Observable {
	if timespan == nil || timespan.duration() == 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "timespan must no be nil"))
	}

	if timeshift == nil {
		timeshift = WithDuration(0)
	}

	var mux sync.Mutex
	var listenMutex sync.Mutex
	buffer := make([]interface{}, 0)
	stop := false
	listen := true

	next, ctx, _ := buildOptionValues(opts...)

	stopped := false

	// First goroutine in charge to check the timespan
	go func() {
		for {
			time.Sleep(timespan.duration())
			mux.Lock()
			if !stop {
				next <- FromValue(buffer)
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
				if i.IsError() {
					mux.Lock()
					if len(buffer) > 0 {
						next <- FromValue(buffer)
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
					buffer = append(buffer, i.Value)
				}
				mux.Unlock()
			}
		}
		mux.Lock()
		if len(buffer) > 0 {
			next <- FromValue(buffer)
		}
		close(next)
		stop = true
		mux.Unlock()
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

func (o *observable) Contains(equal Predicate, opts ...Option) Single {
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if equal(item.Value) {
			dst <- FromValue(true)
			stop()
			return
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		dst <- FromValue(false)
	}, opts...)
}

func (o *observable) Count(opts ...Option) Single {
	var count int64
	return newSingleFromOperator(o, func(_ Item, dst chan<- Item, _ func()) {
		count++
	}, func(_ Item, dst chan<- Item, stop func()) {
		count++
		dst <- FromValue(count)
		stop()
	}, defaultEndFuncOperator, opts...)
}

func (o *observable) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}

func (o *observable) Filter(apply Predicate, opts ...Option) Observable {
	return newObservableFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if apply(item.Value) {
			dst <- item
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator)
}

func (o *observable) ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc) {
	handler := func(ctx context.Context, src <-chan Item, dst chan<- Item) {
		for {
			select {
			case <-ctx.Done():
				doneFunc()
				return
			case i, ok := <-src:
				if !ok {
					doneFunc()
					return
				}
				if i.IsError() {
					errFunc(i.Err)
					break
				}
				nextFunc(i.Value)
			}
		}
	}
	newObservableFromHandler(ctx, o, handler)
}

func (o *observable) Map(apply Func, opts ...Option) Observable {
	return newObservableFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			stop()
		}
		dst <- FromValue(res)
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

func (o *observable) Marshal(marshaler Marshaler, opts ...Option) Observable {
	return o.Map(func(i interface{}) (interface{}, error) {
		return marshaler(i)
	}, opts...)
}

func (o *observable) Retry(count int, opts ...Option) Observable {
	next, ctx, _ := buildOptionValues(opts...)

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
				if i.IsError() {
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

	return &observable{
		iterable: newChannelIterable(next),
	}
}

func (o *observable) SkipWhile(apply Predicate, opts ...Option) Observable {
	skip := true

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if !skip {
			dst <- item
		} else {
			if !apply(item.Value) {
				skip = false
				dst <- item
			}
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

func (o *observable) Take(nth uint, opts ...Option) Observable {
	takeCount := 0

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		if takeCount < int(nth) {
			takeCount++
			dst <- item
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

func (o *observable) TakeLast(nth uint, opts ...Option) Observable {
	n := int(nth)
	r := ring.New(n)
	count := 0

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		count++
		r.Value = item.Value
		r = r.Next()
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
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
			dst <- FromValue(r.Value)
			r = r.Next()
		}
	}, opts...)
}

func (o *observable) ToSlice(opts ...Option) Single {
	s := make([]interface{}, 0)
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, stop func()) {
		s = append(s, item.Value)
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		dst <- FromValue(s)
	}, opts...)
}

func (o *observable) Unmarshal(unmarshaler Unmarshaler, factory func() interface{}, opts ...Option) Observable {
	return o.Map(func(i interface{}) (interface{}, error) {
		v := factory()
		err := unmarshaler(i.([]byte), v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}, opts...)
}
