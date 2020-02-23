# Contains Operator

## Overview

Determine whether an Observable emits a particular item or not.

![](http://reactivex.io/documentation/operators/images/contains.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).Contains(func(i interface{}) bool {
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

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)