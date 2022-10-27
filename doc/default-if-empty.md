# DefaultIfEmpty

> Emits a given value if the source Observable completes without emitting any next value, otherwise mirrors the source Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/defaultIfEmpty.png)

`DefaultIfEmpty` emits the values emitted by the source Observable or a specified default value if the source Observable is empty (completes without having emitted any next value).

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 3),
    rxgo.DefaultIfEmpty[uint](100),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 1
// Next -> 2
// Next -> 3
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Empty[any](),
    rxgo.DefaultIfEmpty[any]("default"),
).SubscribeSync(func(v any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> default
// Complete!
```
