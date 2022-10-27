# ElementAt

> Emits the single value at the specified index in a sequence of emissions from the source Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/elementAt.png)

`ElementAt` returns an Observable that emits the item at the specified index in the source Observable, or a default value if that index is out of range and the default argument is provided. If the default argument is not given and the index is out of range, the output Observable will emit an `ErrArgumentOutOfRange` error.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 100),
    rxgo.ElementAt[uint](2),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 3
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 10),
    rxgo.ElementAt[uint](88),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    // it will throws `rxgo.ErrArgumentOutOfRange`
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: argument out of range
```
