# IgnoreElements

> Ignores all items emitted by the source Observable and only passes calls of `Complete` or `Error`.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/ignoreElements.png)

The `IgnoreElements` operator suppresses all items emitted by the source Observable, but allows its termination notification (either `Error` or `Complete`) to pass through unchanged.

If you do not care about the items being emitted by an Observable, but you do want to be notified when it completes or when it terminates with an error, you can apply the `IgnoreElements` operator to the Observable, which will ensure that it will never call its observers next handlers.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.IgnoreElements[uint](),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Complete!
```
