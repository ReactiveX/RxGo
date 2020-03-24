# Interval Operator

## Overview

Create an Observable that emits a sequence of integers spaced by a particular time interval.

![](http://reactivex.io/documentation/operators/images/interval.png)

## Example

```go
observable := rxgo.Interval(rxgo.WithDuration(5 * time.Second))
```

Output:

```
0 // After 5 seconds
1 // After 10 seconds
2 // After 15 seconds
3 // After 20 seconds
...
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)