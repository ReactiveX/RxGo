# Count

> Counts the number of emissions on the source and emits that number when the source completes.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/count.png)

`Count` transforms an Observable that emits values into an Observable that emits a single value that represents the number of values emitted by the source Observable. If the source Observable terminates with an error, count will pass this error notification along without emitting a value first. If the source Observable does not terminate at all, count will neither emit a value nor terminate. This operator takes an optional predicate function as argument, in which case the output emission will represent the number of source values that matched true with the predicate.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](0, 7),
    rxgo.Count[uint](),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 7
// Complete!
```
