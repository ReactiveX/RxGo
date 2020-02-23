# LastOrDefault Operator

## Overview

Similar to `Last`, but you pass it a default item that it can emit if the source Observable fails to emit any items.

![](http://reactivex.io/documentation/operators/images/lastOrDefault.png)

## Example

```go
observable := rxgo.Empty().LastOrDefault(1)
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