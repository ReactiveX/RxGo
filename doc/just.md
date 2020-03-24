# Just Operator

## Overview

Convert an object or a set of objects into an Observable that emits that or those objects.

![](http://reactivex.io/documentation/operators/images/just.png)

## Examples

### Single Item

```go
observable := rxgo.Just(1)()
```

Output:

```
1
```

### Multiple Items

```go
observable := rxgo.Just(1, 2, 3)()
```

Output:

```
1
2
3
```

### Channel

```go
externalCh := make(chan int)
observable := rxgo.Just(externalCh)()
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithPublishStrategy](options.md#withpublishstrategy)