# DistinctUntilChanged

> Returns a result Observable that emits all values pushed by the source observable if they are distinct in comparison to the last value the result observable emitted.

## Description

When provided without parameters or with the first parameter (comparator), it behaves like this:

It will always emit the first value from the source.
For all subsequent values pushed by the source, they will be compared to the previously emitted values using the provided comparator or an === equality check.
If the value pushed by the source is determined to be unequal by this check, that value is emitted and becomes the new "previously emitted value" internally.
When the second parameter (keySelector) is provided, the behavior changes:

It will always emit the first value from the source.
The keySelector will be run against all values, including the first value.
For all values after the first, the selected key will be compared against the key selected from the previously emitted value using the comparator.
If the keys are determined to be unequal by this check, the value (not the key), is emitted and the selected key from that value is saved for future comparisons against other keys.

## Example 1

```go
rxgo.Pipe1(
	rxgo.Of2("a", "a", "b", "a", "c", "c", "d"),
	rxgo.DistinctUntilChanged[string](),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> a
// Next -> b
// Next -> c
// Next -> d
// Complete!
```

## Example 2

```go
type build struct {
	engineVersion       string
	transmissionVersion string
}

rxgo.Pipe1(
	rxgo.Of2(
		build{engineVersion: "1.1.0", transmissionVersion: "1.2.0"},
		build{engineVersion: "1.1.0", transmissionVersion: "1.4.0"},
		build{engineVersion: "1.3.0", transmissionVersion: "1.4.0"},
		build{engineVersion: "1.3.0", transmissionVersion: "1.5.0"},
		build{engineVersion: "2.0.0", transmissionVersion: "1.5.0"},
	),
	rxgo.DistinctUntilChanged(func(prev, curr build) bool {
		return (prev.engineVersion == curr.engineVersion ||
			prev.transmissionVersion == curr.transmissionVersion)
	}),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> {engineVersion: "1.1.0", transmissionVersion: "1.2.0"}
// Next -> {engineVersion: "1.3.0", transmissionVersion: "1.4.0"}
// Next -> {engineVersion: "2.0.0", transmissionVersion: "1.5.0"}
// Complete!
```
