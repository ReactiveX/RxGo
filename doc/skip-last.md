# SkipLast

> Skip a specified number of values before the completion of an observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/skipLast.png)

Returns an observable that will emit values as soon as it can, given a number of skipped values. For example, if you `SkipLast(3)` on a source, when the source emits its fourth value, the first value the source emitted will finally be emitted from the returned observable, as it is no longer part of what needs to be skipped.

All values emitted by the result of `SkipLast(N)` will be delayed by N emissions, as each value is held in a buffer until enough values have been emitted that that the buffered value may finally be sent to the consumer.

After subscribing, unsubscribing will not result in the emission of the buffered skipped values.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](0, 10),
	rxgo.SkipLast(3),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 0
// Next -> 1
// Next -> 2
// Next -> 3
// Next -> 4
// Next -> 5
// Next -> 6
// Next -> 7 (8, 9 and 10 are skipped)
// Complete!
```
