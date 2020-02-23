# FlatMap Operator

## Overview

Transform the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.

![](http://reactivex.io/documentation/operators/images/flatMap.c.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).FlatMap(func(i rxgo.Item) rxgo.Observable {
	return rxgo.Just([]interface{}{
		i.V.(int) * 10,
		i.V.(int) * 100,
	})
})
```

Output:

```
10
100
20
200
30
300
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