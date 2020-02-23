# All Operator

## Overview

Determine whether all items emitted by an Observable meet some criteria.

![](http://reactivex.io/documentation/operators/images/all.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4}).
	All(func(i interface{}) bool {
		// Check all items are less than 10
		return i.(int) < 10
	})
```

Output:

```
true
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

### WithPool

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool