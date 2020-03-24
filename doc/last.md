# Last Operator

## Overview

Emit only the last item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/last.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().Last()
```

Output:

```
3
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)