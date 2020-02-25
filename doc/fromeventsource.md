# FromEventSource Operator

## Overview

Create a hot observable from a channel.

The items are consumed as soon as the observable is created. An Observer will see only the items since the moment he subscribed to the Observable.

## Example

```go
ch := make(chan rxgo.Item)
observable := rxgo.FromEventSource(ch)
```

## Options

* WithBackPressureStrategy

    * Block (default): block until the Observer is ready to consume the next item using `rxgo.WithBackPressureStrategy(rxgo.Block)`

    * Drop: drop the item if the Observer isn't ready using `rxgo.WithBackPressureStrategy(rxgo.Drop)`

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)