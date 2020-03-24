# Distinct Operator

## Overview

Suppress duplicate items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/distinct.png)

## Example

```go
observable := rxgo.Just(1, 2, 2, 3, 4, 4, 5)().
	Distinct(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	})
```

Output:

```
1
2
3
4
5
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