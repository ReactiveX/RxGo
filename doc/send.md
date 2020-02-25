# Send Operator

## Overview

Send the Observable items to a given channel that will closed once the operation completes.

## Example

```go
ch := make(chan rxgo.Item)
rxgo.Just(1, 2, 3)().Send(ch)
for item := range ch {
	fmt.Println(item.V)
}
```

Output:

```
1
2
3
```

## Options

* [WithContext](options.md#withcontext)