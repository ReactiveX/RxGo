# Start Operator

## Overview

Create an Observable that emits the return value of a function.

![](http://reactivex.io/documentation/operators/images/start.png)

## Example

```go
observable := rxgo.Start([]rxgo.Supplier{func(ctx context.Context) rxgo.Item {
	return rxgo.Of(1)
}, func(ctx context.Context) rxgo.Item {
	return rxgo.Of(2)
}})
```

Output:

```
1
2
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)