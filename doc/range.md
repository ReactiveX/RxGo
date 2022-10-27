# Range

> Creates an Observable that emits a sequence of numbers within a specified range.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/range.png)

**Range** operator emits a range of sequential integers, in order, where you select the start of the range and its length.

## Example

```go
rxgo.Range[uint8](1, 10).
SubscribeSync(func(v uint8) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 1
// Next -> 2
// Next -> 3
// Next -> 4
// Next -> 5
// Next -> 6
// Next -> 7
// Next -> 8
// Next -> 9
// Next -> 10
// Complete!
```
