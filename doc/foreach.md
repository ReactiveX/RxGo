# ForEach Operator

## Overview

Subscribe to an Observable and register `OnNext`, `OnError` and `OnCompleted` actions.

It returns a `<-chan struct{}` that closes once the Observable terminates.

## Example

```go
<-rxgo.Just(1, errors.New("foo"))().
	ForEach(
		func(i interface{}) {
			fmt.Printf("next: %v\n", i)
		}, func(err error) {
			fmt.Printf("error: %v\n", err)
		}, func() {
			fmt.Println("done")
		})
```

Output:

```
next: 1
error: foo
done
```

## Options

* [WithContext](options.md#withcontext)

# FromChannel Operator

## Overview

Create a cold observable from a channel.

## Example

```go
ch := make(chan rxgo.Item)
observable := rxgo.FromChannel(ch)
```

The items are buffered in the channel until an Observer subscribes.