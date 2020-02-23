# StartWithIterable Operator

## Overview

Emit a specified Iterable before beginning to emit the items from the source Observable.

![](http://reactivex.io/documentation/operators/images/startWith.png)

## Example

```go
observable := rxgo.Just([]interface{}{3, 4}).StartWith(
	rxgo.Just([]interface{}{1, 2}))
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