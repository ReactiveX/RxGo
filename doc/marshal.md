# Marshal Operator

## Overview

Transform the items emitted by an Observable by applying a marshaller function (`func(interface{}) ([]byte, error)`) to each item.

## Example

```go
type customer struct {
	ID int `json:"id"`
}

observable := rxgo.Just(
	customer{
		ID: 1,
	},
	customer{
		ID: 2,
	},
)().Marshal(json.Marshal)
```

Output:

```
{"id":1}
{"id":2}
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

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)