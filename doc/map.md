
# Map Operator

## Overview

Transform the items emitted by an Observable by applying a function to each item.

![](http://reactivex.io/documentation/operators/images/map.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
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

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)

### WithPool

[Detail](options.md#withpool)

### WithCPUPool

[Detail](options.md#withcpupool)

### Serialize

[Detail](options.md#serialize)

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)