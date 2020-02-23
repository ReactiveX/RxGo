# Merge Operator

## Overview

Combine multiple Observables into one by merging their emissions.

![](http://reactivex.io/documentation/operators/images/merge.png)

## Example

```go
observable := rxgo.Merge([]rxgo.Observable{
	rxgo.Just([]interface{}{1, 2}),
	rxgo.Just([]interface{}{3, 4}),
})
```

Output:

```
1
2
3
4
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