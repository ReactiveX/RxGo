# Reduce Operator

## Overview

Apply a function to each item emitted by an Observable, sequentially, and emit the final value.

![](http://reactivex.io/documentation/operators/images/reduce.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().
	Reduce(func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
		if acc == nil {
			return elem, nil
		}
		return acc.(int) + elem.(int), nil
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