# First

> Emits only the first value (or the first value that meets some condition) emitted by the source Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/first.png)

If called with no arguments, `First` emits the first value of the source Observable, then completes. If called with a predicate function, `First` emits the first value of the source that matches the specified condition. Throws an error if `defaultValue` was not provided and a matching element is not found.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.First[uint](nil),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 1
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Empty[string](),
    rxgo.First[string](nil, "defaultValue"),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> defaultValue
// Complete!
```

## Example 3

```go
rxgo.Pipe1(
    rxgo.Empty[string](),
    rxgo.First[string](nil),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: empty value
```

## Example 4

```go
rxgo.Pipe1(
    rxgo.Range[uint8](1, 10),
    rxgo.First[uint8](func(v uint8, _ uint) bool {
        return v > 50
    }),
).SubscribeSync(func(v uint8) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: empty value
```
