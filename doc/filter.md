# Filter Operator

## Overview

Emit only those items from an Observable that pass a predicate test.

![](http://reactivex.io/documentation/operators/images/filter.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
	Filter(func(i interface{}) bool {
		return i != 2
	})
```

Output:

```
1
3
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