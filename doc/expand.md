# Expand

> Recursively projects each source value to an Observable which is merged in the output Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/expand.png)

Returns an Observable that emits items based on applying a function that you supply to each item emitted by the source Observable, where that function returns an Observable, and then merging those resulting Observables and emitting the results of this merger. **Expand** will re-emit on the output Observable every source value. Then, each output value is given to the project function which returns an inner Observable to be merged on the output Observable. Those output values resulting from the projection are also given to the project function to produce new output values. This is how expand behaves recursively.

## Example 1

```go
rxgo.Pipe2(
    rxgo.Range[uint8](1, 5),
    rxgo.Expand(func(v uint8, _ uint) rxgo.Observable[string] {
        return rxgo.Of2(fmt.Sprintf("Number(%d)", v))
    }),
    rxgo.Take[Either[uint8, string]](5),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 2
// Complete!
```
