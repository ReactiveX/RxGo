# Empty

> A simple Observable that emits no items to the Observer and immediately emits a complete notification.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/empty.png)

A simple Observable that only emits the complete notification. It can be used for composing with other Observables, such as in a `MergeMap`.

## Example

```go
rxgo.Empty[any]().
SubscribeSync(func(v any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Complete!
```
