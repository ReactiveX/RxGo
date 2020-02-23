# Debounce Operator

## Overview

Only emit an item from an Observable if a particular timespan has passed without it emitting another item.

![](http://reactivex.io/documentation/operators/images/debounce.png)

## Example

```go
observable.Debounce(rxgo.WithDuration(250 * time.Millisecond))
```

Output: each item emitted by the Observable if not item has been emitted after 250 milliseconds. 

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

https://github.com/ReactiveX/RxGo/wiki/Options#withpool

### WithCPUPool

https://github.com/ReactiveX/RxGo/wiki/Options#withcpupool