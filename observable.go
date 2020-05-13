package rxgo

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/emirpasic/gods/trees/binaryheap"
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
	BufferWithTime(timespan Duration, opts ...Option) Observable
	BufferWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable
	Connect() (context.Context, Disposable)
	Contains(equal Predicate, opts ...Option) Single
	Count(opts ...Option) Single
	Debounce(timespan Duration, opts ...Option) Observable
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
	Join(joiner Func2, right Observable, timeExtractor func(interface{}) time.Time, window Duration, opts ...Option) Observable
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
	Retry(count int, shouldRetry func(error) bool, opts ...Option) Observable
	Run(opts ...Option) Disposed
	Sample(iterable Iterable, opts ...Option) Observable
	Scan(apply Func2, opts ...Option) Observable
	SequenceEqual(iterable Iterable, opts ...Option) Single
	Send(output chan<- Item, opts ...Option)
	Serialize(from int, identifier func(interface{}) int, opts ...Option) Observable
	Skip(nth uint, opts ...Option) Observable
	SkipLast(nth uint, opts ...Option) Observable
	SkipWhile(apply Predicate, opts ...Option) Observable
	StartWith(iterable Iterable, opts ...Option) Observable
	SumFloat32(opts ...Option) OptionalSingle
	SumFloat64(opts ...Option) OptionalSingle
	SumInt64(opts ...Option) OptionalSingle
	Take(nth uint, opts ...Option) Observable
	TakeLast(nth uint, opts ...Option) Observable
	TakeUntil(apply Predicate, opts ...Option) Observable
	TakeWhile(apply Predicate, opts ...Option) Observable
	TimeInterval(opts ...Option) Observable
	Timestamp(opts ...Option) Observable
	ToMap(keySelector Func, opts ...Option) Single
	ToMapWithValueSelector(keySelector, valueSelector Func, opts ...Option) Single
	ToSlice(initialCapacity int, opts ...Option) ([]interface{}, error)
	Unmarshal(unmarshaller Unmarshaller, factory func() interface{}, opts ...Option) Observable
	WindowWithCount(count int, opts ...Option) Observable
	WindowWithTime(timespan Duration, opts ...Option) Observable
	WindowWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable
	ZipFromIterable(iterable Iterable, zipper Func2, opts ...Option) Observable
}

// ObservableImpl implements Observable.
type ObservableImpl struct {
	iterable Iterable
}

func defaultErrorFuncOperator(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	item.SendContext(ctx, dst)
	operatorOptions.stop()
}

func customObservableOperator(f func(ctx context.Context, next chan Item, option Option, opts ...Option), opts ...Option) Observable {
	option := parseOptions(opts...)

	if option.isEagerObservation() {
		next := option.buildChannel()
		ctx := option.buildContext()
		go f(ctx, next, option, opts...)
		return &ObservableImpl{iterable: newChannelIterable(next)}
	}

	return &ObservableImpl{
		iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
			mergedOptions := append(opts, propagatedOptions...)
			option := parseOptions(mergedOptions...)
			next := option.buildChannel()
			ctx := option.buildContext()
			go f(ctx, next, option, mergedOptions...)
			return next
		}),
	}
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
			runSequential(ctx, next, iterable, operatorFactory, option, opts...)
		} else {
			runParallel(ctx, next, iterable.Observe(opts...), operatorFactory, bypassGather, option, opts...)
		}
		return &ObservableImpl{iterable: newChannelIterable(next)}
	}

	if forceSeq || !parallel {
		return &ObservableImpl{
			iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
				mergedOptions := append(opts, propagatedOptions...)
				option := parseOptions(mergedOptions...)

				next := option.buildChannel()
				ctx := option.buildContext()
				runSequential(ctx, next, iterable, operatorFactory, option, mergedOptions...)
				return next
			}),
		}
	}

	if serialized, f := option.isSerialized(); serialized {
		firstItemIDCh := make(chan Item, 1)
		fromCh := make(chan Item, 1)
		obs := &ObservableImpl{
			iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
				mergedOptions := append(opts, propagatedOptions...)
				option := parseOptions(mergedOptions...)

				next := option.buildChannel()
				ctx := option.buildContext()
				observe := iterable.Observe(opts...)
				go func() {
					select {
					case <-ctx.Done():
						return
					case firstItemID := <-firstItemIDCh:
						if firstItemID.Error() {
							firstItemID.SendContext(ctx, fromCh)
							return
						}
						Of(firstItemID.V.(int)).SendContext(ctx, fromCh)
						runParallel(ctx, next, observe, operatorFactory, bypassGather, option, mergedOptions...)
					}
				}()
				runFirstItem(ctx, f, firstItemIDCh, observe, next, operatorFactory, bypassGather, option, mergedOptions...)
				return next
			}),
		}
		return obs.serialize(fromCh, f)
	}

	return &ObservableImpl{
		iterable: newFactoryIterable(func(propagatedOptions ...Option) <-chan Item {
			mergedOptions := append(opts, propagatedOptions...)
			option := parseOptions(mergedOptions...)

			next := option.buildChannel()
			ctx := option.buildContext()
			runParallel(ctx, next, iterable.Observe(mergedOptions...), operatorFactory, bypassGather, option, mergedOptions...)
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
			runSequential(ctx, next, iterable, operatorFactory, option, opts...)
		} else {
			runParallel(ctx, next, iterable.Observe(opts...), operatorFactory, bypassGather, option, opts...)
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
				runSequential(ctx, next, iterable, operatorFactory, option, mergedOptions...)
			} else {
				runParallel(ctx, next, iterable.Observe(mergedOptions...), operatorFactory, bypassGather, option, mergedOptions...)
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
			runSequential(ctx, next, iterable, operatorFactory, option, opts...)
		} else {
			runParallel(ctx, next, iterable.Observe(opts...), operatorFactory, bypassGather, option, opts...)
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
				runSequential(ctx, next, iterable, operatorFactory, option, mergedOptions...)
			} else {
				runParallel(ctx, next, iterable.Observe(mergedOptions...), operatorFactory, bypassGather, option, mergedOptions...)
			}
			return next
		}),
	}
}

func runSequential(ctx context.Context, next chan Item, iterable Iterable, operatorFactory func() operator, option Option, opts ...Option) {
	observe := iterable.Observe(opts...)
	go func() {
		op := operatorFactory()
		stopped := false
		operator := operatorOptions{
			stop: func() {
				if option.getErrorStrategy() == StopOnError {
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

func runParallel(ctx context.Context, next chan Item, observe <-chan Item, operatorFactory func() operator, bypassGather bool, option Option, opts ...Option) {
	wg := sync.WaitGroup{}
	_, pool := option.getPool()
	wg.Add(pool)

	var gather chan Item
	if bypassGather {
		gather = next
	} else {
		gather = make(chan Item, 1)

		// Gather
		go func() {
			op := operatorFactory()
			stopped := false
			operator := operatorOptions{
				stop: func() {
					if option.getErrorStrategy() == StopOnError {
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

	// Scatter
	for i := 0; i < pool; i++ {
		go func() {
			op := operatorFactory()
			stopped := false
			operator := operatorOptions{
				stop: func() {
					if option.getErrorStrategy() == StopOnError {
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
							Of(op).SendContext(ctx, gather)
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
		close(gather)
	}()
}

func runFirstItem(ctx context.Context, f func(interface{}) int, notif chan Item, observe <-chan Item, next chan Item, operatorFactory func() operator, bypassGather bool, option Option, opts ...Option) {
	go func() {
		op := operatorFactory()
		stopped := false
		operator := operatorOptions{
			stop: func() {
				if option.getErrorStrategy() == StopOnError {
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
					i.SendContext(ctx, notif)
				} else {
					op.next(ctx, i, next, operator)
					Of(f(i.V)).SendContext(ctx, notif)
				}
			}
		}
		op.end(ctx, next)
	}()
}

func (o *ObservableImpl) serialize(fromCh chan Item, identifier func(interface{}) int, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()

	ctx := option.buildContext()
	mutex := sync.Mutex{}
	minHeap := binaryheap.NewWith(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})
	status := make(map[int]interface{})
	notif := make(chan struct{})

	var from int
	var counter int64
	src := o.Observe(opts...)
	go func() {
		select {
		case <-ctx.Done():
			close(next)
			return
		case item := <-fromCh:
			if item.Error() {
				item.SendContext(ctx, next)
				close(next)
				return
			}
			from = item.V.(int)
			minHeap.Push(from)
			counter = int64(from)

			// Scatter
			go func() {
				defer close(notif)

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
						select {
						case <-ctx.Done():
							return
						case notif <- struct{}{}:
						}
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
									Of(itemValue).SendContext(ctx, next)
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
		}
	}()

	return &ObservableImpl{
		iterable: newChannelIterable(next),
	}
}
