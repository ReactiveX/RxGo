# Contains Operator

## Overview

Determine whether an Observable emits a particular item or not.

![](http://reactivex.io/documentation/operators/images/contains.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().Contains(func(i interface{}) bool {
	return i == 2
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

### WithPool

[Detail](options.md#withpool)

### WithCPUPool

[Detail](options.md#withcpupool)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)