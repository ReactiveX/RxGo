# Find

> Emits only the first value emitted by the source Observable that meets some condition.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/find.png)

**Find** searches for the first item in the source Observable that matches the specified condition embodied by the predicate, and returns the first occurrence in the source. Unlike first, the predicate is required in **Find**, and does not emit an error if a valid value is not found (emits undefined instead).

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.Find(func(v, idx uint) bool {
        return v == 2
    }),
).SubscribeSync(func(v rxgo.Optional[uint]) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> Some[uint](2)
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Of2("a", "b", "c", "d", "e"),
    rxgo.Find(func(v string, _ uint) bool {
        return v == "d"
    }),
).SubscribeSync(func(v rxgo.Optional[string]) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> Some[string]("d")
// Complete!
```

## Example 3

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 100),
    rxgo.Find(func(v, idx uint) bool {
        return v > 200
    }),
).SubscribeSync(func(v rxgo.Optional[uint]) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> None[uint]()
// Complete!
```
