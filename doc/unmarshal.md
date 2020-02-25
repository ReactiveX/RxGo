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

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)