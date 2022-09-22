# ToSlice Operator

> Collects all source emissions and emits them as an slice when the source completes.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/toArray.png)

ToSlice will wait until the source Observable completes before emitting the slice containing all emissions. When the source Observable errors no slice will be emitted.

## Example

```go
rxgo.Pipe2(
	Interval[uint](time.Second),
	Take[uint](10),
	ToSlice[uint](),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
// Complete!
```
