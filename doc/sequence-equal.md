# SequenceEqual

> Compares all values of two observables in sequence using an optional comparator function and returns an observable of a single boolean value representing whether or not the two sequences are equal.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/sequenceEqual.png)

`SequenceEqual` subscribes to two observables and buffers incoming values from each observable. Whenever either observable emits a value, the value is buffered and the buffers are shifted and compared from the bottom up; If any value pair doesn't match, the returned observable will emit false and complete. If one of the observables completes, the operator will wait for the other observable to complete; If the other observable emits before completing, the returned observable will emit false and complete. If one observable never completes or emits after the other completes, the returned observable will never complete.

## Example

```go
rxgo.Pipe1(
	rxgo.Range[uint](1, 10),
	rxgo.SequenceEqual(Range[uint](1, 10)),
).SubscribeSync(func(v bool) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> true
// Complete!
```
