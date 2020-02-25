# SequenceEqual Operator

## Overview

Determine whether two Observables emit the same sequence of items.

![](http://reactivex.io/documentation/operators/images/sequenceEqual.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().
	SequenceEqual(rxgo.Just(1, 2, 42, 4, 5)())
```

Output:

```
false
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)