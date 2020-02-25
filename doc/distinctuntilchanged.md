# DistinctUntilChanged Operator

## Overview

Suppress consecutive duplicate items in the original Observable.

## Example

```go
observable := rxgo.Just(1, 2, 2, 1, 1, 3)().
	DistinctUntilChanged(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	})
```

Output:

```
1
2
1
3
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)