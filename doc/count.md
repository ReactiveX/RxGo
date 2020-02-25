# Count Operator

## Overview

Count the number of items emitted by the source Observable and emit only this value.

![](http://reactivex.io/documentation/operators/images/Count.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).Count()
```

Output:

```
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