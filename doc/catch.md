# Catch Operator

## Overview

Recover from an error by continuing the sequence without error.

## Instances

* `OnErrorResumeNext`: instructs an Observable to pass control to another Observable rather than invoking onError if it encounters an error.
* `OnErrorReturn`: instructs an Observable to emit an item (returned by a specified function) rather than invoking onError if it encounters an error.
* `OnErrorReturnItem`: instructs on Observable to emit an item if it encounters an error.

## Example

### OnErrorResumeNext

```go
observable := rxgo.Just(1, 2, errors.New("foo"))().
	OnErrorResumeNext(func(e error) rxgo.Observable {
		return rxgo.Just(3, 4)()
	})
```

Output:

```
1
2
3
4
```

### OnErrorReturn

```go
observable := rxgo.Just(1, errors.New("2"), 3, errors.New("4"), 5)().
	OnErrorReturn(func(err error) interface{} {
		return err.Error()
	})
```

Output:

```
1
2
3
4
5
```

### OnErrorReturnItem

```go
observable := rxgo.Just(1, errors.New("2"), 3, errors.New("4"), 5)().
	OnErrorReturnItem("foo")
```

Output:

```
1
foo
3
foo
5
```

## Options

* [WithBufferedChannel](options.md#withbufferedchannel)

* [WithContext](options.md#withcontext)

* [WithObservationStrategy](options.md#withobservationstrategy)

* [WithErrorStrategy](options.md#witherrorstrategy)