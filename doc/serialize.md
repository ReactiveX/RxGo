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

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)