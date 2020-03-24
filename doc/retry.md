# Retry Operator

## Overview

Implements a retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.

It accepts a `shouldRetry func(error) bool` function to determine whether an error should by retried.

![](http://reactivex.io/documentation/operators/images/retry.png)

## Example

```go
observable := rxgo.Just(1, 2, errors.New("foo"))().
	Retry(2, func(err error) bool {
		return err.Error() == "foo"
	})
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

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)

* [WithPublishStrategy](options.md#withpublishstrategy)