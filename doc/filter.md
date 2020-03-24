# Filter Operator

## Overview

Emit only those items from an Observable that pass a predicate test.

![](http://reactivex.io/documentation/operators/images/filter.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().
	Filter(func(i interface{}) bool {
		return i != 2
	})
```

Output:

```
1
3
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

### Serialize

[Detail](options.md#serialize)

* [WithPublishStrategy](options.md#withpublishstrategy)