# Join Operator

## Overview

Combine items emitted by two Observables whenever an item from one Observable is emitted during a time window defined according to an item emitted by the other Observable. 

The time is extracted using a timeExtractor function.

![](http://reactivex.io/documentation/operators/images/join.c.png)

## Example

```go
observable := rxgo.Just(
	map[string]int64{"tt": 1, "V": 1},
	map[string]int64{"tt": 4, "V": 2},
	map[string]int64{"tt": 7, "V": 3},
)().Join(func(ctx context.Context, l interface{}, r interface{}) (interface{}, error) {
	return map[string]interface{}{
		"l": l,
		"r": r,
	}, nil
}, rxgo.Just(
	map[string]int64{"tt": 2, "V": 5},
	map[string]int64{"tt": 3, "V": 6},
	map[string]int64{"tt": 5, "V": 7},
)(), func(i interface{}) time.Time {
	return time.Unix(i.(map[string]int64)["tt"], 0)
}, rxgo.WithDuration(2))
```

Output:

```
map[l:map[V:1 tt:1] r:map[V:5 tt:2]]
map[l:map[V:1 tt:1] r:map[V:6 tt:3]]
map[l:map[V:2 tt:4] r:map[V:5 tt:2]]
map[l:map[V:2 tt:4] r:map[V:6 tt:3]]
map[l:map[V:2 tt:4] r:map[V:7 tt:5]]
map[l:map[V:3 tt:7] r:map[V:7 tt:5]]
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)