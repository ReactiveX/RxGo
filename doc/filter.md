# Filter Operator

## Overview

Emit only those items from an Observable that pass a predicate test.

![](http://reactivex.io/documentation/operators/images/filter.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().
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

[Detail](options.md#withpool)

### WithCPUPool

[Detail](options.md#withcpupool)

### Serialize

[Detail](options.md#serialize)

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)