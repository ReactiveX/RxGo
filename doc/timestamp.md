# Timestamp Operator

## Overview

Determine whether all items emitted by an Observable meet some criteria.

![](http://reactivex.io/documentation/operators/images/timestamp.c.png)

## Example

```go
observe := rxgo.Just(1, 2, 3)().Timestamp().Observe()
var timestampItem rxgo.TimestampItem
timestampItem = (<-observe).V.(rxgo.TimestampItem)
fmt.Println(timestampItem)
```

Output:

```
{2020-02-23 15:26:02.231197 +0000 UTC 1}
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