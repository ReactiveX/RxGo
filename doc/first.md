# First

> Emits only the first value (or the first value that meets some condition) emitted by the source Observable.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/first.png)

If called with no arguments, `First` emits the first value of the source Observable, then completes. If called with a predicate function, `First` emits the first value of the source that matches the specified condition. Throws an error if `defaultValue` was not provided and a matching element is not found.

## Example

```go
observable := rxgo.Just(1, 2, 3)().First()
```

Output:

```
true
```

## Options

- [WithBufferedChannel](options.md#withbufferedchannel)

- [WithContext](options.md#withcontext)

- [WithObservationStrategy](options.md#withobservationstrategy)

- [WithPublishStrategy](options.md#withpublishstrategy)
