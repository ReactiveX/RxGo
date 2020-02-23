# DistinctUntilChanged Operator

## Overview

Suppress consecutive duplicate items in the original Observable.

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 2, 1, 1, 3}).
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

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)