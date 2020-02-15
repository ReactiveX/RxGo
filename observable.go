package rxgo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Observable is the basic observable interface.
type Observable interface {
	Iterable
	All(ctx context.Context, predicate Predicate) Single
	AverageFloat32(ctx context.Context) Single
	AverageFloat64(ctx context.Context) Single
	AverageInt(ctx context.Context) Single
	AverageInt8(ctx context.Context) Single
	AverageInt16(ctx context.Context) Single
	AverageInt32(ctx context.Context) Single
	AverageInt64(ctx context.Context) Single
	BufferWithCount(ctx context.Context, count, skip int) Observable
	BufferWithTime(ctx context.Context, timespan, timeshift Duration) Observable
	Filter(ctx context.Context, apply Predicate) Observable
	ForEach(ctx context.Context, nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc)
	Map(ctx context.Context, apply Func) Observable
	// TODO Add backoff retry
	Retry(ctx context.Context, count int) Observable
	SkipWhile(ctx context.Context, apply Predicate) Observable
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

func defaultEndFuncOperator(_ Item, _ chan<- Item, _ func()) {}

func operator(ctx context.Context, iterable Iterable, nextFunc, errFunc, endFunc Operator) chan Item {
	next := make(chan Item)

	stopped := false
	stop := func() {
		stopped = true
	}

	go func() {
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
		endFunc(FromValue(nil), next, nil)
		close(next)
	}()

	return next
}

// TODO Options
func newObservableFromOperator(ctx context.Context, iterable Iterable, nextFunc, errFunc, endFunc Operator) Observable {
	next := operator(ctx, iterable, nextFunc, errFunc, endFunc)
	return &observable{
		iterable: newChannelIterable(next),
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

func (o *observable) All(ctx context.Context, predicate Predicate) Single {
	all := true
	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if !predicate(item.Value) {
			dst <- FromValue(false)
			all = false
			stop()
		}
	}, defaultErrorFuncOperator, func(item Item, dst chan<- Item, stop func()) {
		if all {
			dst <- FromValue(true)
		}
	})
}

func (o *observable) AverageFloat32(ctx context.Context) Single {
	var sum float32
	var count float32

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(float32); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: float32, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) AverageFloat64(ctx context.Context) Single {
	var sum float64
	var count float64

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(float64); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: float64, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) AverageInt(ctx context.Context) Single {
	var sum int
	var count int

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) AverageInt8(ctx context.Context) Single {
	var sum int8
	var count int8

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int8); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int8, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) AverageInt16(ctx context.Context) Single {
	var sum int16
	var count int16

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int16); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int16, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) AverageInt32(ctx context.Context) Single {
	var sum int32
	var count int32

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int32); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int32, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) AverageInt64(ctx context.Context) Single {
	var sum int64
	var count int64

	return newSingleFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if v, ok := item.Value.(int64); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int64, got: %t", item)))
			stop()
		}
	}, defaultErrorFuncOperator, func(_ Item, dst chan<- Item, _ func()) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	})
}

func (o *observable) BufferWithCount(ctx context.Context, count, skip int) Observable {
	if count <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "count must be positive"))
	}
	if skip <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "skip must be positive"))
	}

	buffer := make([]interface{}, count)
	iCount := 0
	iSkip := 0

	return newObservableFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
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
		stop()
	}, func(_ Item, dst chan<- Item, _ func()) {
		if iCount != 0 {
			dst <- FromValue(buffer[:iCount])
		}
	})
}

func (o *observable) BufferWithTime(ctx context.Context, timespan, timeshift Duration) Observable {
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

	next := make(chan Item)

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

// TODO Options?
func (o *observable) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe()
}

func (o *observable) Filter(ctx context.Context, apply Predicate) Observable {
	return newObservableFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
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
					return
				}
				nextFunc(i.Value)
			}
		}
	}
	newObservableFromHandler(ctx, o, handler)
}

func (o *observable) Map(ctx context.Context, apply Func) Observable {
	return newObservableFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			stop()
		}
		dst <- FromValue(res)
	}, defaultErrorFuncOperator, defaultEndFuncOperator)
}

func (o *observable) Retry(ctx context.Context, count int) Observable {
	next := make(chan Item)

	stopped := false
	stop := func() {
		stopped = true
	}

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
					count--
					if count < 0 {
						next <- i
						stop()
						return
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

func (o *observable) SkipWhile(ctx context.Context, apply Predicate) Observable {
	skip := true

	return newObservableFromOperator(ctx, o, func(item Item, dst chan<- Item, stop func()) {
		if !skip {
			dst <- item
		} else {
			if !apply(item.Value) {
				skip = false
				dst <- item
			}
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator)
}
