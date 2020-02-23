# GroupBy Operator

## Overview

Divide an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key.

![](http://reactivex.io/documentation/operators/images/groupBy.c.png)

## Example

```go
count := 3
observable := rxgo.Range(0, 10).GroupBy(count, func(item rxgo.Item) int {
	return item.V.(int) % count
}, rxgo.WithBufferedChannel(10))

for i := range observable.Observe() {
	fmt.Println("New observable:")

	for i := range i.V.(rxgo.Observable).Observe() {
		fmt.Printf("item: %v\n", i.V)
	}
}
```

Output:

```
New observable:
item: 0
item: 3
item: 6
item: 9
New observable:
item: 1
item: 4
item: 7
item: 10
New observable:
item: 2
item: 5
item: 8
```

## Options

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)