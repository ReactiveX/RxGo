# CombineLatestWith

> Create an observable that combines the latest values from all passed observables and the source into arrays and emits them.

## Desceiption

Returns an observable, that when subscribed to, will subscribe to the source observable and all sources provided as arguments. Once all sources emit at least one value, all of the latest values will be emitted as an array. After that, every time any source emits a value, all of the latest values will be emitted as an array.

This is a useful operator for eagerly calculating values based off of changed inputs.

## Example

```go
rxgo.Pipe2(
	rxgo.Interval(time.Second),
	rxgo.CombineLatestWith(
		rxgo.Range[uint](1, 10),
		rxgo.Of2[uint](88),
	),
	rxgo.Take[[]uint](5),
).SubscribeSync(func(v []uint) {
	log.Println("Next ->", v)
}, func(err error) {
	log.Println("Error ->", err)
}, func() {
	log.Println("Complete!")
})

// Output:
// Next -> [0 10 88] // after 1s
// Next -> [1 10 88] // after 2s
// Next -> [2 10 88] // after 3s
// Next -> [3 10 88] // after 4s
// Next -> [4 10 88] // after 5s
// Complete!
```
