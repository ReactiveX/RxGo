# RxGo
[![Join the chat at https://gitter.im/ReactiveX/RxGo](https://badges.gitter.im/ReactiveX/RxGo.svg)](https://gitter.im/ReactiveX/RxGo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ReactiveX/RxGo.svg?branch=master)](https://travis-ci.org/ReactiveX/RxGo)
[![Coverage Status](https://coveralls.io/repos/github/ReactiveX/RxGo/badge.svg?branch=master)](https://coveralls.io/github/ReactiveX/RxGo?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/reactivex/rxgo)](https://goreportcard.com/report/github.com/reactivex/rxgo)

Reactive Extensions for the Go Language

## ReactiveX

[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with Observable streams. This is the official ReactiveX API for the Go language.

ReactiveX is a new, alternative way of asynchronous programming to callbacks, promises and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system. A stream of events is called an [Observable](http://reactivex.io/documentation/contract.html).

An operator is basically a function that defines an Observable, how and when it should emit data. The list of operators covered is available [here](README.md#supported-operators-in-rxgo).

## RxGo

The RxGo implementation is based on the [pipelines](https://blog.golang.org/pipelines) concept. In a nutshell, a pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function.

![](res/rx.png)

Let's see at a concrete example with each box being an operator:
* We create a static Observable based on a fixed list of items using `Just` operator.
* We define a transformation function (convert a circle into a square) using `Map` operator.
* We filter each yellow square using `Filter` operator.

In this example, the final items are sent in a channel, available to a consumer. There are many ways to consume or to produce data using RxGo. Publishing the results in a channel is only one of them.

Each operator is a transformation stage. By default, everything is sequential. Yet, we can leverage modern CPU architectures by defining multiple instances of the same operator. Each operator instance being a goroutine connected to a common channel.

The philosophy of RxGo is to implement the ReactiveX concepts and leverage the main Go primitives (channels, goroutines, etc.) so that the integration between the two worlds is as smooth as possible.

## Installation of RxGo v2

```
go get -u github.com/reactivex/rxgo/v2@develop
```

## Getting Started

The following documentation gives an overview of RxGo. If you need more information, please check at the Wiki.

### Hello World

Let's create our first Observable and consume an item:

```go
observable := rxgo.Just([]rxgo.Item{rxgo.Of("Hello, World!")})
ch := observable.Observe()
item := <-ch
fmt.Println(item.V)
```

The `Just` operator creates an Observable from a static list of items. `Of(value)` creates an item from a given value. If we want to create an item from an error, we have to use `Error(err)`. This is a difference with the v1 that was accepting directly a value or an error without having to wrap it. What's the rationale for this change? It is to prepare RxGo for the generics feature coming (hopefully) in Go 2.

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

### Deep Dive

Let's implement a more complete example. We will create an Observable from a channel and implement three operators (`Map`, `Filter` and `Reduce`):

```go
// Create the input channel
ch := make(chan rxgo.Item)
// Create the data producer
go producer(ch)

// Create an Observable
observable := rxgo.FromChannel(ch).
    Map(func(item interface{}) (interface{}, error) {
        if num, ok := item.(int); ok {
            return num * 2, nil
        }
        return nil, errors.New("input error")
    }).
    Filter(func(item interface{}) bool {
        return item != 4
    }).
    Reduce(func(acc interface{}, item interface{}) (interface{}, error) {
        if acc == nil {
            return 1, nil
        }
        return acc.(int) * item.(int), nil
    })

// Observe it
product := <-observable.Observe()

fmt.Println(product)
```

We started by defining a `chan rxgo.Item` that will act as an input channel. We also created a `producer` goroutine that will be in charge to produce the inputs for our Observable.

The Observable was created using `FromChannel` operator. This is basically a wrapper on top of an existing channel. Then, we implemented 3 operators:
* `Map` to double each input.
* `Filter` to filter items equals to 4.
* `Reduce` to perform a reduction based on the product of each item.

In the end, `Observe` returns the output channel. We consume the output using the standard `<-` operator and we print the value.

In this example, the Observable does produce zero or one item. This is what we would expect from a basic `Reduce` operator. 

An Observable that produces exactly one item is a Single (e.g. `Count`, `FirstOrDefault`, etc.). An Observable that produces zero or one item is called an OptionalSingle (e.g. `Reduce`, `Max`, etc.).

The operators producing a Single or OptionalSingle emit an item (or no item in the case of an OptionalSingle) when the stream is closed. In this example, the action that triggers the closure of the stream is when the `input` channel will be closed (most likely by the producer).

## Observable Types

### Hot vs Cold Observables

In the Rx world, there is a distinction between hot and cold Observable. When the data is produced by the Observable itself, it is a cold Observable. When the data is produced outside the Observable, it is a hot Observable. Usually, when we don't want to create a producer over and over again, we favour a hot Observable.

In RxGo, there is a similar concept.

First, let's create a cold Observable using `FromChannel` operator and see the implications:

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

It means, the first Observer consumed already all the items.

On the other hand, let's create a hot Observable using `Defer` operator:

```go
observable := rxgo.Defer([]Producer{func(_ context.Context, ch chan<- rxgo.Item, done func()) {
    for i := 0; i < 3; i++ {
        ch <- rxgo.Of(i)
    }
    done()
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

In the case of a hot observable created with `Defer`, the stream is reproducible. Depending on our use case, we may favour one or the other approach.

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
observable := rxgo.FromChannel(ch).Map(transform, rxgo.WithEagerObservation())
```

In this case, the `Map` operator is triggered whenever an item is produced even without any Observer.

### Sequential vs Parallel Operators

By default, each operator is sequential. One operator being one goroutine instance. We can override it using the following option:

```go
observable.Map(transform, rxgo.WithPool(32))
```

In this example, we create a pool of 32 goroutines that consume items concurrently from the same channel. If the operation is CPU-bound, we can use the `WithCPUPool()` option that creates a pool based on the number of logical CPUs.

## Supported Operators in RxGo

### Creating Observables
* [Create](http://reactivex.io/documentation/operators/create.html) - create an Observable from scratch by calling observer methods programmatically
* [Defer](http://reactivex.io/documentation/operators/defer.html) - do not create the Observable until the observer subscribes, and create a fresh Observable for each observer
* [Empty/Never](http://reactivex.io/documentation/operators/empty-never-throw.html) — create Observables that have very precise and limited behaviour
* FromChannel — create an Observable based on a lazy channel
* FromEventSource — create an Observable based on an eager channel
* [Interval](http://reactivex.io/documentation/operators/interval.html) — create an Observable that emits a sequence of integers spaced by a particular time interval
* [Just](http://reactivex.io/documentation/operators/just.html) — convert a set of objects into an Observable that emits that or those objects
* JustItem — convert one object into a Single that emits this object
* [Range](http://reactivex.io/documentation/operators/range.html) — create an Observable that emits a range of sequential integers
* [Repeat](http://reactivex.io/documentation/operators/repeat.html) — create an Observable that emits a particular item or sequence of items repeatedly
* [Start](http://reactivex.io/documentation/operators/start.html) — create an Observable that emits the return value of a function
* [Timer](http://reactivex.io/documentation/operators/timer.html) — create an Observable that emits a single item after a given delay

### Transforming Observables
* [BufferWithCount/BufferWithTime/BufferWithTimeOrCount](http://reactivex.io/documentation/operators/buffer.html) — periodically gather items from an Observable into bundles and emit these bundles rather than emitting the items one at a time
* [FlatMap](http://reactivex.io/documentation/operators/flatmap.html) — transform the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable
* [GroupBy](http://reactivex.io/documentation/operators/groupby.html) — divide an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key
* [Map](http://reactivex.io/documentation/operators/map.html) — transform the items emitted by an Observable by applying a function to each item
* Marshal - transform the items emitted by an Observable by applying a marshalling function to each item
* [Scan](http://reactivex.io/documentation/operators/scan.html) — apply a function to each item emitted by an Observable, sequentially, and emit each successive value
* Unmarshal - transform the items emitted by an Observable by applying an unmarshalling function to each item

### Filtering Observables
* [Distinct/DistinctUntilChanged](http://reactivex.io/documentation/operators/distinct.html) — suppress duplicate items emitted by an Observable
* [ElementAt](http://reactivex.io/documentation/operators/elementat.html) — emit only item n emitted by an Observable
* [Filter](http://reactivex.io/documentation/operators/filter.html) — emit only those items from an Observable that pass a predicate test
* [First/FirstOrDefault](http://reactivex.io/documentation/operators/first.html) — emit only the first item or the first item that meets a condition, from an Observable
* [IgnoreElements](http://reactivex.io/documentation/operators/ignoreelements.html) — do not emit any items from an Observable but mirror its termination notification
* [Last/LastOrDefault](http://reactivex.io/documentation/operators/last.html) — emit only the last item emitted by an Observable
* [Sample](http://reactivex.io/documentation/operators/sample.html) — emit the most recent item emitted by an Observable within periodic time intervals
* [Skip](http://reactivex.io/documentation/operators/skip.html) — suppress the first n items emitted by an Observable
* [SkipLast](http://reactivex.io/documentation/operators/skiplast.html) — suppress the last n items emitted by an Observable
* [Take](http://reactivex.io/documentation/operators/take.html) — emit only the first n items emitted by an Observable
* [TakeLast](http://reactivex.io/documentation/operators/takelast.html) — emit only the last n items emitted by an Observable

### Combining Observables
* [CombineLatest](http://reactivex.io/documentation/operators/combinelatest.html) — when an item is emitted by either of two Observables, combine the latest item emitted by each Observable via a specified function and emit items based on the results of this function
* [Merge](http://reactivex.io/documentation/operators/merge.html) — combine multiple Observables into one by merging their emissions
* [StartWithIterable](http://reactivex.io/documentation/operators/startwith.html) — emit a specified sequence of items before beginning to emit the items from the source Iterable
* [ZipFromIterable](http://reactivex.io/documentation/operators/zip.html) — combine the emissions of multiple Observables together via a specified function and emit single items for each combination based on the results of this function

### Error Handling Operators
* [OnErrorResumeNext/OnErrorReturn/OnErrorReturnItem](http://reactivex.io/documentation/operators/catch.html) — recover from an onError notification by continuing the sequence without error
* [Retry](http://reactivex.io/documentation/operators/retry.html) — if a source Observable sends an onError notification, resubscribe to it in the hopes that it will complete without error

### Observable Utility Operators
* [DoOnNext/DoOnError/DoOnCompleted](http://reactivex.io/documentation/operators/do.html) - register an action to take upon a variety of Observable lifecycle events
* Run - create an Observer without consuming the emitted items
* Send - send the Observable items in a specific channel

### Conditional and Boolean Operators
* [All](http://reactivex.io/documentation/operators/all.html) — determine whether all items emitted by an Observable meet some criteria
* [Amb](http://reactivex.io/documentation/operators/amb.html) — given two or more source Observables, emit all of the items from only the first of these Observables to emit an item
* [Contains](http://reactivex.io/documentation/operators/contains.html) — determine whether an Observable emits a particular item or not
* [DefaultIfEmpty](http://reactivex.io/documentation/operators/defaultifempty.html) — emit items from the source Observable, or a default item if the source Observable emits nothing
* [SequenceEqual](http://reactivex.io/documentation/operators/sequenceequal.html) — determine whether two Observables emit the same sequence of items
* [SkipWhile](http://reactivex.io/documentation/operators/skipwhile.html) — discard items emitted by an Observable until a specified condition becomes false
* [TakeUntil](http://reactivex.io/documentation/operators/takeuntil.html) — discard items emitted by an Observable after a second Observable emits an item or terminates
* [TakeWhile](http://reactivex.io/documentation/operators/takewhile.html) — discard items emitted by an Observable after a specified condition becomes false

### Mathematical and Aggregate Operators
* [AverageFloat32/AverageFloat64/AverageInt/AverageInt8/AverageInt16/AverageInt32/AverageInt64](http://reactivex.io/documentation/operators/average.html) — calculates the average of numbers emitted by an Observable and emits this average
* [Concat](http://reactivex.io/documentation/operators/concat.html) — emit the emissions from two or more Observables without interleaving them
* [Count](http://reactivex.io/documentation/operators/count.html) — count the number of items emitted by the source Observable and emit only this value
* [Max](http://reactivex.io/documentation/operators/max.html) — determine, and emit, the maximum-valued item emitted by an Observable
* [Min](http://reactivex.io/documentation/operators/min.html) — determine, and emit, the minimum-valued item emitted by an Observable
* [Reduce](http://reactivex.io/documentation/operators/reduce.html) — apply a function to each item emitted by an Observable, sequentially, and emit the final value
* [SumFloat32/SumFloat64/SumInt64](http://reactivex.io/documentation/operators/sum.html) — calculate the sum of numbers emitted by an Observable and emit this sum

### Operators to Convert Observables
* [ToMap/ToMapWithValueSelector/ToSlice](http://reactivex.io/documentation/operators/to.html) — convert an Observable into another object or data structure

## Contributions

All contributions are very welcome! Be sure you check out the [Contributions](https://github.com/ReactiveX/RxGo/wiki/Contributions) and [Roadmap](https://github.com/ReactiveX/RxGo/wiki/Roadmap) pages first.
