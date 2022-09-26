# Distinct

> Returns an Observable that emits all items emitted by the source Observable that are distinct by comparison from previous items.

## Description

If a keySelector function is provided, then it will project each value from the source observable into a new value that it will check for equality with previously projected values. If the keySelector function is not provided, it will use each value from the source observable directly with an equality check against previous values.

In JavaScript runtimes that support Set, this operator will use a Set to improve performance of the distinct value checking.

In other runtimes, this operator will use a minimal implementation of Set that relies on an Array and indexOf under the hood, so performance will degrade as more values are checked for distinction. Even in newer browsers, a long-running distinct use might result in memory leaks. To help alleviate this in some scenarios, an optional flushes parameter is also provided so that the internal Set can be "flushed", basically clearing it of values.

## Example 1

```go
rxgo.Pipe1(
	rxgo.Of2[uint](1, 1, 2, 2, 2, 1, 2, 3, 4, 3, 2, 1),
	rxgo.Distinct[uint](),
).SubscribeSync(func(v []uint) {
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
// Complete!
```

## Example 2

```go
type user struct {
	name string
	age  uint
}

rxgo.Pipe1(
	rxgo.Of2(
		user{name: "Foo", age: 4},
		user{name: "Bar", age: 7},
		user{name: "Foo", age: 5},
	),
	rxgo.Distinct(func(v user) string {
		return v.name
	}),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> {age: 4, name: "Foo"}
// Next -> {age: 7, name: "Bar"}
// Complete!
```
