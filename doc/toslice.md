# ToSlice Operator

## Overview

Transform the Observable items into a slice. It accepts a capacity that will be used as the initial capacity of the slice produced.

## Example

```go
s, err := rxgo.Just(1, 2, 3)().ToSlice(3)
if err != nil {
	return err
}
fmt.Println(s)
```

Output:

```
[1 2 3]
```

## Options

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)