# Timer

## Overview

Creates an observable that will wait for a specified time period before emitting the number 0.

![](http://reactivex.io/documentation/operators/images/timer.png)

## Example 1

Wait 3 seconds and start another observable

You might want to use timer to delay subscription to an observable by a set amount of time.

```go
rxgo.Timer[uint](time.Second * 3).SubscribeSync(func(v uint) {
    log.Println("Timer ->", v)
}, nil, func() {
    log.Println("Complete!")
})

// Output:
// Timer -> 0 # after 3s
// Complete!
```

## Example 2

Start an interval that starts right away
Since interval waits for the passed delay before starting, sometimes that's not ideal. You may want to start an interval immediately. timer works well for this. Here we have both side-by-side so you can see them in comparison.

Note that this observable will never complete.

```go
rxgo.Timer[uint](0, time.Second).SubscribeSync(func(v uint) {
    log.Println("Timer ->", v)
}, nil, func() {
    log.Println("Complete!")
})
// 0 - after 0ms
// 1 - after 1s
// 2 - after 2s
// ...

rxgo.Interval(time.Second).SubscribeSync(func(v uint) {
    log.Println("Interval ->", v)
}, nil, nil)
// 0 - after 1s
// 1 - after 2s
// 2 - after 3s
// ...
```
