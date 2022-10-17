# ThrowIfEmpty

> If the source observable completes without emitting a value, it will emit an error. The error will be created at that time by the optional errorFactory argument, otherwise, the error will be `ErrEmpty`.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/throwIfEmpty.png)

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 3),
	rxgo.ThrowIfEmpty[uint](),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 1
// Next -> 2
// Next -> 3
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Empty[any](),
    rxgo.ThrowIfEmpty[any](),
).SubscribeSync(func(v any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: empty value (ErrEmpty)
```

## Example 3

```go
var err = errors.New("something wrong")

rxgo.Pipe1(
    rxgo.Empty[any](),
    rxgo.ThrowIfEmpty[any](func() error {
        return err
    }),
).SubscribeSync(func(v any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> something wrong
```
