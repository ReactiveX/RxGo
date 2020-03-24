# Take Operator

## Overview

Emit only the first n items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/take.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().Take(2)
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