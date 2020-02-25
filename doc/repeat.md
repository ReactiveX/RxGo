# Repeat Operator

## Overview

Create an Observable that emits a particular item multiple times at a particular frequency.

![](http://reactivex.io/documentation/operators/images/repeat.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().
	Repeat(3, rxgo.WithDuration(time.Second))
```

Output:

```
// Immediately
1
2
3
// After 1 second
1
2
3
// After 2 seconds
1
2
3
...
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)