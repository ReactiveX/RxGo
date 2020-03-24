# SkipLast Operator

## Overview

Suppress the last n items emitted by an Observable.

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().SkipLast(2)
```

Output:

```
1
2
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)