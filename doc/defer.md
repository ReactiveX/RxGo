# Defer Operator

## Overview

do not create the Observable until the observer subscribes, and create a fresh Observable for each observer.

![](http://reactivex.io/documentation/operators/images/defer.png)

## Example

```go
observable := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item, done func()) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Of(3)
	done()
}})
```

Output:

```
1
2
3
```

There are two ways to close the Observable:
* Closing the `next` channel.
* Calling the `done()` function.

Yet, as we can pass multiple producers, using the `done()` function is the recommended approach.

## Options

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)