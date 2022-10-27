# Never

> An Observable that emits no items to the Observer and never completes.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/never.png)

A simple Observable that emits neither values nor errors nor the completion notification. It can be used for testing purposes or for composing with other Observables. Please note that by never emitting a complete notification, this Observable keeps the subscription from being disposed automatically. Subscriptions need to be manually disposed.

## Example

```go
rxgo.Never[any]().
SubscribeSync(func(v any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// ... never complete
```
