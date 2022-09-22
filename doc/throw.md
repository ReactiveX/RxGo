# Thrown Operator

## Overview

Creates an observable that will create an error instance and push it to the consumer as an error immediately upon subscription.

![](http://reactivex.io/documentation/operators/images/throw.c.png)

## Example

```go
rxgo.Throw[string](func() error {
    return errors.New("foo")
}).SubscribeSync(nil, func(err error) {
    log.Println("Error ->", err)
}, nil)

// Output:
// Error -> foo
```

This creation function is useful for creating an observable that will create an error and error every time it is subscribed to. Generally, inside of most operators when you might want to return an errored observable, this is unnecessary. In most cases, such as in the inner return of ConcatMap, MergeMap, Defer, and many others, you can simply throw the error, and RxGo will pick that up and notify the consumer of the error.
