# Defer Operator

> Creates an Observable that, on subscribe, calls an Observable factory to make an Observable for each new Observer.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/defer.png)

`Defer` allows you to create an Observable only when the Observer subscribes. It waits until an Observer subscribes to it, calls the given factory function to get an Observable -- where a factory function typically generates a new Observable -- and subscribes the Observer to this Observable. In case the factory function returns a falsy value, then `Empty` is used as Observable instead. Last but not least, an exception during the factory function call is transferred to the Observer by calling error.

## Example

```go
rxgo.Defer(func() rxgo.Observable[uint] {
	if rand.Intn(10) > 5 {
		return rxgo.Interval(time.Second)
	}
	return rxgo.Throw[uint](func() error {
		return errors.New("failed")
	})
}).SubscribeSync(func(v uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// If the result of `rand.Intn(10)` is greater than 5 it will emit ascending numbers, one every second(1000ms);
// else it will return an `Empty` observable and complete.
//
// Output 1:
// Complete!
//
// Output 2:
// Next -> 0 # after 1s
// Next -> 1 # after 1s
// Next -> 2 # after 1s
// ...
```
