# Timer Operator

## Overview

Create an Observable that completes after a specified delay.

![](http://reactivex.io/documentation/operators/images/timer.png)

## Example

```go
observable := rxgo.Timer(rxgo.WithDuration(5 * time.Second))
```

Output:

```
{} // After 5 seconds
```

## Options

* [WithContext](options.md#withcontext)