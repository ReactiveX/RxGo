# Retry

> Returns an Observable that mirrors the source Observable with the exception of an error.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/retry.png)

If the source Observable calls error, this method will resubscribe to the source Observable for a maximum of count resubscriptions rather than propagating the error call.

Any and all items emitted by the source Observable will be emitted by the resulting Observable, even those emitted during failed subscriptions. For example, if an Observable fails at first but emits [1, 2] then succeeds the second time and emits: [1, 2, 3, 4, 5, complete] then the complete stream of emissions and notifications would be: [1, 2, 1, 2, 3, 4, 5, complete].

## Example

```go
var err = errors.New("failed")

rxgo.Pipe2(
	rxgo.Range[uint8](1, 5),
	rxgo.Map(func(v uint8, _ uint) (uint8, error) {
		if v == 4 {
			return 0, err
		}
		return v, nil
	}),
	rxgo.Retry[uint8, uint](2),
).SubscribeSync(func(v uint8) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> 1
// Next -> 2
// Next -> 3
// Next -> 1
// Next -> 2
// Next -> 3
// Next -> 1
// Next -> 2
// Next -> 3
// Error -> failed
```
