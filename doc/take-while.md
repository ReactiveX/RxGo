# TakeWhile

> Emits values emitted by the source Observable so long as each value satisfies the given predicate, and then completes as soon as this predicate is not satisfied.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/takeWhile.png)

**TakeWhile** subscribes and begins mirroring the source Observable. Each value emitted on the source is given to the `predicate` function which returns a boolean, representing a condition to be satisfied by the source values. The output Observable emits the source values until such time as the predicate returns false, at which point takeWhile stops mirroring the source Observable and completes the output Observable.

## Example 1

```go
rxgo.Pipe1(
	rxgo.Range[uint](1, 100),
	rxgo.TakeWhile(func(v uint, _ uint) bool {
		return v >= 50
	}),
).SubscribeSync(func(v uint) {
	log.Println("Next ->", v)
}, func(err error) {
	log.Println("Error ->", err)
}, func() {
	log.Println("Complete!")
})

// Output :
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
	rxgo.Range[uint](1, 100),
	rxgo.TakeWhile(func(v uint, _ uint) bool {
		return v <= 5
	}),
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
// Next -> 4
// Next -> 5
// Complete!
```
