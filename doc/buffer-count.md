# BufferCount

> Buffers the source Observable values until the size hits the maximum bufferSize given.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/bufferCount.png)

Buffers a number of values from the source Observable by `bufferSize` then emits the buffer and clears it, and starts a new buffer each `startBufferEvery` values. If `startBufferEvery` is not provided, then new buffers are started immediately at the start of the source and when each buffer closes and is emitted.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](0, 7),
    rxgo.BufferCount[uint](3, 1),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> [0, 1, 2]
// Next -> [1, 2, 3]
// Next -> [2, 3, 4]
// Next -> [3, 4, 5]
// Next -> [4, 5, 6]
// Next -> [5, 6]
// Next -> [6]
// Complete!
```
