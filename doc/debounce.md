# Debounce Operator

## Overview

Only emit an item from an Observable if a particular timespan has passed without it emitting another item.

![](http://reactivex.io/documentation/operators/images/debounce.png)

## Example

```go
observable.Debounce(rxgo.WithDuration(250 * time.Millisecond))
```

Output: each item emitted by the Observable if not item has been emitted after 250 milliseconds. 

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)