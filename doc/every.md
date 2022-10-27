# Every

> Returns an Observable that emits whether or not every item of the source satisfies the condition specified.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/every.png)

If all values pass predicate before the source completes, emits true before completion, otherwise emit false, then complete.

## Example

```go
rxgo.Pipe1(
    rxgo.Range[uint](1, 7),
    rxgo.Every(func(value, index uint) bool {
        return value < 10
    }),
).SubscribeSync(func(v bool) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> true
// Complete!
```
