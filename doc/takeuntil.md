# TakeUntil Operator

## Overview

Discard any items emitted by an Observable after a second Observable emits an item or terminates.

![](http://reactivex.io/documentation/operators/images/takeUntil.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().TakeUntil(func(i interface{}) bool {
	return i == 3
})
```

Output:

```
1
2
3
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