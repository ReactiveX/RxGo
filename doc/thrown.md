# Thrown Operator

## Overview

Create an Observable that emits no items and terminates with an error.

![](http://reactivex.io/documentation/operators/images/throw.c.png)

## Example

```go
observable := rxgo.Thrown(errors.New("foo"))
```

Output:

```
foo
```