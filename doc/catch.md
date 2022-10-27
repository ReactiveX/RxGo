# Catch

> Catches errors on the observable to be handled by returning a new observable or throwing an error.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/catch.png)

This operator handles errors, but forwards along all other events to the resulting observable. If the source observable terminates with an error, it will map that error to a new observable, subscribe to it, and forward all of its events to the resulting observable.

## Example 1

```go
rxgo.Pipe2(
    rxgo.Range[uint](1, 5),
    rxgo.Map(func(v, _ uint) (uint, error) {
        if v == 4 {
            return 0, errors.New("four")
        }
        return v, nil
    }),
    rxgo.Catch(func(err error, caught rxgo.Observable[uint]) rxgo.Observable[uint] {
        return rxgo.Range[uint](1, 5)
    }),
).SubscribeSync(func(v uint) {
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
// Next -> 4
// Next -> 5
// Complete!
```

## Example 2

```go
rxgo.Pipe3(
    rxgo.Range[uint](1, 5),
    rxgo.Map(func(v, _ uint) (uint, error) {
        if v == 4 {
            return 0, errors.New("four")
        }
        return v, nil
    }),
    rxgo.Catch(func(err error, caught rxgo.Observable[uint]) rxgo.Observable[uint] {
        return rxgo.Range[uint](1, 5)
    }),
    rxgo.Take[uint](5),
).SubscribeSync(func(v uint) {
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
// Complete!
```
