# Timestamp Operator

## Overview

Attach a timestamp to each item emitted by an Observable.

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

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)