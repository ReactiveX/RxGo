# IgnoreElements Operator

## Overview

Do not emit any items from an Observable but mirror its termination notification.

![](http://reactivex.io/documentation/operators/images/ignoreElements.c.png)

## Example

```go
observable := rxgo.Just(1, 2, errors.New("foo"))().
	IgnoreElements()
```

Output:

```
foo
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