# Do Operator

## Overview

Register an action to take upon a variety of Observable lifecycle events.

![](http://reactivex.io/documentation/operators/images/do.c.png)

## Instances

* `DoOnNext`
* `DoOnError`
* `DoOnCompleted`

Each one returns a `<-chan struct{}` that closes once the Observable terminates.

## Example

### DoOnNext

```go
<-rxgo.Just(1, 2, 3)().
	DoOnNext(func(i interface{}) {
		fmt.Println(i)
	})
```

Output:

```
1
2
3
```

### DoOnError

```go
<-rxgo.Just(1, 2, errors.New("foo"))().
	DoOnError(func(err error) {
		fmt.Println(err)
	})
```

Output:

```
foo
```

### DoOnCompleted

```go
<-rxgo.Just(1, 2, 3)().
	DoOnCompleted(func() {
		fmt.Println("done")
	})
```

Output:

```
done
```

## Options

* [WithContext](options.md#withcontext)