# ToMap Operator

## Overview

Transform the Observable items into a Single emitting a map. It accepts a function that transforms each item into its corresponding key in the map.

## Example

```go
observable := rxgo.Just(1, 2, 3)().
	ToMap(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) * 10, nil
	})
```

Output:

```
map[10:1 20:2 30:3]
```

## Options

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)