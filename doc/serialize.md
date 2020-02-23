# Serialize Operator

## Overview

Force an Observable to make serialized calls and to be well-behaved. 
It takes the starting index and a function that transforms an item value into an index. 

![](http://reactivex.io/documentation/operators/images/serialize.c.png)

## Example

```go
observable := rxgo.Range(0, 1_000_000, rxgo.WithBufferedChannel(capacity)).
	Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}, rxgo.WithCPUPool(), rxgo.WithBufferedChannel(capacity)).
	Serialize(0, func(i interface{}) int {
		return i.(int)
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