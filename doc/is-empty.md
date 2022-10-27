# IsEmpty

> Emits false if the input Observable emits any values, or emits true if the input Observable completes without emitting any values.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/isEmpty.png)

`IsEmpty` transforms an Observable that emits values into an Observable that emits a single boolean value representing whether or not any values were emitted by the source Observable. As soon as the source Observable emits a value, `IsEmpty` will emit a false and complete. If the source Observable completes having not emitted anything, `IsEmpty` will emit a true and complete.

A similar effect could be achieved with count, but `IsEmpty` can emit a false value sooner.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 3),
	rxgo.IsEmpty[uint](),
).SubscribeSync(func(v bool) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> false
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Empty[any](),
	rxgo.IsEmpty[any](),
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
