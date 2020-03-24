# Error Operator

## Overview

Return the eventual Observable error. 

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