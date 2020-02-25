# Run Operator

## Overview

Create an Observer without consuming the emitted items.

It returns a `<-chan struct{}` that closes once the Observable terminates.

## Example

```go
<-rxgo.Just(1, 2, errors.New("foo"))().Run()
```

## Options

* [WithContext](options.md#withcontext)