# ElementAt Operator

## Overview

Emit only item n emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/elementAt.png)

## Example

```go
observable := rxgo.Just(0, 1, 2, 3, 4)().ElementAt(2)
```

Output:

```
2
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)