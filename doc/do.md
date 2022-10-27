# Do

> Used to perform side-effects for notifications from the source observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/tap.png)

**Do** is designed to allow the developer a designated place to perform side effects. While you could perform side-effects inside of a map or a mergeMap, that would make their mapping functions impure, which isn't always a big deal, but will make it so you can't do things like memoize those functions. The **Do** operator is designed solely for such side-effects to help you remove side-effects from other operations.

For any notification, next, error, or complete, **Do** will call the appropriate callback you have provided to it, via a function reference, or a partial observer, then pass that notification down the stream.

The observable returned by **Do** is an exact mirror of the source, with one exception: Any error that occurs -- synchronously -- in a handler provided to **Do** will be emitted as an error from the returned observable.

**Be careful! You can mutate objects as they pass through the Do operator's handlers.**

## Example

```go
rxgo.Pipe1(
	rxgo.Range[uint](1, 5),
	rxgo.Do(NewObserver(func(v uint) {
		log.Println("DoNext ->", v)
	}, func(err error) {
		log.Println("DoError ->", err)
	}, func() {
		log.Println("DoComplete!")
	})),
).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Debug -> 1
// Next -> 1
// Debug -> 2
// Next -> 2
// Debug -> 3
// Next -> 3
// Debug -> 4
// Next -> 4
// Debug -> 5
// Next -> 5
// Complete!
```
