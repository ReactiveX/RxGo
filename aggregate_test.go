package rxgo

import (
	"errors"
	"testing"
)

func TestCount(t *testing.T) {
	t.Run("Count with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Count[any]()), uint(0), nil, true)
	})

	t.Run("Count everything from Range(0,7)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](0, 7), Count[uint]()), uint(7), nil, true)
	})

	t.Run("Count from Range(1,7) with condition", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 7), Count(func(i uint, _ uint) bool {
			return i%2 == 1
		})), uint(4), nil, true)
	})
}

func TestMax(t *testing.T) {

	t.Run("Max with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Max[any]()), nil, nil, true)
	})
}

func TestMin(t *testing.T) {
	type human struct {
		age  int
		name string
	}

	t.Run("Min with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Min[any]()), nil, nil, true)
	})

	t.Run("Min with struct", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Scheduled(
			human{age: 7, name: "Foo"},
			human{age: 5, name: "Bar"},
			human{age: 9, name: "Beer"},
		), Min(func(a, b human) int8 {
			if a.age < b.age {
				return -1
			}
			return 1
		})), human{age: 5, name: "Bar"}, nil, true)
	})
}

func TestReduce(t *testing.T) {
	t.Run("Reduce with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			EMPTY[uint](),
			Reduce(func(acc, cur, _ uint) (uint, error) {
				return acc + cur, nil
			}, 0),
		), uint(0), nil, true)
	})

	t.Run("Reduce with error", func(t *testing.T) {
		var err = errors.New("something happened")
		checkObservableResult(t, Pipe1(
			Range[uint](1, 18),
			Reduce(func(acc, cur, idx uint) (uint, error) {
				if idx > 5 {
					return 0, err
				}
				return acc + cur, nil
			}, 0),
		), uint(0), err, false)
	})

	t.Run("Reduce with zero default value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Scheduled[uint](1, 3, 5),
			Reduce(func(acc, cur, _ uint) (uint, error) {
				return acc + cur, nil
			}, 0),
		), uint(9), nil, true)
	})
}
