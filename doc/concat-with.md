# ConcatWith

> Emits all of the values from the source observable, then, once it completes, subscribes to each observable source provided, one at a time, emitting all of their values, and not subscribing to the next one until it completes.

## Example

```go
rxgo.Pipe1(
    rxgo.Of2[any]("a", "b", "c", "d"),
    rxgo.ConcatWith(
        rxgo.Of2[any](1, 2, 88),
        rxgo.Of2[any](88.1991, true, false),
    ),
).SubscribeSync(func(v any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> a
// Next -> b
// Next -> c
// Next -> d
// Next -> 1
// Next -> 2
// Next -> 88
// Next -> 88.1991
// Next -> true
// Next -> false
// Complete!
```
