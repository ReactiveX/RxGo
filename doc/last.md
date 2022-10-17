# Last

> Returns an Observable that emits only the last item emitted by the source Observable. It optionally takes a predicate function as a parameter, in which case, rather than emitting the last item from the source Observable, the resulting Observable will emit the last item from the source Observable that satisfies the predicate.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/last.png)

It will throw an error if the source completes without notification or one that matches the predicate. It returns the last value or if a predicate is provided last value that matches the predicate. It returns the given default value if no notification is emitted or matches the predicate.

## Example 1

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 5),
    rxgo.Last[uint](nil),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 5
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
    rxgo.Empty[string](),
    rxgo.Last[string](nil, "defaultValue"),
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
    rxgo.Last[string](nil),
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
    rxgo.Empty[string](),
    rxgo.Last(func(value string, _ uint) bool {
		return value == "a"
	}),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Error -> rxgo: no values match
```
