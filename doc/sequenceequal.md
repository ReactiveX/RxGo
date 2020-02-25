# SequenceEqual Operator

## Overview

Determine whether two Observables emit the same sequence of items.

![](http://reactivex.io/documentation/operators/images/sequenceEqual.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().
	SequenceEqual(rxgo.Just(1, 2, 42, 4, 5)())
```

Output:

```
false
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