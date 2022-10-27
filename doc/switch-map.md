# SwitchMap

> Projects each source value to an Observable which is merged in the output Observable, emitting values only from the most recently projected Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/switchMap.png)

Returns an Observable that emits items based on applying a function that you supply to each item emitted by the source Observable, where that function returns an (so-called "inner") Observable. Each time it observes one of these inner Observables, the output Observable begins emitting the items emitted by that inner Observable. When a new inner Observable is emitted, switchMap stops emitting items from the earlier-emitted inner Observable and begins emitting items from the new one. It continues to behave like this for subsequent inner Observables.

## Example

```go
rxgo.Pipe1(
    rxgo.Interval(time.Second),
    rxgo.SwitchMap(func(v, _ uint) rxgo.Observable[uint] {
        return rxgo.Interval(time.Millisecond * 500)
    }),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 0 // after 1.5s
// Next -> 0 // after 1s
// Next -> 0 // after 1s
// ...
```
