package rxgo

import (
	"errors"
	"testing"
)

func TestDefaultIfEmpty(t *testing.T) {
	t.Run("DefaultIfEmpty with any", func(t *testing.T) {
		str := "hello world"
		checkObservableResult(t, Pipe1(EMPTY[any](), DefaultIfEmpty[any](str)), any(str), nil, true)
	})

	t.Run("DefaultIfEmpty with non-empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 3), DefaultIfEmpty[uint](100)), []uint{1, 2, 3}, nil, true)
	})
}

func TestEvery(t *testing.T) {
	t.Run("Every with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[uint](), Every(func(value, index uint) bool {
			return value < 10
		})), true, nil, true)
	})

	t.Run("Every with all value match the condition", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 7), Every(func(value, index uint) bool {
			return value < 10
		})), true, nil, true)
	})

	t.Run("Every with not all value match the condition", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 7), Every(func(value, index uint) bool {
			return value < 5
		})), false, nil, true)
	})
}

func TestFind(t *testing.T) {
	t.Run("Find with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Find(func(a any, u uint) bool {
			return a == nil
		})), None[any](), nil, true)
	})

	t.Run("Find with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Scheduled("a", "b", "c", "d", "e"),
			Find(func(v string, u uint) bool {
				return v == "c"
			}),
		), Some("c"), nil, true)
	})
}

func TestFindIndex(t *testing.T) {
	t.Run("FindIndex with value that doesn't exist", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), FindIndex(func(a any, u uint) bool {
			return a == nil
		})), -1, nil, true)
	})

	t.Run("FindIndex with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Scheduled("a", "b", "c", "d", "e"),
			FindIndex(func(v string, u uint) bool {
				return v == "c"
			}),
		), 2, nil, true)
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), IsEmpty[any]()), true, nil, true)
	})

	t.Run("IsEmpty with error", func(t *testing.T) {
		var err = errors.New("something wrong")
		checkObservableResult(t, Pipe1(Scheduled[any](err), IsEmpty[any]()), false, err, false)
	})

	t.Run("IsEmpty with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 3), IsEmpty[uint]()), false, nil, true)
	})
}

func TestThrowIfEmpty(t *testing.T) {
	t.Run("ThrowIfEmpty with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), ThrowIfEmpty[any]()), nil, ErrEmpty, false)
	})

	t.Run("ThrowIfEmpty with error factory", func(t *testing.T) {
		var err = errors.New("something wrong")
		checkObservableResult(t, Pipe1(EMPTY[any](), ThrowIfEmpty[any](func() error {
			return err
		})), nil, err, false)
	})

	t.Run("ThrowIfEmpty with value", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 3), ThrowIfEmpty[uint]()), []uint{1, 2, 3}, nil, true)
	})
}
