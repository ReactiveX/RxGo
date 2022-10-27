# ZipWith

> Subscribes to the source, and the observable inputs provided as arguments, and combines their values, by index, into arrays.

## Description

What is meant by "combine by index": The first value from each will be made into a single array, then emitted, then the second value from each will be combined into a single array and emitted, then the third value from each will be combined into a single array and emitted, and so on.

This will continue until it is no longer able to combine values of the same index into an array.

After the last value from any one completed source is emitted in an array, the resulting observable will complete, as there is no way to continue "zipping" values together by index.

Use-cases for this operator are limited. There are memory concerns if one of the streams is emitting values at a much faster rate than the others. Usage should likely be limited to streams that emit at a similar pace, or finite streams of known length.

## Example

```go
rxgo.Pipe1(
    rxgo.Of2[any](27, 25, 29),
    rxgo.ZipWith(
        rxgo.Of2[any]("Foo", "Bar", "Beer"),
        rxgo.Of2[any](true, true, false),
    ),
).SubscribeSync(func(v []any) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> [27 Foo true]
// Next -> [25 Bar true]
// Next -> [29 Beer false]
// Complete!
```
