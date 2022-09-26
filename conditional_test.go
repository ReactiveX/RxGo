package rxgo

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDefaultIfEmpty(t *testing.T) {
	t.Run("DefaultIfEmpty with any", func(t *testing.T) {
		str := "hello world"
		checkObservableResult(t, Pipe1(
			Empty[any](),
			DefaultIfEmpty[any](str),
		), any(str), nil, true)
	})

	t.Run("DefaultIfEmpty with error", func(t *testing.T) {
		var err = fmt.Errorf("some error")
		checkObservableResult(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			DefaultIfEmpty[any]("hello world"),
		), nil, err, false)
	})

	t.Run("DefaultIfEmpty with non-empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint](1, 3),
			DefaultIfEmpty[uint](100),
		), []uint{1, 2, 3}, nil, true)
	})
}

func TestEvery(t *testing.T) {
	t.Run("Every with Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[uint](),
			Every(func(value, index uint) bool {
				return value < 10
			}),
		), true, nil, true)
	})

	t.Run("Every with error", func(t *testing.T) {
		var err = fmt.Errorf("some error")
		checkObservableResult(t, Pipe1(
			Throw[uint](func() error {
				return err
			}),
			Every(func(value, index uint) bool {
				return value < 10
			}),
		), false, err, false)
	})

	t.Run("Every with all value match the condition", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Range[uint](1, 7),
			Every(func(value, index uint) bool {
				return value < 10
			}),
		), true, nil, true)
	})

	t.Run("Every with not all value match the condition", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Range[uint](1, 7),
			Every(func(value, index uint) bool {
				return value < 5
			}),
		), false, nil, true)
	})
}

func TestFind(t *testing.T) {
	t.Run("Find with Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[any](),
			Find(func(a any, u uint) bool {
				return a == nil
			}),
		), None[any](), nil, true)
	})

	t.Run("Find with error", func(t *testing.T) {
		var err = errors.New("some error")
		checkObservableResult(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			Find(func(a any, u uint) bool {
				return a == "not found"
			}),
		), nil, err, false)
	})

	t.Run("Find with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Of2("a", "b", "c", "d", "e"),
			Find(func(v string, u uint) bool {
				return v == "c"
			}),
		), Some("c"), nil, true)
	})

	t.Run("Find with Range", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Range[uint8](1, 10),
			Find(func(v uint8, u uint) bool {
				return v == 5
			}),
		), Some[uint8](5), nil, true)
	})
}

func TestFindIndex(t *testing.T) {
	t.Run("FindIndex with value that doesn't exist", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[any](),
			FindIndex(func(a any, u uint) bool {
				return a == nil
			}),
		), -1, nil, true)
	})

	t.Run("FindIndex with error", func(t *testing.T) {
		var err = errors.New("some error")
		checkObservableResult(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			FindIndex(func(a any, u uint) bool {
				return a == nil
			}),
		), int(0), err, false)
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
	t.Run("IsEmpty with Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[any](),
			IsEmpty[any](),
		), true, nil, true)
	})

	t.Run("IsEmpty with error", func(t *testing.T) {
		var err = errors.New("something wrong")
		checkObservableResult(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			IsEmpty[any](),
		), false, err, false)
	})

	t.Run("IsEmpty with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Range[uint](1, 3),
			IsEmpty[uint](),
		), false, nil, true)
	})
}

func TestSequenceEqual(t *testing.T) {
	t.Run("SequenceEqual with one Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[uint](),
			SequenceEqual(Range[uint](1, 10)),
		), false, nil, true)
	})

	t.Run("SequenceEqual with all Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[any](),
			SequenceEqual(Empty[any]()),
		), true, nil, true)
	})

	t.Run("SequenceEqual with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResult(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			SequenceEqual(Of2[any](1, 100, 88)),
		), false, err, false)
	})

	t.Run("SequenceEqual with values (tally)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Range[uint](1, 10),
			SequenceEqual(Range[uint](1, 10)),
		), true, nil, true)
	})

	// eventually it will be tally
	t.Run("SequenceEqual with values (tally but different pacing)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Pipe1(Interval(time.Millisecond), Take[uint](3)),
			SequenceEqual(Pipe1(Interval(time.Millisecond*5), Take[uint](3))),
		), true, nil, true)
	})

	t.Run("SequenceEqual with values (not tally)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Of2("a", "i", "o", "z"),
			SequenceEqual(Of2("a", "i", "v", "z")),
		), false, nil, true)
	})
}

func TestThrowIfEmpty(t *testing.T) {
	t.Run("ThrowIfEmpty with Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[any](),
			ThrowIfEmpty[any](),
		), nil, ErrEmpty, false)
	})

	t.Run("ThrowIfEmpty with error factory", func(t *testing.T) {
		var err = errors.New("something wrong")
		checkObservableResult(t, Pipe1(
			Empty[any](),
			ThrowIfEmpty[any](func() error {
				return err
			}),
		), nil, err, false)
	})

	t.Run("ThrowIfEmpty with value", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint](1, 3),
			ThrowIfEmpty[uint](),
		), []uint{1, 2, 3}, nil, true)
	})
}
