# Assert API

## Overview

There is a public API to facilitate writing unit tests while using RxGo. This is based on `rxgo.Assert`. It automatically creates an Observer on an Iterable (Observable, Single or OptionalSingle).

```go
func TestMap(t *testing.T) {
	err := errors.New("foo")
	observable := rxgo.Just(1, 2, 3)().
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			if i == 3 {
				return nil, err
			}
			return i, nil
		})

	rxgo.Assert(context.Background(), t, observable,
		rxgo.HasItems(1, 2),
		rxgo.HasError(err))
}
```

```
=== RUN   TestMap
--- PASS: TestMap (0.00s)
PASS
```

## Predicates

The API accepts the following list of predicates:

### HasItems

Check whether an Observable produced a given list of items.

```go
rxgo.Assert(ctx, t, observable, rxgo.HasItems(1, 2, 3))
```

### HasItem

Check whether an Single or OptionalSingle produced a given item.

```go
rxgo.Assert(ctx, t, single, rxgo.HasItem(1))
```

### HasItemsNoOrder

Check whether an Observable produced a given item regardless of the ordering (useful when we use a pool option)

```go
rxgo.Assert(ctx, t, single, rxgo.HasItemsNoOrder(1, 2, 3))
```

### IsEmpty

Check whether an Iterable produced no item(s).

```go
rxgo.Assert(ctx, t, single, rxgo.IsEmpty())
```

### IsNotEmpty

Check whether an Iterable produced any item(s).

```go
rxgo.Assert(ctx, t, single, rxgo.IsNotEmpty())
```

### HasError

Check whether an Iterable produced a given error.

```go
rxgo.Assert(ctx, t, single, rxgo.HasError(expectedErr))
```

### HasAnError

Check whether an Iterable produced an error.

```go
rxgo.Assert(ctx, t, single, rxgo.HasAnError())
```

### HasErrors

Check whether an Iterable produced a given list of errors (useful when we use `rxgo.ContinueOnError` strategy).

```go
rxgo.Assert(ctx, t, single, rxgo.HasErrors(expectedErr1, expectedErr2, expectedErr3))
```

### HasNoError

Check whether an Iterable did not produced any error.

```go
rxgo.Assert(ctx, t, single, rxgo.HasNoError())
```

### CustomPredicate

Implement a custom predicate.

```go
rxgo.Assert(ctx, t, observable, rxgo.CustomPredicate(func(items []interface{}) error {
	if len(items) != 3 {
		return errors.New("wrong number of items")
	}
	return nil
}))
```