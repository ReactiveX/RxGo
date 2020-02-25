# FlatMap Operator

## Overview

Transform the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.

![](http://reactivex.io/documentation/operators/images/flatMap.c.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().FlatMap(func(i rxgo.Item) rxgo.Observable {
	return rxgo.Just(i.V.(int) * 10, i.V.(int) * 100)()
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

[Detail](options.md#withpool)

### WithCPUPool

[Detail](options.md#withcpupool)

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)