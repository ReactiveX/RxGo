# FindIndex

> Emits only the index of the first value emitted by the source Observable that meets some condition.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/findIndex.png)

**FindIndex** searches for the first item in the source Observable that matches the specified condition embodied by the predicate, and returns the (zero-based) index of the first occurrence in the source. Unlike first, the predicate is required in **FindIndex**, and does not emit an error if a valid value is not found.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.FindIndex(func(v, idx uint) bool {
        return v == 2
    }),
).SubscribeSync(func(v int) {
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
    rxgo.Of2("a", "b", "c", "d", "e"),
    rxgo.FindIndex(func(v string, _ uint) bool {
        return v == "d"
    }),
).SubscribeSync(func(v int) {
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

## Example 3

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 100),
    rxgo.FindIndex(func(v, idx uint) bool {
        return v > 200
    }),
).SubscribeSync(func(v int) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> -1
// Complete!
```
