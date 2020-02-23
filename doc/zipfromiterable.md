# ZipFromIterable Operator

## Overview

Merge the emissions of an Iterable via a specified function and emit single items for each combination based on the results of this function.

![](http://reactivex.io/documentation/operators/images/zip.o.png)

## Example

```go
observable1 := rxgo.Just([]interface{}{1, 2, 3})
observable2 := rxgo.Just([]interface{}{10, 20, 30})
zipper := func(_ context.Context, i1 interface{}, i2 interface{}) (interface{}, error) {
	return i1.(int) + i2.(int), nil
}
zippedObservable := observable1.ZipFromIterable(observable2, zipper)
```

Output:

```
11
12
13
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