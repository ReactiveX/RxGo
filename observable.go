package rxgo

import (
	"context"
	"sync"

	"github.com/cenkalti/backoff/v4"
)

// Observable is the standard interface for Observables.
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
	BackOffRetry(backOffCfg backoff.BackOff, opts ...Option) Observable
	BufferWithCount(count int, opts ...Option) Observable
	BufferWithTime(timespan, timeshift Duration, opts ...Option) Observable
	BufferWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable
	Contains(equal Predicate, opts ...Option) Single
	Count(opts ...Option) Single
	DefaultIfEmpty(defaultValue interface{}, opts ...Option) Observable
	Distinct(apply Func, opts ...Option) Observable
	DistinctUntilChanged(apply Func, opts ...Option) Observable
	DoOnCompleted(completedFunc CompletedFunc, opts ...Option) Disposed
	DoOnError(errFunc ErrFunc, opts ...Option) Disposed
	DoOnNext(nextFunc NextFunc, opts ...Option) Disposed
	ElementAt(index uint, opts ...Option) Single
	Error(opts ...Option) error
	Errors(opts ...Option) []error
	Filter(apply Predicate, opts ...Option) Observable
	First(opts ...Option) OptionalSingle
	FirstOrDefault(defaultValue interface{}, opts ...Option) Single
	FlatMap(apply ItemToObservable, opts ...Option) Observable
	ForEach(nextFunc NextFunc, errFunc ErrFunc, completedFunc CompletedFunc, opts ...Option) Disposed
	GroupBy(length int, distribution func(Item) int, opts ...Option) Observable
	IgnoreElements(opts ...Option) Observable
	Last(opts ...Option) OptionalSingle
	LastOrDefault(defaultValue interface{}, opts ...Option) Single
	Map(apply Func, opts ...Option) Observable
	Marshal(marshaller Marshaller, opts ...Option) Observable
	Max(comparator Comparator, opts ...Option) OptionalSingle
	Min(comparator Comparator, opts ...Option) OptionalSingle
	OnErrorResumeNext(resumeSequence ErrorToObservable, opts ...Option) Observable
	OnErrorReturn(resumeFunc ErrorFunc, opts ...Option) Observable
	OnErrorReturnItem(resume interface{}, opts ...Option) Observable
	Reduce(apply Func2, opts ...Option) OptionalSingle
	Repeat(count int64, frequency Duration, opts ...Option) Observable
	Retry(count int, opts ...Option) Observable
	Run(opts ...Option) Disposed
	Sample(iterable Iterable, opts ...Option) Observable
	Scan(apply Func2, opts ...Option) Observable
	SequenceEqual(iterable Iterable, opts ...Option) Single
	Send(output chan<- Item, opts ...Option)
	Serialize(from int, identifier func(interface{}) int, opts ...Option) Observable
	Skip(nth uint, opts ...Option) Observable
	SkipLast(nth uint, opts ...Option) Observable
	SkipWhile(apply Predicate, opts ...Option) Observable
	StartWithIterable(iterable Iterable, opts ...Option) Observable
	SumFloat32(opts ...Option) OptionalSingle
	SumFloat64(opts ...Option) OptionalSingle
	SumInt64(opts ...Option) OptionalSingle
	Take(nth uint, opts ...Option) Observable
	TakeLast(nth uint, opts ...Option) Observable
	TakeUntil(apply Predicate, opts ...Option) Observable
	TakeWhile(apply Predicate, opts ...Option) Observable
	ToMap(keySelector Func, opts ...Option) Single
	ToMapWithValueSelector(keySelector, valueSelector Func, opts ...Option) Single
	ToSlice(initialCapacity int, opts ...Option) ([]interface{}, error)
	Unmarshal(unmarshaller Unmarshaller, factory func() interface{}, opts ...Option) Observable
	WindowWithCount(count int, opts ...Option) Observable
	ZipFromIterable(iterable Iterable, zipper Func2, opts ...Option) Observable
}

// ObservableImpl implements Observable.
type ObservableImpl struct {
	iterable Iterable
}

func defaultErrorFuncOperator(_ context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	dst <- item
	operatorOptions.stop()
}

type operator interface {
	next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions)
	err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions)
	end(ctx context.Context, dst chan<- Item)
	gatherNext(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions)
}

func observable(iterable Iterable, operatorFactory func() operator, forceSeq, bypassGather bool, opts ...Option) Observable {
	option := parseOptions(opts...)
	parallel, _ := option.getPool()

	if option.isEagerObservation() {
		next := option.buildChannel()
		ctx := option.buildContext()
		if forceSeq || !parallel {
			runSeq(ctx, next, iterable, operatorFactory, option, opts...)
		} else {
			runPar(ctx, next, iterable, operatorFactory, bypassGather, option, opts...)
		}
		return &ObservableImpl{iterable: newChannelIterable(next)}
	}

	return &ObservableImpl{
		iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
			mergedOptions := append(opts, propagatedOptions...)
			option = parseOptions(mergedOptions...)

			next := option.buildChannel()
			ctx := option.buildContext()
			if forceSeq || !parallel {
				runSeq(ctx, next, iterable, operatorFactory, option, mergedOptions...)
			} else {
				runPar(ctx, next, iterable, operatorFactory, bypassGather, option, mergedOptions...)
			}
			return next
		}),
	}
}

func single(iterable Iterable, operatorFactory func() operator, forceSeq, bypassGather bool, opts ...Option) Single {
	option := parseOptions(opts...)
	parallel, _ := option.getPool()

	if option.isEagerObservation() {
		next := option.buildChannel()
		ctx := option.buildContext()
		if forceSeq || !parallel {
			runSeq(ctx, next, iterable, operatorFactory, option, opts...)
		} else {
			runPar(ctx, next, iterable, operatorFactory, bypassGather, option, opts...)
		}
		return &SingleImpl{iterable: newChannelIterable(next)}
	}

	return &SingleImpl{
		iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
			mergedOptions := append(opts, propagatedOptions...)
			option = parseOptions(mergedOptions...)

			next := option.buildChannel()
			ctx := option.buildContext()
			if forceSeq || !parallel {
				runSeq(ctx, next, iterable, operatorFactory, option, mergedOptions...)
			} else {
				runPar(ctx, next, iterable, operatorFactory, bypassGather, option, mergedOptions...)
			}
			return next
		}),
	}
}

func optionalSingle(iterable Iterable, operatorFactory func() operator, forceSeq, bypassGather bool, opts ...Option) OptionalSingle {
	option := parseOptions(opts...)
	parallel, _ := option.getPool()

	if option.isEagerObservation() {
		next := option.buildChannel()
		ctx := option.buildContext()
		if forceSeq || !parallel {
			runSeq(ctx, next, iterable, operatorFactory, option, opts...)
		} else {
			runPar(ctx, next, iterable, operatorFactory, bypassGather, option, opts...)
		}
		return &OptionalSingleImpl{iterable: newChannelIterable(next)}
	}

	return &OptionalSingleImpl{
		iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
			mergedOptions := append(opts, propagatedOptions...)
			option = parseOptions(mergedOptions...)

			next := option.buildChannel()
			ctx := option.buildContext()
			if forceSeq || !parallel {
				runSeq(ctx, next, iterable, operatorFactory, option, mergedOptions...)
			} else {
				runPar(ctx, next, iterable, operatorFactory, bypassGather, option, mergedOptions...)
			}
			return next
		}),
	}
}

func runSeq(ctx context.Context, next chan Item, iterable Iterable, operatorFactory func() operator, option Option, opts ...Option) {
	go func() {
		op := operatorFactory()
		stopped := false
		observe := iterable.Observe(opts...)
		operator := operatorOptions{
			stop: func() {
				if option.getErrorStrategy() == Stop {
					stopped = true
				}
			},
			resetIterable: func(newIterable Iterable) {
				observe = newIterable.Observe(opts...)
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
				if i.Error() {
					op.err(ctx, i, next, operator)
				} else {
					op.next(ctx, i, next, operator)
				}
			}
		}
		op.end(ctx, next)
		close(next)
	}()
}

func runPar(ctx context.Context, next chan Item, iterable Iterable, operatorFactory func() operator, bypassGather bool, option Option, opts ...Option) {
	observe := iterable.Observe(opts...)

	wg := sync.WaitGroup{}
	_, pool := option.getPool()
	wg.Add(pool)
	ctx, cancel := context.WithCancel(ctx)

	var gather chan Item
	if bypassGather {
		gather = next
	} else {
		gather = make(chan Item, 1)

		go func() {
			op := operatorFactory()
			stopped := false
			operator := operatorOptions{
				stop: func() {
					if option.getErrorStrategy() == Stop {
						stopped = true
					}
				},
				resetIterable: func(newIterable Iterable) {
					observe = newIterable.Observe(opts...)
				},
			}
			for item := range gather {
				if stopped {
					break
				}
				if item.Error() {
					op.err(ctx, item, next, operator)
				} else {
					op.gatherNext(ctx, item, next, operator)
				}
			}
			op.end(ctx, next)
			close(next)
		}()
	}

	for i := 0; i < pool; i++ {
		go func() {
			op := operatorFactory()
			stopped := false
			operator := operatorOptions{
				stop: func() {
					if option.getErrorStrategy() == Stop {
						stopped = true
					}
				},
				resetIterable: func(newIterable Iterable) {
					observe = newIterable.Observe(opts...)
				},
			}
			defer wg.Done()
			for !stopped {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-observe:
					if !ok {
						if !bypassGather {
							gather <- Of(op)
						}
						return
					}
					if item.Error() {
						op.err(ctx, item, gather, operator)
					} else {
						op.next(ctx, item, gather, operator)
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		cancel()
		close(gather)
	}()
}
