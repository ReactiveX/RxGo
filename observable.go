package rxgo

import (
	"context"
	"sync"

	"github.com/tevino/abool"
)

// Observable is the basic observable interface.
// TODO Throttling
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
	BufferWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable
	BufferWithTime(timespan, timeshift Duration, opts ...Option) Observable
	Contains(equal Predicate, opts ...Option) Single
	Count(opts ...Option) Single
	DefaultIfEmpty(defaultValue interface{}, opts ...Option) Observable
	Distinct(apply Func, opts ...Option) Observable
	DistinctUntilChanged(apply Func, opts ...Option) Observable
	ElementAt(index uint, opts ...Option) Single
	Filter(apply Predicate, opts ...Option) Observable
	First(opts ...Option) OptionalSingle
	FirstOrDefault(defaultValue interface{}, opts ...Option) Single
	ForEach(nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc, opts ...Option)
	IgnoreElements(opts ...Option) Observable
	Last(opts ...Option) OptionalSingle
	LastOrDefault(defaultValue interface{}, opts ...Option) Single
	Map(apply Func, opts ...Option) Observable
	Marshal(marshaler Marshaler, opts ...Option) Observable
	Max(comparator Comparator, opts ...Option) OptionalSingle
	Min(comparator Comparator, opts ...Option) OptionalSingle
	// TODO Add backoff retry
	OnErrorResumeNext(resumeSequence ErrorToObservable, opts ...Option) Observable
	OnErrorReturn(resumeFunc ErrorFunc, opts ...Option) Observable
	OnErrorReturnItem(resume interface{}, opts ...Option) Observable
	Reduce(apply Func2, opts ...Option) OptionalSingle
	Repeat(count int64, frequency Duration, opts ...Option) Observable
	Retry(count int, opts ...Option) Observable
	Sample(obs Observable, opts ...Option) Observable
	Scan(apply Func2, opts ...Option) Observable
	SequenceEqual(obs Observable, opts ...Option) Single
	Send(output chan<- Item, opts ...Option)
	Skip(nth uint, opts ...Option) Observable
	SkipLast(nth uint, opts ...Option) Observable
	SkipWhile(apply Predicate, opts ...Option) Observable
	StartWithIterable(iterable Iterable, opts ...Option) Observable
	StartWithObservable(observable Observable, opts ...Option) Observable
	SumFloat32(opts ...Option) Single
	SumFloat64(opts ...Option) Single
	SumInt64(opts ...Option) Single
	Take(nth uint, opts ...Option) Observable
	TakeLast(nth uint, opts ...Option) Observable
	TakeUntil(apply Predicate, opts ...Option) Observable
	TakeWhile(apply Predicate, opts ...Option) Observable
	Timeout(opts ...Option) Observable
	ToMap(keySelector Func, opts ...Option) Single
	ToMapWithValueSelector(keySelector, valueSelector Func, opts ...Option) Single
	ToSlice(opts ...Option) Single
	Unmarshal(unmarshaler Unmarshaler, factory func() interface{}, opts ...Option) Observable
	ZipFromObservable(publisher Observable, zipper Func2, opts ...Option) Observable
}

type observable struct {
	iterable Iterable
}

func defaultNextFuncOperator(_ context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	dst <- item
}

func defaultErrorFuncOperator(_ context.Context, item Item, dst chan<- Item, operator operatorOptions) {
	dst <- item
	operator.stop()
}

func defaultEndFuncOperator(_ context.Context, _ chan<- Item) {}

func operator(iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd, opts ...Option) Iterable {
	option := parseOptions(opts...)

	if option.withEagerObservation() {
		next := option.buildChannel()
		ctx := option.buildContext()
		if withPool, pool := option.withPool(); withPool {
			parallel(ctx, pool, next, iterable, nextFunc, errFunc, endFunc)
		} else {
			seq(ctx, next, iterable, nextFunc, errFunc, endFunc)
		}

		return newChannelIterable(next)
	}

	return &observable{
		iterable: newColdIterable(func() <-chan Item {
			next := option.buildChannel()
			ctx := option.buildContext()
			if withPool, pool := option.withPool(); withPool {
				parallel(ctx, pool, next, iterable, nextFunc, errFunc, endFunc)
			} else {
				seq(ctx, next, iterable, nextFunc, errFunc, endFunc)
			}
			return next
		}),
	}
}

func seq(ctx context.Context, next chan Item, iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd) {
	go func() {
		stopped := false
		observe := iterable.Observe()
		operator := operatorOptions{
			stop: func() {
				stopped = true
			},
			resetIterable: func(newIterable Iterable) {
				observe = newIterable.Observe()
			},
		}

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
					errFunc(ctx, i, next, operator)
				} else {
					nextFunc(ctx, i, next, operator)
				}
			}
		}
		endFunc(ctx, next)
		close(next)
	}()
}

func parallel(ctx context.Context, pool int, next chan Item, iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd) {
	stopped := abool.New()
	observe := iterable.Observe()
	operator := operatorOptions{
		stop: func() {
			stopped.Set()
		},
		// TODO Can we implement a reset strategy with a parallel implementation
		resetIterable: func(_ Iterable) {},
	}

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
						errFunc(ctx, i, next, operator)
					} else {
						nextFunc(ctx, i, next, operator)
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		endFunc(ctx, next)
		close(next)
	}()
}

func newObservableFromOperator(iterable Iterable, nextFunc, errFunc operatorItem, endFunc operatorEnd, opts ...Option) Observable {
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
