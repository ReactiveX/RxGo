# ElementAt Operator

## Overview

Emit only item n emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/elementAt.png)

## Example

```go
observable := rxgo.Just(0, 1, 2, 3, 4)().ElementAt(2)
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

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)