# TakeUntil Operator

## Overview

Discard any items emitted by an Observable after a second Observable emits an item or terminates.

![](http://reactivex.io/documentation/operators/images/takeUntil.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().TakeUntil(func(i interface{}) bool {
	return i == 3
})
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

* [WithPublishStrategy](options.md#withpublishstrategy)