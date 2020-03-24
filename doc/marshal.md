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

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPool](options.md#withpool)

* [WithCPUPool](options.md#withcpupool)

* [WithPublishStrategy](options.md#withpublishstrategy)