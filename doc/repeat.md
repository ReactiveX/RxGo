# Repeat Operator

## Overview

Create an Observable that emits a particular item multiple times at a particular frequency.

![](http://reactivex.io/documentation/operators/images/repeat.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, 3}).
	Repeat(3, rxgo.WithDuration(time.Second))
```

Output:

```
// Immediately
1
2
3
// After 1 second
1
2
3
// After 2 seconds
1
2
3
...
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