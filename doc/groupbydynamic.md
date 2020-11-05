# GroupBy Operator

## Overview

Divides an Observable into a dynamic set of Observables that each emit GroupedObservable from the original Observable, organized by key.

`GroupByDyDynamic` differs from [GroupBy](groupby.md) for two reasons:
 * We don't need to pass a fixed set length.
 * The distribution function is a `func(rxgo.Item) string` instead of a `func(rxgo.Item) int`. The rationale is because of possible collisions. For example, if our distribution function produces 128-bit UUIDs, there is a collision risk if such a UUID has to be casted into an int.   

![](http://reactivex.io/documentation/operators/images/groupBy.c.png)

## Example

```go
count := 3
observable := rxgo.Range(0, 10).GroupByDynamic(func(item rxgo.Item) string {
    return strconv.Itoa(item.V.(int) % count)
}, rxgo.WithBufferedChannel(10))

for i := range observable.Observe() {
    groupedObservable := i.V.(rxgo.GroupedObservable)
    fmt.Printf("New observable: %d\n", groupedObservable.Key)

    for i := range groupedObservable.Observe() {
        fmt.Printf("item: %v\n", i.V)
    }
}
```

Output:

```
New observable: 0
item: 0
item: 3
item: 6
item: 9
New observable: 1
item: 1
item: 4
item: 7
item: 10
New observable: 2
item: 2
item: 5
item: 8
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)