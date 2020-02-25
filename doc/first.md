# First Operator

## Overview

Emit only the first item emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/first.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).First()
```

Output:

```
true
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