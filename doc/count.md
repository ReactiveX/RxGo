# Count Operator

## Overview

Count the number of items emitted by the source Observable and emit only this value.

![](http://reactivex.io/documentation/operators/images/Count.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().Count()
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