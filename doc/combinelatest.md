# CombineLatest Operator

## Overview

When an item is emitted by either of two Observables, combine the latest item emitted by each Observable via a specified function and emit items based on the results of this function.

![](http://reactivex.io/documentation/operators/images/combineLatest.png)

## Example

```go
observable := rxgo.CombineLatest(func(i ...interface{}) interface{} {
	sum := 0
	for _, v := range i {
		if v == nil {
			continue
		}
		sum += v.(int)
	}
	return sum
}, []rxgo.Observable{
	rxgo.Just(1, 2)(),
	rxgo.Just(10, 11)(),
})
```

Output:

```
12
13
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)