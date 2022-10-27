# Interval

> Creates an Observable that emits sequential numbers every specified interval of time.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/interval.png)

`Interval` returns an Observable that emits an infinite sequence of ascending integers, with a constant interval of time of your choosing between those emissions. The first emission is not sent immediately, but only after the first period has passed.

## Example

```go
rxgo.Interval(5 * time.Second).
SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 0 // after 5s
// Next -> 1 // after 10s
// Next -> 2 // after 15s
// Next -> 3 // after 20s
// Next -> 4 // after 25s
// ...
```
