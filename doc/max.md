# Max

> The Max operator operates on an Observable that emits numbers (or items that can be compared with a provided function), and when source Observable completes it emits a single item: the item with the largest value.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/max.png)

## Example

```go
rxgo.Pipe1(
	rxgo.Of2[uint](5, 4, 8, 2, 0),
	rxgo.Max[uint](),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 8
// Complete!
```
