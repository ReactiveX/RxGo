# BufferTime

> Buffers the source Observable values for a specific time period.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/bufferTime.png)

Buffers values from the source for a specific time duration `bufferTimeSpan`. Unless the optional argument bufferCreationInterval is given, it emits and resets the buffer every `bufferTimeSpan` milliseconds. If bufferCreationInterval is given, this operator opens the buffer every bufferCreationInterval milliseconds and closes (emits and resets) the buffer every `bufferTimeSpan` milliseconds. When the optional argument maxBufferSize is specified, the buffer will be closed either after `bufferTimeSpan` milliseconds or when it contains maxBufferSize elements.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](0, 7),
    rxgo.BufferTime[uint](time.Second),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> [0, 1, 2] // after 1s
// Next -> [1, 2, 3] // after 2s
// Next -> [2, 3, 4] // after 3s
// Next -> [3, 4, 5] // after 1s
// Next -> [4, 5, 6] // after 1s
// Next -> [5, 6]    // after 1s
// Next -> [6]       // after 1s
// Complete!
```
