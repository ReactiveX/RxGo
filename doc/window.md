# Window Operator

## Overview

Periodically subdivide items from an Observable into Observable windows and emit these windows rather than emitting the items one at a time.

## Instances

* `WindowWithCount`

![](http://reactivex.io/documentation/operators/images/window3.png)

* `WindowWithTime`

![](http://reactivex.io/documentation/operators/images/window5.png)

* `WindowWithTimeOrCount`

![](http://reactivex.io/documentation/operators/images/window6.png)

## Example

```go
observe := rxgo.Just(1, 2, 3)().WindowWithCount(2).Observe()

fmt.Println("First Observable")
for item := range (<-observe).V.(rxgo.Observable).Observe() {
	if item.Error() {
		return item.E
	}
	fmt.Println(item.V)
}

fmt.Println("Second Observable")
for item := range (<-observe).V.(rxgo.Observable).Observe() {
	if item.Error() {
		return item.E
	}
	fmt.Println(item.V)
}
```

Output:

```
First Observable
1
2
Second Observable
3
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)