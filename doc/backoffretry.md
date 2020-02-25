# BackOffRetry Operator

## Overview

Implements a backoff retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.

The backoff configuration relies on [github.com/cenkalti/backoff/v4](github.com/cenkalti/backoff/v4).

![](http://reactivex.io/documentation/operators/images/retry.png)

## Example

```go
// Backoff retry configuration
backOffCfg := backoff.NewExponentialBackOff()
backOffCfg.InitialInterval = 10 * time.Millisecond

observable := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item, done func()) {
	next <- rxgo.Of(1)
	next <- rxgo.Of(2)
	next <- rxgo.Error(errors.New("foo"))
	done()
}}).BackOffRetry(backoff.WithMaxRetries(backOffCfg, 2))
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