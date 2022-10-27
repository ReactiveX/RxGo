# ForkJoin

> Accepts an Array of ObservableInput or a dictionary Object of ObservableInput and returns an Observable that emits either an array of values in the exact same order as the passed array, or a dictionary of values in the same shape as the passed dictionary.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/forkJoin.png)

**ForkJoin** is an operator that takes any number of input observables which can be passed either as an array or a dictionary of input observables. If no input observables are provided (e.g. an empty array is passed), then the resulting stream will complete immediately.

**ForkJoin** will wait for all passed observables to emit and complete and then it will emit an array or an object with last values from corresponding observables.

If you pass an array of n observables to the operator, then the resulting array will have n values, where the first value is the last one emitted by the first observable, second value is the last one emitted by the second observable and so on.

If you pass a dictionary of observables to the operator, then the resulting objects will have the same keys as the dictionary passed, with their last values they have emitted located at the corresponding key.

That means **ForkJoin** will not emit more than once and it will complete after that. If you need to emit combined values not only at the end of the lifecycle of passed observables, but also throughout it, try out combineLatest or zip instead.

In order for the resulting array to have the same length as the number of input observables, whenever any of the given observables completes without emitting any value, **ForkJoin** will complete at that moment as well and it will not emit anything either, even if it already has some last values from other observables. Conversely, if there is an observable that never completes, **ForkJoin** will never complete either, unless at any point some other observable completes without emitting a value, which brings us back to the previous case. Overall, in order for **ForkJoin** to emit a value, all given observables have to emit something at least once and complete.

If any given observable errors at some point, **ForkJoin** will error as well and immediately unsubscribe from the other observables.

Optionally **ForkJoin** accepts a resultSelector function, that will be called with values which normally would land in the emitted array. Whatever is returned by the resultSelector, will appear in the output observable instead. This means that the default resultSelector can be thought of as a function that takes all its arguments and puts them into an array. Note that the resultSelector will be called only when **ForkJoin** is supposed to emit a result.

## Example

```go
rxgo.ForkJoin(
    rxgo.Of2[uint](1, 88, 2, 7215251),
    rxgo.Pipe1(
        rxgo.Interval(time.Millisecond*10),
        rxgo.Take[uint](3),
    ),
).SubscribeSync(func(v []uint) {
    log.Println("Next ->", v)
}, func(err error) {
    log.Println("Error ->", err)
}, func() {
    log.Println("Complete!")
})

// Output:
// Next -> [7215251, 2]
// Complete!
```
