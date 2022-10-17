# TakeLast

> Waits for the source to complete, then emits the last N values from the source, as specified by the count argument.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/takeLast.png)

`TakeLast` results in an observable that will hold values up to count values in memory, until the source completes. It then pushes all values in memory to the consumer, in the order they were received from the source, then notifies the consumer that it is complete.

If for some reason the source completes before the count supplied to `TakeLast` is reached, all values received until that point are emitted, and then completion is notified.

Warning: Using `TakeLast` with an observable that never completes will result in an observable that never emits a value.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 100),
	rxgo.TakeLast(3),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 98
// Next -> 99
// Next -> 100
// Complete!
```
