# Scan Operator

## Overview

Apply a function to each item emitted by an Observable, sequentially, and emit each successive value.

![](http://reactivex.io/documentation/operators/images/scan.png)

## Example

```go
observable := rxgo.Just(1, 2, 3, 4, 5)().
    Scan(func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
        if acc == nil {
            return elem, nil
        }
        return acc.(int) + elem.(int), nil
    })
```

Output:

```
1
3
6
10
15
```

## Options

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)