# TimeInterval Operator

## Overview

Convert an Observable that emits items into one that emits indications of the amount of time elapsed between those emissions.

![](http://reactivex.io/documentation/operators/images/timeInterval.c.png)

## Example

```go
observable := rxgo.Interval(rxgo.WithDuration(time.Second)).TimeInterval()
```

Output:

```
1.002664s
1.004267s
1.00044s
...
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)