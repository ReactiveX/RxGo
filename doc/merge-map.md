# MergeMap

> Projects each source value to an Observable which is merged in the output Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/mergeMap.png)

Returns an Observable that emits items based on applying a function that you supply to each item emitted by the source Observable, where that function returns an Observable, and then merging those resulting Observables and emitting the results of this merger.

# Example

```go
rxgo.Pipe1(
    rxgo.Of2("a", "b", "c"),
    rxgo.MergeMap(func(x string, _ uint) rxgo.Observable[string] {
        return rxgo.Pipe1(
            rxgo.Interval(time.Second),
            rxgo.Map(func(y, _ uint) (string, error) {
                return fmt.Sprintf("%s%d", x, y), nil
            }),
        )
    }),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> b0
// Next -> a0
// Next -> c0
// Next -> b1
// Next -> c1
// Next -> a1
// Next -> a2
// Next -> c2
// Next -> b2
// ...
```
