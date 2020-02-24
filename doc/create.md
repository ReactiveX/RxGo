# Create Operator

## Overview

Create an Observable from scratch by calling observer methods programmatically.

![](http://reactivex.io/documentation/operators/images/create.png)

## Example

```go
observable := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Of(3)
}})
```

Output:

```
1
2
3
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