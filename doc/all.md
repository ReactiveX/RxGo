# All Operator

## Overview

Determine whether all items emitted by an Observable meet some criteria.

![](http://reactivex.io/documentation/operators/images/all.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4)().
	All(func(i interface{}) bool {
		// Check all items are less than 10
		return i.(int) < 10
	})
```

Output:

```
true
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)