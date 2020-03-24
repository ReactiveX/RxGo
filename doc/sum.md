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
observable := rxgo.Just(1, 2, 3, 4)().SumInt64()
```

Output:

```
10
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)