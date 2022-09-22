# Map

> Applies a given project function to each value emitted by the source Observable, and emits the resulting values as an Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/map.png)

This operator applies a projection to each value and emits that projection in the output Observable.

## Example 1

```go
rxgo.Pipe1(
	rxgo.Range[uint8](1, 5),
	rxgo.Map(func(v uint8, index uint) (string, error) {
		return fmt.Sprintf("%d", v), nil
	}),
).SubscribeSync(func(v string) {
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
// Next -> 4
// Next -> 5
// Complete!
```

## Example 2

```go
rxgo.Pipe1(
	rxgo.Range[uint8](1, 10),
	rxgo.Map(func(v uint8, index uint) (string, error) {
		if v > 5 {
			return "", errors.New("the value is greater than 5")
		}
		return fmt.Sprintf("%d", v), nil
	}),
).SubscribeSync(func(v string) {
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
// Next -> 4
// Next -> 5
// Error -> the value is greater than 5
```
