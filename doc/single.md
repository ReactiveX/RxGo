# Single

> Returns an observable that asserts that only one value is emitted from the observable that matches the predicate. If no predicate is provided, then it will assert that the observable only emits one value.

## Description

In the event that the observable is empty, it will throw an `ErrEmpty`.

In the event that two values are found that match the predicate, or when there are two values emitted and no predicate, it will throw a `ErrSequence`.

In the event that no values match the predicate, if one is provided, it will throw a `ErrNotFound`.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 10),
	rxgo.Single(func(value, index uint, source rxgo.Observable[uint]) bool {
		return value == 2
	}),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 2
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Empty[string](),
    rxgo.Single[string](),
).SubscribeSync(func(v string) {
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
rxgo.Pipe1(
    rxgo.Range[uint](1, 10),
	rxgo.Single(func(value, index uint, source rxgo.Observable[uint]) bool {
        return value > 2
    }),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: too many values match (ErrSequence)
```

## Example 4

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 10),
	rxgo.Single(func(value, index uint, source rxgo.Observable[uint]) bool {
		return value > 100
	}),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: no values match (ErrNotFound)
```
