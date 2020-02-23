# Timer Operator

## Overview

Create an Observable that emits a single item after a given delay.

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

### WithContext

[Detail](options.md#withcontext)