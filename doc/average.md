# Average Operator

## Overview

Calculate the average of numbers emitted by an Observable and emits this average.

![](http://reactivex.io/documentation/operators/images/average.png)

## Instances

* `AverageFloat32`
* `AverageFloat64`
* `AverageInt`
* `AverageInt8`
* `AverageInt16`
* `AverageInt32`
* `AverageInt64`

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4}).AverageInt()
```

Output:

```
2
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