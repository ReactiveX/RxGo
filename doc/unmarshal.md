# Unmarshal Operator

## Overview

Transform the items emitted by an Observable by applying an unmarshaller function (`func([]byte, interface{}) error`) to each item. It takes a factory function that initializes the target structure.

## Example

```go
observable := rxgo.Just(
	[]byte(`{"id":1}`),
	[]byte(`{"id":2}`),
)().Unmarshal(json.Unmarshal,
	func() interface{} {
		return &customer{}
	})
```

Output:

```
&{ID:1}
&{ID:2}
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