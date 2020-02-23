# Retry Operator

## Overview

Implements a retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.

![](http://reactivex.io/documentation/operators/images/retry.png)

## Example

```go
observable := rxgo.Just([]interface{}{1, 2, errors.New("foo")}).Retry(2)
```

Output:

```
1
2
1
2
1
2
foo
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