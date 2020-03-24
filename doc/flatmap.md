# FlatMap Operator

## Overview

Transform the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.

![](http://reactivex.io/documentation/operators/images/flatMap.c.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().FlatMap(func(i rxgo.Item) rxgo.Observable {
	return rxgo.Just(i.V.(int) * 10, i.V.(int) * 100)()
})
```

Output:

```
10
100
20
200
30
300
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)