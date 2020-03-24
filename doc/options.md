# Operator Options

Most of the operators accept a list of functional options to impact the Observable behaviour.

As an example:

```go
observable.Map(func(_ context.Context, i interface{}) (interface{}, error) {
	return i.(int) * 10, nil
}, rxgo.WithContext(),
   rxgo.WithCPUPool(),
   rxgo.WihtBufferedChannel(1))
```

## WithBufferedChannel

Configure the capacity of the output channel.

```go
rxgo.WithBufferedChannel(1) // Create a buffered channel with a 1 capacity
```

## WithContext

Allows passing a context. The Observable will listen to its done signal to close itself.

```go
rxgo.WithContext(ctx)
```

## WithObservationStrategy

* Lazy (default): consume when an Observer starts to subscribe.

```go
rxgo.WithObservation(rxgo.Lazy)
```

* Eager: consumer when the Observable is created:

```go
rxgo.WithObservation(rxgo.Eager)
```

## WithErrorStrategy

* StopOnError (default): stop processing if the Observable produces an error.

```go
rxgo.WithErrorStrategy(rxgo.StopOnError)
```

* ContinueOnError: continue processing items if the Observable produces an error.

```go
rxgo.WithErrorStrategy(rxgo.ContinueOnError)
```

This strategy is propagated to the parent(s) Observable(s).

## WithPool

Convert the operator in a parallel operator and specify the number of concurrent goroutines.

```go
rxgo.WithPool(8) // Creates a pool of 8 goroutines
```

## WithCPUPool

Convert the operator in a parallel operator and specify the number of concurrent goroutines as `runtime.NumCPU()`.

```go
rxgo.WithCPUPool()
```

## Serialize

Force an Observable to produce items sequentially.

```go
rxgo.Serialize()
```

This option should be used in coordination with `rxgo.WithPool(n)` or `rxgo.WithCPUPool()`.

## WithPublishStrategy

Create a [Connectable Observable](../README.md#connectable-observable).

```go
rxgo.WithPublishStrategy()
```

This option is propagated to the parent(s) Observable(s).