# RxGo
[![Join the chat at https://gitter.im/ReactiveX/RxGo](https://badges.gitter.im/ReactiveX/RxGo.svg)](https://gitter.im/ReactiveX/RxGo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ReactiveX/RxGo.svg?branch=master)](https://travis-ci.org/ReactiveX/RxGo)
[![Coverage Status](https://coveralls.io/repos/github/ReactiveX/RxGo/badge.svg?branch=master)](https://coveralls.io/github/ReactiveX/RxGo?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/reactivex/rxgo)](https://goreportcard.com/report/github.com/reactivex/rxgo)

Reactive Extensions for the Go Language

## ReactiveX

[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with Observable streams. This is the official ReactiveX API for the Go language.

ReactiveX is a new, alternative way of asynchronous programming to callbacks, promises and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system. A stream of events is called an [Observable](http://reactivex.io/documentation/contract.html).

An operator is a function that defines an Observable, how and when it should emit data. The list of operators covered is available [here](README.md#supported-operators-in-rxgo).

## RxGo

The RxGo implementation is based on the concept of [pipelines](https://blog.golang.org/pipelines). In a nutshell, a pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function.

![](doc/rx.png)

Let's see at a concrete example with each box being an operator:
* We create a static Observable based on a fixed list of items using `Just` operator.
* We define a transformation function (convert a circle into a square) using `Map` operator.
* We filter each yellow square using `Filter` operator.

In this example, the final items are sent in a channel, available to a consumer. There are many ways to consume or to produce data using RxGo. Publishing the results in a channel is only one of them.

Each operator is a transformation stage. By default, everything is sequential. Yet, we can leverage modern CPU architectures by defining multiple instances of the same operator. Each operator instance being a goroutine connected to a common channel.

The philosophy of RxGo is to implement the ReactiveX concepts and leverage the main Go primitives (channels, goroutines, etc.) so that the integration between the two worlds is as smooth as possible.

## Installation of RxGo v2

```
go get -u github.com/reactivex/rxgo/v2
```

## Getting Started

### Hello World

Let's create our first Observable and consume an item:

```go
observable := rxgo.Just("Hello, World!")()
ch := observable.Observe()
item := <-ch
fmt.Println(item.V)
```

The `Just` operator creates an Observable from a static list of items. `Of(value)` creates an item from a given value. If we want to create an item from an error, we have to use `Error(err)`. This is a difference with the v1 that was accepting directly a value or an error without having to wrap it. What's the rationale for this change? It is to prepare RxGo for the generics feature coming (hopefully) in Go 2.

By the way, the `Just` operator uses currying as a syntactic sugar. This way, it accepts multiple items in the first parameter list and multiple options in the second parameter list. We'll see below how to specify options.

Once the Observable is created, we can observe it using `Observe()`. By default, an Observable is lazy in the sense that it emits items only once a subscription is made. `Observe()` returns a `<-chan rxgo.Item`.

We consumed an item from this channel and printed its value of the item using `item.V`. 

An item is a wrapper on top of a value or an error. We may want to check the type first like this:

```go
item := <-ch
if item.Error() {
    return item.E
}
fmt.Println(item.V)
``` 

`item.Error()` returns a boolean indicating whether an item contains an error. Then, we use either `item.E` to get the error or `item.V` to get the value.

By default, an Observable is stopped once an error is produced. However, there are special operators to deal with errors (e.g. `OnError`, `Retry`, etc.)

It is also possible to consume items using callbacks:

```go
observable.ForEach(func(v interface{}) {
    fmt.Printf("received: %v\n", v)
}, func(err error) {
    fmt.Printf("error: %e\n", err)
}, func() {
    fmt.Println("observable is closed")
})
```

In this example, we passed 3 functions:
* A `NextFunc` triggered when a value item is emitted.
* An `ErrFunc` triggered when an error item is emitted.
* A `CompletedFunc` triggered once the Observable is completed.

`ForEach` is non-blocking. Yet, it returns a notification channel that will be closed once the Observable completes. Hence, to make the previous code blocking, we simply need to use `<-`:

```go
<-observable.ForEach(...)
```

### Real-World Example

Let's say we want to implement a stream that consumes the following `Customer` structure:
```go
type Customer struct {
	ID             int
	Name, LastName string
	Age            int
	TaxNumber      string
}
```

We create an producer that will emit `Customer`s to a given `chan rxgo.Item` and create an Observable from it:
```go
// Create the input channel
ch := make(chan rxgo.Item)
// Data producer
go producer(ch)

// Create an Observable
observable := rxgo.FromChannel(ch)
```

Then, we need to perform the two following operations:
* Filter the customers whose age is below 18.
* Enrich each customer with a tax number. Retrieving a tax number is done for example by an IO-bound function doing an external REST call.

As the enriching step is IO-bound, it might be interesting to parallelize it within a given pool of goroutines.
Yet, for some reason, all the `Customer` items need to be produced sequentially based on its `ID`.

```go
observable.
	Filter(func(item interface{}) bool {
		// Filter operation
		customer := item.(Customer)
		return customer.Age > 18
	}).
	Map(func(_ context.Context, item interface{}) (interface{}, error) {
		// Enrich operation
		customer := item.(Customer)
		taxNumber, err := getTaxNumber(customer)
		if err != nil {
			return nil, err
		}
		customer.TaxNumber = taxNumber
		return customer, nil
	},
		// Create multiple instances of the map operator
		rxgo.WithPool(pool),
		// Serialize the items emitted by their Customer.ID
		rxgo.Serialize(func(item interface{}) int {
			customer := item.(Customer)
			return customer.ID
		}), rxgo.WithBufferedChannel(1))
``` 

In the end, we consume the items using `ForEach()` or `Observe()` for example. `Observe()` returns a `<-chan Item`:
```go
for customer := range observable.Observe() {
	if customer.Error() {
		return err
	}
	fmt.Println(customer)
}
```

## Observable Types

### Hot vs Cold Observables

In the Rx world, there is a distinction between cold and hot Observables. When the data is produced by the Observable itself, it is a cold Observable. When the data is produced outside the Observable, it is a hot Observable. Usually, when we don't want to create a producer over and over again, we favour a hot Observable.

In RxGo, there is a similar concept.

First, let's create a **hot** Observable using `FromChannel` operator and see the implications:

```go
ch := make(chan rxgo.Item)
go func() {
    for i := 0; i < 3; i++ {
        ch <- rxgo.Of(i)
    }
    close(ch)
}()
observable := rxgo.FromChannel(ch)

// First Observer
for item := range observable.Observe() {
    fmt.Println(item.V)
}

// Second Observer
for item := range observable.Observe() {
    fmt.Println(item.V)
}
```

The result of this execution is:

```
0
1
2
```

It means, the first Observer already consumed all items. And nothing left for others.  
Though this behavior can be altered with [Connectable](#connectable-observable) Observables.  
The main point here is the goroutine produced those items.

On the other hand, let's create a **cold** Observable using `Defer` operator:

```go
observable := rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
    for i := 0; i < 3; i++ {
        ch <- rxgo.Of(i)
    }
}})

// First Observer
for item := range observable.Observe() {
    fmt.Println(item.V)
}

// Second Observer
for item := range observable.Observe() {
    fmt.Println(item.V)
}
```

Now, the result is:

```go
0
1
2
0
1
2
```

In the case of a cold observable, the stream was created independently for every observer.

Again, **hot** vs **cold** Observables are not about how you consume items, it's about where data is produced.  
Good example for hot Observable are price ticks from a trading exchange.  
And if you teach an Observable to fetch products from a database, then yield them one by one, you will create the **cold** Observable.

### Backpressure

There is another operator called `FromEventSource` that creates an Observable from a channel. The difference between `FromChannel` operator is that as soon as the Observable is created, it starts to emit items regardless if there is an Observer or not. Hence, the items emitted by an Observable without Observer(s) are lost (whilst they are buffered with `FromChannel` operator).

A use case with `FromEventSource` operator is for example telemetry. We may not be interested in all the data produced from the very beginning of a stream. Only the data since we started to observe it.

Once we start observing an Observable created with `FromEventSource`, we can configure the backpressure strategy. By default, it is blocking (there is a guaranteed delivery for the items emitted after we observe it). We can override this strategy this way:

```go
observable := rxgo.FromEventSource(input, rxgo.WithBackPressureStrategy(rxgo.Drop))
```

The `Drop` strategy means that if the pipeline after `FromEventSource` was not ready to consume an item, this item is dropped.

By default, a channel connecting operators is non-buffered. We can override this behaviour like this:

```go
observable.Map(transform, rxgo.WithBufferedChannel(42))
```

Each operator has an `opts ...Option` parameter allowing to pass such options.

### Lazy vs Eager Observation

The default observation strategy is lazy. It means the items emitted by an Observable are processed by an operator once we start observing it. We can change this behaviour this way:

```go
observable := rxgo.FromChannel(ch).Map(transform, rxgo.WithObservation(rxgo.Eager))
```

In this case, the `Map` operator is triggered whenever an item is produced even without any Observer.

### Sequential vs Parallel Operators

By default, each operator is sequential. One operator being one goroutine instance. We can override it using the following option:

```go
observable.Map(transform, rxgo.WithPool(32))
```

In this example, we create a pool of 32 goroutines that consume items concurrently from the same channel. If the operation is CPU-bound, we can use the `WithCPUPool()` option that creates a pool based on the number of logical CPUs.

### Connectable Observable

A Connectable Observable resembles an ordinary Observable, except that it does not begin emitting items when it is subscribed to, but only when its connect() method is called. In this way, you can wait for all intended Subscribers to subscribe to the Observable before the Observable begins emitting items.

Let's create a Connectable Observable using `rxgo.WithPublishStrategy`:

```go
ch := make(chan rxgo.Item)
go func() {
	ch <- rxgo.Of(1)
	ch <- rxgo.Of(2)
	ch <- rxgo.Of(3)
	close(ch)
}()
observable := rxgo.FromChannel(ch, rxgo.WithPublishStrategy())
```

Then, we create two Observers:

```go
observable.Map(func(_ context.Context, i interface{}) (interface{}, error) {
	return i.(int) + 1, nil
}).DoOnNext(func(i interface{}) {
	fmt.Printf("First observer: %d\n", i)
})

observable.Map(func(_ context.Context, i interface{}) (interface{}, error) {
	return i.(int) * 2, nil
}).DoOnNext(func(i interface{}) {
	fmt.Printf("Second observer: %d\n", i)
})
```

If `observable` was not a Connectable Observable, as `DoOnNext` creates an Observer, the source Observable would have begun emitting items. Yet, in the case of a Connectable Observable, we have to call `Connect()`:

```go
observable.Connect()
``` 

Once `Connect()` is called, the Connectable Observable begin to emit items.

There is another important change with a regular Observable. A Connectable Observable publishes its items. It means, all the Observers receive a copy of the items.

Here is an example with a regular Observable:

```go
ch := make(chan rxgo.Item)
go func() {
	ch <- rxgo.Of(1)
	ch <- rxgo.Of(2)
	ch <- rxgo.Of(3)
	close(ch)
}()
// Create a regular Observable
observable := rxgo.FromChannel(ch)

// Create the first Observer
observable.DoOnNext(func(i interface{}) {
	fmt.Printf("First observer: %d\n", i)
})

// Create the second Observer
observable.DoOnNext(func(i interface{}) {
	fmt.Printf("Second observer: %d\n", i)
})
```

```
First observer: 1
First observer: 2
First observer: 3
```

Now, with a Connectable Observable:

```go
ch := make(chan rxgo.Item)
go func() {
	ch <- rxgo.Of(1)
	ch <- rxgo.Of(2)
	ch <- rxgo.Of(3)
	close(ch)
}()
// Create a Connectable Observable
observable := rxgo.FromChannel(ch, rxgo.WithPublishStrategy())

// Create the first Observer
observable.DoOnNext(func(i interface{}) {
	fmt.Printf("First observer: %d\n", i)
})

// Create the second Observer
observable.DoOnNext(func(i interface{}) {
	fmt.Printf("Second observer: %d\n", i)
})

disposed, cancel := observable.Connect()
go func() {
	// Do something
	time.Sleep(time.Second)
	// Then cancel the subscription
	cancel()
}()
// Wait for the subscription to be disposed
<-disposed
```

```
Second observer: 1
First observer: 1
First observer: 2
First observer: 3
Second observer: 2
Second observer: 3
```

### Observable, Single and Optional Single

An Iterable is an object that can be observed using `Observe(opts ...Option) <-chan Item`.

An Iterable can be either:
* An Observable: emit 0 or multiple items
* A Single: emit 1 item
* An Optional Single: emit 0 or 1 item

## Documentation

Package documentation: [https://pkg.go.dev/github.com/reactivex/rxgo/v2](https://pkg.go.dev/github.com/reactivex/rxgo/v2)

### Assert API

How to use the [assert API](doc/assert.md) to write unit tests while using RxGo.

### Operator Options

[Operator options](doc/options.md)

### Creating Observables
* [Create](doc/create.md) — create an Observable from scratch by calling observer methods programmatically
* [Defer](doc/defer.md) — do not create the Observable until the observer subscribes, and create a fresh Observable for each observer
* [Empty](doc/empty.md)/[Never](doc/never.md)/[Thrown](doc/thrown.md) — create Observables that have very precise and limited behaviour
* [FromChannel](doc/fromchannel.md) — create an Observable based on a lazy channel
* [FromEventSource](doc/fromeventsource.md) — create an Observable based on an eager channel
* [Interval](doc/interval.md) — create an Observable that emits a sequence of integers spaced by a particular time interval
* [Just](doc/just.md) — convert a set of objects into an Observable that emits that or those objects
* [JustItem](doc/justitem.md) — convert one object into a Single that emits this object
* [Range](doc/range.md) — create an Observable that emits a range of sequential integers
* [Repeat](doc/repeat.md) — create an Observable that emits a particular item or sequence of items repeatedly
* [Start](doc/start.md) — create an Observable that emits the return value of a function
* [Timer](doc/timer.md) — create an Observable that completes after a specified delay

### Transforming Observables
* [Buffer](doc/buffer.md) — periodically gather items from an Observable into bundles and emit these bundles rather than emitting the items one at a time
* [FlatMap](doc/flatmap.md) — transform the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable
* [GroupBy](doc/groupby.md) — divide an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key
* [Map](doc/map.md) — transform the items emitted by an Observable by applying a function to each item
* [Marshal](doc/marshal.md) — transform the items emitted by an Observable by applying a marshalling function to each item
* [Scan](doc/scan.md) — apply a function to each item emitted by an Observable, sequentially, and emit each successive value
* [Unmarshal](doc/unmarshal.md) — transform the items emitted by an Observable by applying an unmarshalling function to each item
* [Window](doc/window.md) — apply a function to each item emitted by an Observable, sequentially, and emit each successive value

### Filtering Observables
* [Debounce](doc/debounce.md) — only emit an item from an Observable if a particular timespan has passed without it emitting another item
* [Distinct](doc/distinct.md)/[DistinctUntilChanged](doc/distinctuntilchanged.md) — suppress duplicate items emitted by an Observable
* [ElementAt](doc/elementat.md) — emit only item n emitted by an Observable
* [Filter](doc/filter.md) — emit only those items from an Observable that pass a predicate test
* [First](doc/first.md)/[FirstOrDefault](doc/firstordefault.md) — emit only the first item or the first item that meets a condition, from an Observable
* [IgnoreElements](doc/ignoreelements.md) — do not emit any items from an Observable but mirror its termination notification
* [Last](doc/last.md)/[LastOrDefault](doc/lastordefault.md) — emit only the last item emitted by an Observable
* [Sample](doc/sample.md) — emit the most recent item emitted by an Observable within periodic time intervals
* [Skip](doc/skip.md) — suppress the first n items emitted by an Observable
* [SkipLast](doc/skiplast.md) — suppress the last n items emitted by an Observable
* [Take](doc/take.md) — emit only the first n items emitted by an Observable
* [TakeLast](doc/takelast.md) — emit only the last n items emitted by an Observable

### Combining Observables
* [CombineLatest](doc/combinelatest.md) — when an item is emitted by either of two Observables, combine the latest item emitted by each Observable via a specified function and emit items based on the results of this function
* [Join](doc/join.md) — combine items emitted by two Observables whenever an item from one Observable is emitted during a time window defined according to an item emitted by the other Observable
* [Merge](doc/merge.md) — combine multiple Observables into one by merging their emissions
* [StartWithIterable](doc/startwithiterable.md) — emit a specified sequence of items before beginning to emit the items from the source Iterable
* [ZipFromIterable](doc/zipfromiterable.md) — combine the emissions of multiple Observables together via a specified function and emit single items for each combination based on the results of this function

### Error Handling Operators
* [Catch](doc/catch.md) — recover from an onError notification by continuing the sequence without error
* [Retry](doc/retry.md)/[BackOffRetry](doc/backoffretry.md) — if a source Observable sends an onError notification, resubscribe to it in the hopes that it will complete without error

### Observable Utility Operators
* [Do](doc/do.md) - register an action to take upon a variety of Observable lifecycle events
* [Run](doc/run.md) — create an Observer without consuming the emitted items
* [Send](doc/send.md) — send the Observable items in a specific channel
* [Serialize](doc/serialize.md) — force an Observable to make serialized calls and to be well-behaved
* [TimeInterval](doc/timeinterval.md) — convert an Observable that emits items into one that emits indications of the amount of time elapsed between those emissions
* [Timestamp](doc/timestamp.md) — attach a timestamp to each item emitted by an Observable

### Conditional and Boolean Operators
* [All](doc/all.md) — determine whether all items emitted by an Observable meet some criteria
* [Amb](doc/amb.md) — given two or more source Observables, emit all of the items from only the first of these Observables to emit an item
* [Contains](doc/contains.md) — determine whether an Observable emits a particular item or not
* [DefaultIfEmpty](doc/defaultifempty.md) — emit items from the source Observable, or a default item if the source Observable emits nothing
* [SequenceEqual](doc/sequenceequal.md) — determine whether two Observables emit the same sequence of items
* [SkipWhile](doc/skipwhile.md) — discard items emitted by an Observable until a specified condition becomes false
* [TakeUntil](doc/takeuntil.md) — discard items emitted by an Observable after a second Observable emits an item or terminates
* [TakeWhile](doc/takewhile.md) — discard items emitted by an Observable after a specified condition becomes false

### Mathematical and Aggregate Operators
* [Average](doc/average.md) — calculates the average of numbers emitted by an Observable and emits this average
* [Concat](doc/concat.md) — emit the emissions from two or more Observables without interleaving them
* [Count](doc/count.md) — count the number of items emitted by the source Observable and emit only this value
* [Max](doc/max.md) — determine, and emit, the maximum-valued item emitted by an Observable
* [Min](doc/min.md) — determine, and emit, the minimum-valued item emitted by an Observable
* [Reduce](doc/reduce.md) — apply a function to each item emitted by an Observable, sequentially, and emit the final value
* [Sum](doc/sum.md) — calculate the sum of numbers emitted by an Observable and emit this sum

### Operators to Convert Observables
* [Error](doc/error.md)/[Errors](doc/errors.md) — convert an observable into an eventual error or list of errors
* [ToMap](doc/tomap.md)/[ToMapWithValueSelector](doc/tomapwithvalueselector.md)/[ToSlice](doc/toslice.md) — convert an Observable into another object or data structure

## Contributions

All contributions are very welcome! Be sure you check out the [Contributions](https://github.com/ReactiveX/RxGo/wiki/Contributions) and [Roadmap](https://github.com/ReactiveX/RxGo/wiki/Roadmap) pages first.
