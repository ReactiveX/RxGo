# Repeat

> Returns an Observable that will resubscribe to the source stream when the source stream completes.

## Description

![](http://reactivex.io/documentation/operators/images/repeat.png)

`Repeat` will output values from a source until the source completes, then it will resubscribe to the source a specified number of times, with a specified delay. Repeat can be particularly useful in combination with closing operators like `Take`, `TakeUntil`, `First`, or `TakeWhile`, as it can be used to restart a source again from scratch.

Repeat is very similar to retry, where retry will resubscribe to the source in the error case, but repeat will resubscribe if the source completes.

Note that `Repeat` will not catch errors. Use `Retry` for that.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 3),
    rxgo.Repeat[uint, uint8](2),
).SubscribeSync(func(v []uint) {
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
// Next -> 1
// Next -> 2
// Next -> 3
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Range[uint8](2, 4),
    rxgo.Repeat[uint8](rxgo.RepeatConfig{
        Count: 2,
        Delay: time.Second,
    }),
).SubscribeSync(func(v uint8) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 2
// Next -> 3
// Next -> 4
// Next -> 5
// Next -> 2
// Next -> 3
// Next -> 4
// Next -> 5
// Complete!
```
