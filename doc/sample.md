# Sample Operator

## Overview

Emit the most recent item emitted by an Observable within periodic time intervals.

![](http://reactivex.io/documentation/operators/images/sample.png)

## Example

```go
sampledObservable := observable1.Sample(observable2)
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithPublishStrategy](options.md#withpublishstrategy)