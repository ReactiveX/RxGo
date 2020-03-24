# StartWithIterable Operator

## Overview

Emit a specified Iterable before beginning to emit the items from the source Observable.

![](http://reactivex.io/documentation/operators/images/startWith.png)

## Example

```go
observable := rxgo.Just(3, 4)().StartWith(
	rxgo.Just(1, 2)())
```

Output:

```
1
2
3
4
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)