# Take

> Emits only the first count values emitted by the source Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/take.png)

`Take` returns an Observable that emits only the first `count` values emitted by the source Observable. If the source emits fewer than `count` values then all of its values are emitted. After that, it completes, regardless if the source completes.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](0, 10),
	rxgo.Take(3),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> 0
// Next -> 1
// Next -> 2
// Complete!
```
