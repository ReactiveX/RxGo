# ConcatMap

> Projects each source value to an Observable which is merged in the output Observable, in a serialized fashion waiting for each one to complete before merging the next.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/concatMap.png)

Returns an Observable that emits items based on applying a function that you supply to each item emitted by the source Observable, where that function returns an (so-called "inner") Observable. Each new inner Observable is concatenated with the previous inner Observable.

Warning: if source values arrive endlessly and faster than their corresponding inner Observables can complete, it will result in memory issues as inner Observables amass in an unbounded buffer waiting for their turn to be subscribed to.

_Note: **ConcatMap** is equivalent to **MergeMap** with concurrency parameter set to 1._

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.ConcatMap(func(x, i uint) rxgo.Observable[string] {
        return rxgo.Pipe2(
            rxgo.Interval(time.Second),
            rxgo.Map(func(y, _ uint) (string, error) {
                return fmt.Sprintf("%v[%d]", x, y), nil
            }),
            rxgo.Take[string](2),
        )
    }),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 1[0] // 1s
// Next -> 1[1] // 2s
// Next -> 2[0] // 3s
// Next -> 2[1] // 4s
// Next -> 3[0] // 5s
// Next -> 3[1] // 6s
// Next -> 4[0] // 7s
// Next -> 4[1] // 8s
// Next -> 5[0] // 9s
// Next -> 5[1] // 10s
// Complete!
```
