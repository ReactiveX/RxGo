# TakeLast Operator

## Overview

Emit only the final n items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/takeLast.n.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().TakeLast(2)
```

Output:

```
4
5
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)