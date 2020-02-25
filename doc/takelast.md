# TakeLast Operator

## Overview

Emit only the final n items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/takeLast.n.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().TakeLast(2)
```

Output:

```
4
5
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

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)