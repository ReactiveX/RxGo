# Range Operator

## Overview

Create an Observable that emits a range of sequential integers.

![](http://reactivex.io/documentation/operators/images/range.png)

## Example

```go
observable := rxgo.Range(0, 3)
```

Output:

```
0
1
2
3
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)