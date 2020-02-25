# SkipWhile Operator

## Overview

Suppress the Observable items while a condition is not met.

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).SkipWhile(func(i interface{}) bool {
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