# Contains Operator

## Overview

Determine whether an Observable emits a particular item or not.

![](http://reactivex.io/documentation/operators/images/contains.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().Contains(func(i interface{}) bool {
	return i == 2
})
```

Output:

```
true
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)