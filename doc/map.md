
# Map Operator

## Overview

Transform the items emitted by an Observable by applying a function to each item.

![](http://reactivex.io/documentation/operators/images/map.png)

## Example

```go
observable := rxgo.Just(1, 2, 3)().
	Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
```

Output:

```
10
20
30
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

### Serialize

[Detail](options.md#serialize)

* [WithPublishStrategy](options.md#withpublishstrategy)