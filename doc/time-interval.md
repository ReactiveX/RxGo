# TimeInterval

> Emits an object containing the current value, and the time that has passed between emitting the current value and the previous value, which is calculated by using the provided scheduler's `time.Now()` method to retrieve the current time at each emission, then calculating the difference.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/timeInterval.png)

Convert an Observable that emits items into one that emits indications of the amount of time elapsed between those emissions.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Interval(time.Second),
    rxgo.WithTimeInterval[uint](),
).SubscribeSync(func(t rxgo.TimeInterval[uint]) {
    log.Println("Next ->", t.Value(), "|", t.Elapsed())
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 0 | 1s
// Next -> 1 | 1s
// Next -> 2 | 1s
// Next -> 3 | 1s
// Next -> 4 | 1s
// Next -> 5 | 1s
// ...
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.WithTimeInterval[uint](),
).SubscribeSync(func(t rxgo.TimeInterval[uint]) {
    log.Println("Next ->", t.Value(), "|", t.Elapsed())
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 1 | 1.2ms
// Next -> 2 | 5ms
// Next -> 3 | 3ms
// Next -> 4 | 8ms
// Next -> 5 | 1.1ms
// Complete!
```
