# Skip

> Returns an Observable that skips the first count items emitted by the source Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/skip.png)

Skips the values until the sent notifications are equal or less than provided skip count. It raises an error if skip count is equal or more than the actual number of emits and source raises an error.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](0, 10),
	rxgo.Skip(3),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 3
// Next -> 4
// Next -> 5
// Next -> 6
// Next -> 7
// Next -> 8
// Next -> 9
// Next -> 10
// Complete!
```
