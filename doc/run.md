# Run Operator

## Overview

Create an Observer without consuming the emitted items.

It returns a `<-chan struct{}` that closes once the Observable terminates.

## Example

```go
<-rxgo.Just([]interface{}{1, 2, errors.New("foo")}).Run()
```

## Options

### WithContext

[Detail](options.md#withcontext)