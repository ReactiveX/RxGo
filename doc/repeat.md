# Repeat

## Overview

Create an Observable that emits a particular item multiple times at a particular frequency.

![](http://reactivex.io/documentation/operators/images/repeat.png)

## Example

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
