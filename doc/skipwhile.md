# SkipWhile Operator

## Overview

Suppress the Observable items while a condition is not met.

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().SkipWhile(func(i interface{}) bool {
	return i != 2
})
```

Output:

```
2
3
4
5
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)