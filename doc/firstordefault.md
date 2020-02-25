# FirstOrDefault Operator

## Overview

Similar to `First`, but we pass a default item that will be emitted if the source Observable fails to emit any items.

![](http://reactivex.io/documentation/operators/images/firstOrDefault.png)

## Example

```go
observable := rxgo.Empty().FirstOrDefault(1)
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

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)