# Repeat

## Overview

Create an Observable that emits a particular item multiple times at a particular frequency.

![](http://reactivex.io/documentation/operators/images/repeat.png)

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
