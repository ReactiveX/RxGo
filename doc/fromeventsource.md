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

### WithBackPressureStrategy

* Block (default): block until the Observer is ready to consume the next item.

```go
rxgo.FromEventSource(ch, rxgo.WithBackPressureStrategy(rxgo.Block))
```

* Drop: drop the item if the Observer isn't ready.

```go
rxgo.FromEventSource(ch, rxgo.WithBackPressureStrategy(rxgo.Drop))
```

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)