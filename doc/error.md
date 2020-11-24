# Error Operator

## Overview

Return the first error thrown by an Observable.

This method is blocking.

## Example

```go
err := rxgo.Just(1, 2, errors.New("foo"))().Error()
fmt.Println(err)
```

Output:

```
foo
```

## Options

* [WithContext](options.md#withcontext)