# Delay

> Delays the emission of items from the source Observable by a given timeout.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/delay.svg)

If the delay argument is a Number, this operator time shifts the source Observable by that amount of time expressed in milliseconds. The relative time intervals between the values are preserved.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint8](1, 5),
    rxgo.Delay[uint8](time.Second),
).SubscribeSync(func(v uint8) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// after 1 second
// Next -> 1
// Next -> 2
// Next -> 3
// Next -> 4
// Next -> 5
// Complete!
```
