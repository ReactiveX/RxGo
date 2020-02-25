# Merge Operator

## Overview

Combine multiple Observables into one by merging their emissions.

![](http://reactivex.io/documentation/operators/images/merge.png)

## Example

```go
observable := rxgo.Merge([]rxgo.Observable{
	rxgo.Just(1, 2)(),
	rxgo.Just(3, 4)(),
})
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