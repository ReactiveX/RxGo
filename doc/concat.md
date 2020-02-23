# Concat Operator

## Overview

Emit the emissions from two or more Observables without interleaving them.

![](http://reactivex.io/documentation/operators/images/concat.png)

## Example

```go
observable := rxgo.Concat([]rxgo.Observable{
	rxgo.Just([]interface{}{1, 2, 3}),
	rxgo.Just([]interface{}{4, 5, 6}),
})
```

Output:

```
1
2
3
4
5
6
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