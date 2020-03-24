# TakeWhile Operator

## Overview

Mirror items emitted by an Observable until a specified condition becomes false.

![](http://reactivex.io/documentation/operators/images/takeWhile.c.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().TakeWhile(func(i interface{}) bool {
	return i != 3
})
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