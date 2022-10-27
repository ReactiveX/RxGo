# SkipWhile

> Returns an Observable that skips all items emitted by the source Observable as long as a specified condition holds true, but emits all further source items as soon as the condition becomes false.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/skipWhile.png)

Skips all the notifications with a truthy predicate. It will not skip the notifications when the predicate is falsy. It can also be skipped using index. Once the predicate is true, it will not be called again.

## Example

```go
rxgo.Pipe1(
	rxgo.Of2("Green Arrow", "SuperMan", "Flash", "SuperGirl", "Black Canary"),
	rxgo.SkipWhile(func(v string, _ uint) bool {
		return v != "SuperGirl"
	}),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output :
// Next -> SuperGirl
// Next -> Black Canary
// Complete!
```
