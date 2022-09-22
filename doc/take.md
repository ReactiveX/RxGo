# Take

## Overview

Emit only the first n items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/take.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().Take(2)
```

Output:

```
1
2
```
