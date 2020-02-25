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
observable := rxgo.Just(1, 2, 3, 4)().AverageInt()
```

Output:

```
2
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)