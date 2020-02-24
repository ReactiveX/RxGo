# Max Operator

## Overview

Determine, and emit, the maximum-valued item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/max.png)

## Example

```go
observable := rxgo.Just([]interface{}{2, 5, 1, 6, 3, 4}).
	Max(func(i1 interface{}, i2 interface{}) int {
		return i1.(int) - i2.(int)
	})
```

Output:

```
6
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

[Detail](options.md#withpool)

### WithCPUPool

[Detail](options.md#withcpupool)