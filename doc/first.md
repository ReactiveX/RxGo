# First Operator

## Overview

Emit only the first item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/first.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().First()
```

Output:

```
true
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)