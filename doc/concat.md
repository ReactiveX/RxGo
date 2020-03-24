# Concat Operator

## Overview

Emit the emissions from two or more Observables without interleaving them.

![](http://reactivex.io/documentation/operators/images/concat.png)

## Example

```go
observable := rxgo.Concat([]rxgo.Observable{
	rxgo.Just(1, 2, 3)(),
	rxgo.Just(4, 5, 6)(),
})
```

Output:

```
1
2
3
4
5
6
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)