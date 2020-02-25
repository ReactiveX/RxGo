# Last Operator

## Overview

Emit only the last item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/last.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).Last()
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