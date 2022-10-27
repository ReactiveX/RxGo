# RaceWith

> Creates an Observable that mirrors the first source Observable to emit a next, error or complete notification from the combination of the Observable to which the operator is applied and supplied Observables.

## Example

```go
rxgo.Pipe2(
    rxgo.Pipe1(
        rxgo.Interval(time.Millisecond*7),
        rxgo.Map(func(v uint, _ uint) (string, error) {
            return fmt.Sprintf("slowest(%v)", v), nil
        }),
    ),
    rxgo.RaceWith(
        rxgo.Pipe1(
            rxgo.Interval(time.Millisecond*3),
            rxgo.Map(func(v uint, _ uint) (string, error) {
                return fmt.Sprintf("fastest(%v)", v), nil
            }),
        ),
        rxgo.Pipe1(
            rxgo.Interval(time.Millisecond*5),
            rxgo.Map(func(v uint, _ uint) (string, error) {
                return fmt.Sprintf("average(%v)", v), nil
            }),
        ),
    ),
    Take[string](5),
).SubscribeSync(func(v string) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> fastest(0)
// Next -> fastest(1)
// Next -> fastest(2)
// Next -> fastest(3)
// Next -> fastest(4)
// Complete!
```
