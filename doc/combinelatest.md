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
	rxgo.Just([]interface{}{1, 2}),
	rxgo.Just([]interface{}{10, 11}),
})
```

Output:

```
12
13
```

## Options

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)