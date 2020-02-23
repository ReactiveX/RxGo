# DefaultIfEmpty Operator

## Overview

Emit items from the source Observable, or a default item if the source Observable emits nothing.

![](http://reactivex.io/documentation/operators/images/defaultIfEmpty.c.png)

## Example

```go
observable := rxgo.Empty().DefaultIfEmpty(1)
```

Output:

```
1
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