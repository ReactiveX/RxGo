# Materialize

> Represents all of the notifications from the source Observable as next emissions marked with their original types within Notification objects.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/materialize.png)

`Materialize` returns an Observable that emits a next notification for each next, error, or complete emission of the source Observable. When the source Observable emits complete, the output Observable will emit next as a Notification of type "complete", and then it will emit complete as well. When the source Observable emits error, the output will emit next as a Notification of type "error", and then complete.

This operator is useful for producing metadata of the source Observable, to be consumed as next emissions. Use it in conjunction with `Dematerialize`.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.Materialize(),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// ObservableNotification { Kind: 0, Value: 1, Err: nil }
// ObservableNotification { Kind: 0, Value: 2, Err: nil }
// ObservableNotification { Kind: 0, Value: 3, Err: nil }
// ObservableNotification { Kind: 0, Value: 4, Err: nil }
// ObservableNotification { Kind: 0, Value: 5, Err: nil }
// ObservableNotification { Kind: 2, Value: nil, Err: nil }
```
