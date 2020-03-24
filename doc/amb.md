# Amb Operator

## Overview

Given two or more source Observables, emit all of the items from only the first of these Observables to emit an item.

![](http://reactivex.io/documentation/operators/images/amb.png)

## Example

```go
observable := rxgo.Amb([]rxgo.Observable{
	rxgo.Just(1, 2, 3)(),
	rxgo.Just(4, 5, 6)(),
})
```

Output:

```
1
2
3
```
or
```
4
5
6
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)