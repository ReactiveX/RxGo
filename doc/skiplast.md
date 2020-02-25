# SkipLast Operator

## Overview

Suppress the last n items emitted by an Observable.

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3, 4, 5}).SkipLast(2)
```

Output:

```
1
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