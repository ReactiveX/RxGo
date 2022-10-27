# Dematerialize

> Converts an Observable of `ObservableNotification` objects into the emissions that they represent.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/dematerialize.png)

`Dematerialize` is assumed to operate an Observable that only emits ObservableNotification objects as next emissions, and does not emit any error. Such Observable is the output of a materialize operation. Those notifications are then unwrapped using the metadata they contain, and emitted as next, error, and complete on the output Observable.

Use this operator in conjunction with `Materialize`.

## Example

```go
rxgo.Pipe1(
    rxgo.Of2[ObservableNotification[string]](rxgo.Next("a"), rxgo.Next("hello"), rxgo.Next("j"), rxgo.Complete()),
    rxgo.Dematerialize(),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> a
// Next -> hello
// Next -> j
// Complete!
```
