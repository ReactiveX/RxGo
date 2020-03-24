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

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)