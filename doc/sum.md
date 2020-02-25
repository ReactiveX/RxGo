# Sum Operator

## Overview

Calculate the sum of numbers emitted by an Observable and emit this sum.

![](http://reactivex.io/documentation/operators/images/sum.f.png)

## Instances

* `SumFloat32`
* `SumFloat64`
* `SumInt64`

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4}).SumInt64()
```

Output:

```
10
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