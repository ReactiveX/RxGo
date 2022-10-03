# Filter

> Filter items emitted by the source Observable by only emitting those that satisfy a specified predicate.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/filter.png)

Takes values from the source Observable, passes them through a predicate function and only emits those values that yielded true.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.Filter(func(v uint, index uint) bool {
		return v != 2
	}),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 1
// Next -> 3
// Next -> 4
// Next -> 5
// Complete!
```
