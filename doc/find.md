# Find Operator

## Overview

Emit the first item passing a predicate then complete.

## Example

```go
observable := rxgo.Just(1, 2, 3)().Find(func(i interface{}) bool {
    return i == 2
})
```

Output:

```
2
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)