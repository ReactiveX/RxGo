# Max Operator

## Overview

Determine, and emit, the maximum-valued item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/max.png)

## Example

```go
observable := rxgo.Just(2, 5, 1, 6, 3, 4)().
	Max(func(i1 interface{}, i2 interface{}) int {
		return i1.(int) - i2.(int)
	})
```

Output:

```
6
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)