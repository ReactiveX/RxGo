# Operators

## Index

TODO

## Operator Options

### WithBufferedChannel

Configure the capacity of the output channel.

```go
rxgo.WithBufferedChannel(1) // Create a buffered channel with a 1 capacity
```

### WithContext

Allows passing a context. The Observable will listen to its done signal to close itself.

```go
rxgo.WithContext(ctx)
```

### WithObservationStrategy

* Lazy (default): consume when an Observer starts to subscribe.

```go
rxgo.WithObservation(rxgo.Lazy)
```

* Eager: consumer when the Observable is created:

```go
rxgo.WithObservation(rxgo.Eager)
```

### WithErrorStrategy

* Stop (default): stop processing if the Observable produces an error.

```go
rxgo.WithErrorStrategy(rxgo.Stop)
```

* Continue: continue processing items if the Observable produces an error.

```go
rxgo.WithErrorStrategy(rxgo.Continue)
```

This strategy is propagated to the parent(s) Observable(s).

### WithPool

Convert the operator in a parallel operator and specify the number of concurrent goroutines.

```go
rxgo.WithPool(8) // Creates a pool of 8 goroutines
```

### WithCPUPool

Convert the operator in a parallel operator and specify the number of concurrent goroutines as `runtime.NumCPU()`.

```go
rxgo.WithCPUPool()
```

## All Operator

### Overview

Determine whether all items emitted by an Observable meet some criteria.

![](http://reactivex.io/documentation/operators/images/all.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4}).
	All(func(i interface{}) bool {
		// Check all items are less than 10
		return i.(int) < 10
	})
```

* Output:

```
true
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Amb Operator

### Overview

Given two or more source Observables, emit all of the items from only the first of these Observables to emit an item.

![](http://reactivex.io/documentation/operators/images/amb.png)

### Example

```go
observable := rxgo.Amb([]rxgo.Observable{
	rxgo.Just([]interface{}{1, 2, 3}),
	rxgo.Just([]interface{}{4, 5, 6}),
})
```

* Output:

```
1
2
3
```
or
```
4
5
6
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Average Operator

### Overview

Calculate the average of numbers emitted by an Observable and emits this average.

![](http://reactivex.io/documentation/operators/images/average.png)

### Instances

* AverageFloat32
* AverageFloat64
* AverageInt
* AverageInt8
* AverageInt16
* AverageInt32
* AverageInt64

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4}).AverageInt()
```

* Output:

```
2
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## BackOffRetry Operator

### Overview

Implements a backoff retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.

The backoff configuration relies on [github.com/cenkalti/backoff/v4](github.com/cenkalti/backoff/v4).

![](http://reactivex.io/documentation/operators/images/retry.png)

### Example

```go
// Backoff retry configuration
backOffCfg := backoff.NewExponentialBackOff()
backOffCfg.InitialInterval = 10 * time.Millisecond

observable := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item, done func()) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Error(errors.New("foo"))
	done()
}}).BackOffRetry(backoff.WithMaxRetries(backOffCfg, 2))
```

* Output:

```
1
2
1
2
1
2
foo
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## BufferWithCount Operator

### Overview

BufferWithCount returns an Observable that emits buffers of items it collects from the source Observable.

![](http://reactivex.io/documentation/operators/images/Buffer.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4}).BufferWithCount(3)
```

* Output:

```
1 2 3
4
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## BufferWithTime Operator

### Overview

BufferWithTime returns an Observable that emits buffers of items it collects from the source Observable. 

The resulting Observable starts a new buffer periodically, as determined by the
timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.

When the source Observable completes or encounters an error, the resulting Observable emits the current buffer and propagates the notification from the source Observable.

### Example

```go
// Create the producer
ch := make(chan rxgo.Item, 1)
go func() {
	i := 0
	for range time.Tick(time.Second) {
		ch <- rxgo.Of(i)
		i++
	}
}()

observable := rxgo.FromChannel(ch).
	BufferWithTime(rxgo.WithDuration(3*time.Second), nil)
```

* Output:

```
0 1 2
3 4 5
6 7 8
...
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## BufferWithTimeOrCount Operator

### Overview

BufferWithTimeOrCount returns an Observable that emits buffers of items it collects from the source Observable either from a given count or at a given time interval.

### Example

```go
// Create the producer
ch := make(chan rxgo.Item, 1)
go func() {
	i := 0
	for range time.Tick(time.Second) {
		ch <- rxgo.Of(i)
		i++
	}
}()

observable := rxgo.FromChannel(ch).
	BufferWithTimeOrCount(rxgo.WithDuration(3*time.Second), 2)
```

* Output:

```
0 1
2 3
4 5
...
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Catch Operator

### Overview

Recover from an error by continuing the sequence without error.

### Instances

* `OnErrorResumeNext`: instructs an Observable to pass control to another Observable rather than invoking onError if it encounters an error.
* `OnErrorReturn`: instructs an Observable to emit an item (returned by a specified function) rather than invoking onError if it encounters an error.
* `OnErrorReturnItem`: instructs on Observable to emit an item if it encounters an error.

### Example

#### OnErrorResumeNext

```go
observable := rxgo.Just([]interface{}{1, 2, errors.New("foo")}).
	OnErrorResumeNext(func(e error) rxgo.Observable {
		return rxgo.Just([]interface{}{3, 4})
	})
```

* Output:

```
1
2
3
4
```

#### OnErrorReturn

```go
observable := rxgo.Just([]interface{}{1, errors.New("2"), 3, errors.New("4"), 5}).
	OnErrorReturn(func(err error) interface{} {
		return err.Error()
	})
```

* Output:

```
1
2
3
4
5
```

#### OnErrorReturnItem

```go
observable := rxgo.Just([]interface{}{1, errors.New("2"), 3, errors.New("4"), 5}).
	OnErrorReturnItem("foo")
```

* Output:

```
1
foo
3
foo
5
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## CombineLatest Operator

### Overview

When an item is emitted by either of two Observables, combine the latest item emitted by each Observable via a specified function and emit items based on the results of this function.

![](http://reactivex.io/documentation/operators/images/combineLatest.png)

### Example

```go
observable := rxgo.CombineLatest(func(i ...interface{}) interface{} {
	sum := 0
	for _, v := range i {
		if v == nil {
			continue
		}
		sum += v.(int)
	}
	return sum
}, []rxgo.Observable{
	rxgo.Just([]interface{}{1, 2}),
	rxgo.Just([]interface{}{10, 11}),
})
```

* Output:

```
12
13
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Concat Operator

### Overview

Emit the emissions from two or more Observables without interleaving them.

![](http://reactivex.io/documentation/operators/images/concat.png)

### Example

```go
observable := rxgo.Concat([]rxgo.Observable{
	rxgo.Just([]interface{}{1, 2, 3}),
	rxgo.Just([]interface{}{4, 5, 6}),
})
```

* Output:

```
1
2
3
4
5
6
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Contains Operator

### Overview

Determine whether an Observable emits a particular item or not.

![](http://reactivex.io/documentation/operators/images/contains.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).Contains(func(i interface{}) bool {
	return i == 2
})
```

* Output:

```
true
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Count Operator

### Overview

Count the number of items emitted by the source Observable and emit only this value.

![](http://reactivex.io/documentation/operators/images/Count.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).Count()
```

* Output:

```
3
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Create Operator

### Overview

Create an Observable from scratch by calling observer methods programmatically.

![](http://reactivex.io/documentation/operators/images/create.png)

### Example

```go
observable := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item, done func()) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Of(3)
	done()
}})
```

* Output:

```
1
2
3
```

There are two ways to close the Observable:
* Closing the `next` channel.
* Calling the `done()` function.

Yet, as we can pass multiple producers, using the `done()` function is the recommended approach.

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## DefaultIfEmpty Operator

### Overview

Emit items from the source Observable, or a default item if the source Observable emits nothing.

![](http://reactivex.io/documentation/operators/images/defaultIfEmpty.c.png)

### Example

```go
observable := rxgo.Empty().DefaultIfEmpty(1)
```

* Output:

```
1
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Defer Operator

### Overview

do not create the Observable until the observer subscribes, and create a fresh Observable for each observer.

![](http://reactivex.io/documentation/operators/images/defer.png)

### Example

```go
observable := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item, done func()) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Of(3)
	done()
}})
```

* Output:

```
1
2
3
```

There are two ways to close the Observable:
* Closing the `next` channel.
* Calling the `done()` function.

Yet, as we can pass multiple producers, using the `done()` function is the recommended approach.

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Distinct Operator

### Overview

Suppress duplicate items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/distinct.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 2, 3, 4, 4, 5}).
	Distinct(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	})
```

* Output:

```
1
2
3
4
5
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## DistinctUntilChanged Operator

### Overview

Suppress consecutive duplicate items in the original Observable.

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 2, 1, 1, 3}).
	DistinctUntilChanged(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	})
```

* Output:

```
1
2
1
3
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Do Operator

### Overview

Register an action to take upon a variety of Observable lifecycle events.

![](http://reactivex.io/documentation/operators/images/do.c.png)

### Instances

* DoOnNext
* DoOnError
* DoOnCompleted

Each one returns a `<-chan struct{}` that closes once the Observable terminates.

### Example

#### DoOnNext

```go
<-rxgo.Just([]interface{}{1, 2, 3}).
	DoOnNext(func(i interface{}) {
		fmt.Println(i)
	})
```

* Output:

```
1
2
3
```

#### DoOnError

```go
<-rxgo.Just([]interface{}{1, 2, errors.New("foo")}).
	DoOnError(func(err error) {
		fmt.Println(err)
	})
```

* Output:

```
foo
```

#### DoOnCompleted

```go
<-rxgo.Just([]interface{}{1, 2, 3}).
	DoOnCompleted(func() {
		fmt.Println("done")
	})
```

* Output:

```
done
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## ElementAt Operator

### Overview

Emit only item n emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/elementAt.png)

### Example

```go
observable := rxgo.Just([]interface{}{0, 1, 2, 3, 4}).ElementAt(2)
```

* Output:

```
2
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Empty Operator

### Overview

Create an Observable that emits no items but terminates normally.

![](http://reactivex.io/documentation/operators/images/empty.png)

### Example

```go
observable := rxgo.Empty()
```

* Output:

```
```

## Error Operator

### Overview

Return the eventual Observable error. 

This method is blocking.

### Example

```go
err := rxgo.Just([]interface{}{1, 2, errors.New("foo")}).Error()
fmt.Println(err)
```

* Output:

```
foo
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Errors Operator

### Overview

Return the eventual Observable errors.

This method is blocking.

### Example

```go
errs := rxgo.Just([]interface{}{
	errors.New("foo"),
	errors.New("bar"),
	errors.New("baz"),
}).Errors(rxgo.WithErrorStrategy(rxgo.Continue))
fmt.Println(errs)
```

* Output:

```
[foo bar baz]
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Filter Operator

### Overview

Emit only those items from an Observable that pass a predicate test.

![](http://reactivex.io/documentation/operators/images/filter.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
	Filter(func(i interface{}) bool {
		return i != 2
	})
```

* Output:

```
1
3
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## First Operator

### Overview

Emit only the first item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/first.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).First()
```

* Output:

```
true
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

## FirstOrDefault Operator

### Overview

Similar to `First`, but we pass a default item that will be emitted if the source Observable fails to emit any items.

![](http://reactivex.io/documentation/operators/images/firstOrDefault.png)

### Example

```go
observable := rxgo.Empty().FirstOrDefault(1)
```

* Output:

```
1
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

## FlatMap Operator

### Overview

Transform the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.

![](http://reactivex.io/documentation/operators/images/flatMap.c.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).FlatMap(func(i rxgo.Item) rxgo.Observable {
	return rxgo.Just([]interface{}{
		i.V.(int) * 10,
		i.V.(int) * 100,
	})
})
```

* Output:

```
10
100
20
200
30
300
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## ForEach Operator

### Overview

Subscribe to an Observable and register `OnNext`, `OnError` and `OnCompleted` actions.

It returns a `<-chan struct{}` that closes once the Observable terminates.

### Example

```go
<-rxgo.Just([]interface{}{1, errors.New("foo")}).
	ForEach(
		func(i interface{}) {
			fmt.Printf("next: %v\n", i)
		}, func(err error) {
			fmt.Printf("error: %v\n", err)
		}, func() {
			fmt.Println("done")
		})
```

* Output:

```
next: 1
error: foo
done
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## FromChannel Operator

### Overview

Create a cold observable from a channel.

### Example

```go
ch := make(chan rxgo.Item)
observable := rxgo.FromChannel(ch)
```

The items are buffered in the channel until an Observer subscribes.

## FromEventSource Operator

### Overview

Create a hot observable from a channel.

### Example

```go
ch := make(chan rxgo.Item)
observable := rxgo.FromEventSource(ch)
```

The items are consumed as soon as the observable is created. An Observer will see only the items since the moment he subscribed to the Observable.

### Options

#### WithBackPressureStrategy

* Block (default): block until the Observer is ready to consume the next item.

```go
rxgo.FromEventSource(ch, rxgo.WithBackPressureStrategy(rxgo.Block))
```

* Drop: drop the item if the Observer isn't ready.

```go
rxgo.FromEventSource(ch, rxgo.WithBackPressureStrategy(rxgo.Drop))
```

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## GroupBy Operator

### Overview

Divide an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key.

![](http://reactivex.io/documentation/operators/images/groupBy.c.png)

### Example

```go
count := 3
observable := rxgo.Range(0, 10).GroupBy(count, func(item rxgo.Item) int {
	return item.V.(int) % count
}, rxgo.WithBufferedChannel(10))

for i := range observable.Observe() {
	fmt.Println("New observable:")

	for i := range i.V.(rxgo.Observable).Observe() {
		fmt.Printf("item: %v\n", i.V)
	}
}
```

* Output:

```
New observable:
item: 0
item: 3
item: 6
item: 9
New observable:
item: 1
item: 4
item: 7
item: 10
New observable:
item: 2
item: 5
item: 8
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## IgnoreElements Operator

### Overview

Do not emit any items from an Observable but mirror its termination notification.

![](http://reactivex.io/documentation/operators/images/ignoreElements.c.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, errors.New("foo")}).
	IgnoreElements()
```

* Output:

```
foo
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Interval Operator

### Overview

Create an Observable that emits a sequence of integers spaced by a particular time interval.

![](http://reactivex.io/documentation/operators/images/interval.png)

### Example

```go
observable := rxgo.Interval(rxgo.WithDuration(5 * time.Second))
```

* Output:

```
0 // After 5 seconds
1 // After 10 seconds
2 // After 15 seconds
3 // After 20 seconds
...
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Just Operator

### Overview

Convert an object or a set of objects into an Observable that emits that or those objects.

![](http://reactivex.io/documentation/operators/images/just.png)

### Examples

#### Single Item

```go
observable := rxgo.Just(1)
```

* Output:

```
1
```

#### Multiple Items

```go
observable := rxgo.Just([]interface{}{1, 2, 3})
```

* Output:

```
1
2
3
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

## JustItem Operator

### Overview

Convert an object into a Single that emits that object.

### Examples

```go
observable := rxgo.Just(1)
```

* Output:

```
1
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

## Last Operator

### Overview

Emit only the last item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/last.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).Last()
```

* Output:

```
3
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## LastOrDefault Operator

### Overview

Similar to `Last`, but you pass it a default item that it can emit if the source Observable fails to emit any items.

![](http://reactivex.io/documentation/operators/images/lastOrDefault.png)

### Example

```go
observable := rxgo.Empty().LastOrDefault(1)
```

* Output:

```
1
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Map Operator

### Overview

Transform the items emitted by an Observable by applying a function to each item.

![](http://reactivex.io/documentation/operators/images/map.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
	Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
```

* Output:

```
10
20
30
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Marshal Operator

### Overview

Transform the items emitted by an Observable by applying a marshaller function (`func(interface{}) ([]byte, error)`) to each item.

### Example

```go
type customer struct {
	ID int `json:"id"`
}

observable := rxgo.Just([]customer{
	{
		ID: 1,
	},
	{
		ID: 2,
	},
}).Marshal(json.Marshal)
```

* Output:

```
{"id":1}
{"id":2}
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Max Operator

### Overview

Determine, and emit, the maximum-valued item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/max.png)

### Example

```go
observable := rxgo.Just([]interface{}{2, 5, 1, 6, 3, 4}).
	Max(func(i1 interface{}, i2 interface{}) int {
		return i1.(int) - i2.(int)
	})
```

* Output:

```
6
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Merge Operator

### Overview

Combine multiple Observables into one by merging their emissions.

![](http://reactivex.io/documentation/operators/images/merge.png)

### Example

```go
observable := rxgo.Merge([]rxgo.Observable{
	rxgo.Just([]interface{}{1, 2}),
	rxgo.Just([]interface{}{3, 4}),
})
```

* Output:

```
1
2
3
4
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Min Operator

### Overview

Determine, and emit, the minimum-valued item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/min.png)

### Example

```go
observable := rxgo.Just([]interface{}{2, 5, 1, 6, 3, 4}).
	Max(func(i1 interface{}, i2 interface{}) int {
		return i1.(int) - i2.(int)
	})
```

* Output:

```
1
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Never Operator

### Overview

Create an Observable that emits no items and does not terminate.

![](http://reactivex.io/documentation/operators/images/never.png)

### Example

```go
observable := rxgo.Never()
```

* Output:

```
```

## Range Operator

### Overview

Create an Observable that emits a range of sequential integers.

![](http://reactivex.io/documentation/operators/images/range.png)

### Example

```go
observable := rxgo.Range(0, 3)
```

* Output:

```
0
1
2
3
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Reduce Operator

### Overview

Apply a function to each item emitted by an Observable, sequentially, and emit the final value.

![](http://reactivex.io/documentation/operators/images/reduce.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
	Reduce(func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
		if acc == nil {
			return elem, nil
		}
		return acc.(int) + elem.(int), nil
	})
```

* Output:

```
6
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Repeat Operator

### Overview

Create an Observable that emits a particular item multiple times at a particular frequency.

![](http://reactivex.io/documentation/operators/images/repeat.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
	Repeat(3, rxgo.WithDuration(time.Second))
```

* Output:

```
// Immediately
1
2
3
// After 1 second
1
2
3
// After 2 seconds
1
2
3
...
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Retry Operator

### Overview

Implements a retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.

![](http://reactivex.io/documentation/operators/images/retry.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, errors.New("foo")}).Retry(2)
```

* Output:

```
1
2
1
2
1
2
foo
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Run Operator

### Overview

Create an Observer without consuming the emitted items.

It returns a `<-chan struct{}` that closes once the Observable terminates.

### Example

```go
<-rxgo.Just([]interface{}{1, 2, errors.New("foo")}).Run()
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Sample Operator

### Overview

Emit the most recent item emitted by an Observable within periodic time intervals.

![](http://reactivex.io/documentation/operators/images/sample.png)

### Example

```go
sampledObservable := observable1.Sample(observable2)
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Scan Operator

### Overview

Apply a function to each item emitted by an Observable, sequentially, and emit each successive value.

![](http://reactivex.io/documentation/operators/images/scan.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).
    Scan(func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
        if acc == nil {
            return elem, nil
        }
        return acc.(int) + elem.(int), nil
    })
```

* Output:

```
1
3
6
10
15
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## SequenceEqual Operator

### Overview

Determine whether two Observables emit the same sequence of items.

![](http://reactivex.io/documentation/operators/images/sequenceEqual.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).
	SequenceEqual(rxgo.Just([]interface{}{1, 2, 42, 4, 5}))
```

* Output:

```
false
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Send Operator

### Overview

Send the Observable items to a given channel that will closed once the operation completes.

### Example

```go
ch := make(chan rxgo.Item)
rxgo.Just([]interface{}{1, 2, 3}).Send(ch)
for item := range ch {
	fmt.Println(item.V)
}
```

* Output:

```
1
2
3
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Serialize Operator

### Overview

Force an Observable to make serialized calls and to be well-behaved. 
It takes the starting index and a function that transforms an item value into an index. 

![](http://reactivex.io/documentation/operators/images/serialize.c.png)

### Example

```go
observable := rxgo.Range(0, 1_000_000, rxgo.WithBufferedChannel(capacity)).
	Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}, rxgo.WithCPUPool(), rxgo.WithBufferedChannel(capacity)).
	Serialize(0, func(i interface{}) int {
		return i.(int)
	})
```

* Output:

```
true
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

## Skip Operator

### Overview

Suppress the first n items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/skip.png)

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).Skip(2)
```

* Output:

```
3
4
5
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## SkipLast Operator

### Overview

Suppress the last n items emitted by an Observable.

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).SkipLast(2)
```

* Output:

```
1
2
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## SkipWhile Operator

### Overview

Suppress the Observable items while a condition is not met.

### Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).SkipWhile(func(i interface{}) bool {
	return i != 2
})
```

* Output:

```
2
3
4
5
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Start Operator

### Overview

Create an Observable that emits the return value of a function.

![](http://reactivex.io/documentation/operators/images/start.png)

### Example

```go
observable := rxgo.Start([]rxgo.Supplier{func(ctx context.Context) rxgo.Item {
	return rxgo.Of(1)
}, func(ctx context.Context) rxgo.Item {
	return rxgo.Of(2)
}})
```

* Output:

```
1
2
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

## Thrown Operator

### Overview

Create an Observable that emits no items and terminates with an error.

![](http://reactivex.io/documentation/operators/images/throw.c.png)

### Example

```go
observable := rxgo.Thrown(errors.New("foo"))
```

* Output:

```
foo
```

## Timer Operator

### Overview

Create an Observable that emits a single item after a given delay.

![](http://reactivex.io/documentation/operators/images/timer.png)

### Example

```go
observable := rxgo.Timer(rxgo.WithDuration(5 * time.Second))
```

* Output:

```
{} // After 5 seconds
```

### Options

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

## Unmarshal Operator

### Overview

Transform the items emitted by an Observable by applying an unmarshaller function (`func([]byte, interface{}) error`) to each item. It takes a factory function that initializes the target structure.

### Example

```go
type customer struct {
	ID int `json:"id"`
}

observable := rxgo.Just([][]byte{
	[]byte(`{"id":1}`),
	[]byte(`{"id":2}`),
}).Unmarshal(json.Unmarshal,
	func() interface{} {
		return &customer{}
	})
```

* Output:

```
&{ID:1}
&{ID:2}
```

### Options

#### WithBufferedChannel

https://github.com/ReactiveX/RxGo/wiki/Options#withbufferedchannel

#### WithContext

https://github.com/ReactiveX/RxGo/wiki/Options#withcontext

#### WithObservationStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#withobservationstrategy

#### WithErrorStrategy

https://github.com/ReactiveX/RxGo/wiki/Options#witherrorstrategy

#### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

#### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool