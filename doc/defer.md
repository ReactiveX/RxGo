# Defer Operator

## Overview

do not create the Observable until the observer subscribes, and create a fresh Observable for each observer.

![](http://reactivex.io/documentation/operators/images/defer.png)

## Example

```go
observable := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Of(3)
}})
```

Output:

```
1
2
3
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)