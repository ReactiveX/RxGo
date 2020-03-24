# Buffer Operator

## Overview

Periodically gather items emitted by an Observable into bundles and emit these bundles rather than emitting the items one at a time

![](http://reactivex.io/documentation/operators/images/Buffer.png)

## Instances

* `BufferWithCount`:

![](http://reactivex.io/documentation/operators/images/bufferWithCount3.png)

```go
observable := rxgo.Just(1, 2, 3, 4)().BufferWithCount(3)
```

Output:

```
1 2 3
4
```

* `BufferWithTime`:

![](http://reactivex.io/documentation/operators/images/bufferWithTime5.png)

```go
// Create the producer
ch := make(chan rxgo.Item, 1)
go func() {
	i := 0
	for range time.Tick(time.Second) {
		ch <- rxgo.Of(i)
		i++
	}
}()

observable := rxgo.FromChannel(ch).
	BufferWithTime(rxgo.WithDuration(3*time.Second), nil)
```

Output:

```
0 1 2
3 4 5
6 7 8
...
```

* `BufferWithTimeOrCount`:

![](http://reactivex.io/documentation/operators/images/bufferWithTimeOrCount6.png)

```go
// Create the producer
ch := make(chan rxgo.Item, 1)
go func() {
	i := 0
	for range time.Tick(time.Second) {
		ch <- rxgo.Of(i)
		i++
	}
}()

observable := rxgo.FromChannel(ch).
	BufferWithTimeOrCount(rxgo.WithDuration(3*time.Second), 2)
```

Output:

```
0 1
2 3
4 5
...
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)