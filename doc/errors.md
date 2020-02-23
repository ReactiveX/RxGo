# Errors Operator

## Overview

Return the eventual Observable errors.

This method is blocking.

## Example

```go
errs := rxgo.Just([]interface{}{
	errors.New("foo"),
	errors.New("bar"),
	errors.New("baz"),
}).Errors(rxgo.WithErrorStrategy(rxgo.Continue))
fmt.Println(errs)
```

Output:

```
[foo bar baz]
```

## Options

### WithContext

[Detail](options.md#withcontext)