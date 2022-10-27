# Reduce

> Applies an accumulator function over the source Observable, and returns the accumulated result when the source completes, given an optional seed value.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/reduce.png)

`Reduce` applies an accumulator function against an accumulation and each value of the source Observable (from the past) to reduce it to a single value, emitted on the output Observable. Note that reduce will only emit one value, only when the source Observable completes. It is equivalent to applying operator `Scan` followed by operator `Last`.

Returns an Observable that applies a specified accumulator function to each item emitted by the source Observable. If a seed value is specified, then that value will be used as the initial value for the accumulator. If no seed value is specified, the first item of the source is used as the seed.

## Example

```go
rxgo.Pipe1(
	rxgo.Range[uint](1, 18),
	rxgo.Reduce(func(acc, cur, idx uint) (uint, error) {
		return acc + cur, nil
	}, 0),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 171
// Complete!
```
