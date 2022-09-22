# Iif

> Checks a boolean at subscription time, and chooses between one of two observable sources

## Description

`Iif` expects a function that returns a boolean (the condition function), and two sources, the `trueResult` and the `falseResult`, and returns an Observable.

At the moment of subscription, the condition function is called. If the result is true, the subscription will be to the source passed as the `trueResult`, otherwise, the subscription will be to the source passed as the `falseResult`.

If you need to check more than two options to choose between more than one observable, have a look at the `Defer` creation method.

## Example

```go
var (
    flag = true
    iif = rxgo.Iif(
        func() bool {
            return flag
        },
        rxgo.Range(1, 10),
        rxgo.Empty[uint](),
    )
)

iif.SubscribeSync(func(v uint) {
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

flag = false
iif.SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Complete!
```
