package rxgo

import (
	"errors"
	"testing"
	"time"
)

func TestDebounceTime(t *testing.T) {
	t.Run("DebounceTime with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			EMPTY[any](),
			DebounceTime[any](time.Millisecond),
		), nil, nil, true)
	})

	t.Run("DebounceTime with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResult(t, Pipe1(
			ThrownError[any](func() error {
				return err
			}), DebounceTime[any](time.Millisecond),
		), nil, err, false)
	})
}

func TestDistinct(t *testing.T) {
	t.Run("Distinct with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Distinct(func(value any) int {
			return value.(int)
		})), nil, nil, true)
	})

	t.Run("Distinct with numbers", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Scheduled(1, 1, 2, 2, 2, 1, 2, 3, 4, 3, 2, 1), Distinct(func(value int) int {
			return value
		})), []int{1, 2, 3, 4}, nil, true)
	})

	t.Run("Distinct with struct", func(t *testing.T) {
		type user struct {
			name string
			age  uint
		}

		checkObservableResults(t, Pipe1(Scheduled(
			user{name: "Foo", age: 4},
			user{name: "Bar", age: 7},
			user{name: "Foo", age: 5},
		), Distinct(func(v user) string {
			return v.name
		})), []user{
			{age: 4, name: "Foo"},
			{age: 7, name: "Bar"},
		}, nil, true)
	})
}

func TestDistinctUntilChanged(t *testing.T) {
	t.Run("DistinctUntilChanged with empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), DistinctUntilChanged[any]()), nil, nil, true)
	})

	t.Run("DistinctUntilChanged with string", func(t *testing.T) {
		checkObservableResults(t,
			Pipe1(Scheduled("a", "a", "b", "a", "c", "c", "d"), DistinctUntilChanged[string]()),
			[]string{"a", "b", "a", "c", "d"}, nil, true)
	})

	t.Run("DistinctUntilChanged with numbers", func(t *testing.T) {
		checkObservableResults(t,
			Pipe1(
				Scheduled(30, 31, 20, 34, 33, 29, 35, 20),
				DistinctUntilChanged(func(prev, current int) bool {
					return current <= prev
				}),
			), []int{30, 31, 34, 35}, nil, true)
	})

	t.Run("DistinctUntilChanged with struct", func(t *testing.T) {
		type build struct {
			engineVersion       string
			transmissionVersion string
		}
		checkObservableResults(t,
			Pipe1(
				Scheduled(
					build{engineVersion: "1.1.0", transmissionVersion: "1.2.0"},
					build{engineVersion: "1.1.0", transmissionVersion: "1.4.0"},
					build{engineVersion: "1.3.0", transmissionVersion: "1.4.0"},
					build{engineVersion: "1.3.0", transmissionVersion: "1.5.0"},
					build{engineVersion: "2.0.0", transmissionVersion: "1.5.0"},
				),
				DistinctUntilChanged(func(prev, curr build) bool {
					return (prev.engineVersion == curr.engineVersion ||
						prev.transmissionVersion == curr.transmissionVersion)
				}),
			),
			[]build{
				{engineVersion: "1.1.0", transmissionVersion: "1.2.0"},
				{engineVersion: "1.3.0", transmissionVersion: "1.4.0"},
				{engineVersion: "2.0.0", transmissionVersion: "1.5.0"},
			}, nil, true)
	})

	t.Run("DistinctUntilChanged with Struct(complex)", func(t *testing.T) {
		type account struct {
			updatedBy string
			data      []string
		}
		checkObservableResults(t,
			Pipe1(
				Scheduled(
					account{updatedBy: "blesh", data: []string{}},
					account{updatedBy: "blesh", data: []string{}},
					account{updatedBy: "jamieson"},
					account{updatedBy: "jamieson"},
					account{updatedBy: "blesh"},
				),
				DistinctUntilChanged[account](),
			),
			[]account{
				{updatedBy: "blesh", data: []string{}},
				{updatedBy: "jamieson"},
				{updatedBy: "blesh"},
			}, nil, true)
	})
}

func TestElementAt(t *testing.T) {
	t.Run("ElementAt with default value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), ElementAt[any](1, 10)), 10, nil, true)
	})

	t.Run("ElementAt position 2", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 100), ElementAt[uint](2)), 3, nil, true)
	})

	t.Run("ElementAt with error (ErrArgumentOutOfRange)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 10), ElementAt[uint](100)), 0, ErrArgumentOutOfRange, false)
	})
}

func TestFilter(t *testing.T) {
	t.Run("Filter with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Filter[any](nil)), nil, nil, true)
	})

	t.Run("Filter with error", func(t *testing.T) {
		var err = errors.New("throw")
		checkObservableResult(t, Pipe1(
			ThrownError[any](func() error {
				return err
			}), Filter[any](nil)), nil, err, false)
	})

	t.Run("Filter with Range(1,100)", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint8](1, 100),
			Filter(func(value uint8, index uint) bool {
				return value <= 10
			}),
		), []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true)
	})

	t.Run("Filter with alphaberts", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled("a", "b", "cd", "kill", "p", "z", "animal"),
			Filter(func(value string, index uint) bool {
				return len(value) == 1
			}),
		), []string{"a", "b", "p", "z"}, nil, true)
	})
}

func TestFirst(t *testing.T) {
	t.Run("First with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), First[any](nil)), nil, ErrEmpty, false)
	})

	t.Run("First with default value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), First[any](nil, "hello default value")), "hello default value", nil, true)
	})

	t.Run("First with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint8](88, 99), First(func(value uint8, index uint) bool {
			return value > 0
		})), uint8(88), nil, true)
	})
}

func TestLast(t *testing.T) {
	t.Run("Last with empty value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Last[any](nil)), nil, ErrEmpty, false)
	})

	t.Run("Last with default value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Last[any](nil, 88)), 88, nil, true)
	})

	t.Run("Last with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint8](1, 72), Last[uint8](nil)), uint8(72), nil, true)
	})

	t.Run("Last with value but not matched", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint8](1, 10), Last(func(value uint8, _ uint) bool {
			return value > 10
		})), uint8(0), ErrNotFound, false)
	})
}

func TestIgnoreElements(t *testing.T) {
	t.Run("IgnoreElements with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), IgnoreElements[any]()), nil, nil, true)
	})

	t.Run("IgnoreElements with ThrownError", func(t *testing.T) {
		var err = errors.New("throw")
		checkObservableResult(t, Pipe1(ThrownError[error](func() error {
			return err
		}), IgnoreElements[error]()), nil, err, false)
	})

	t.Run("IgnoreElements with Range(1,7)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 7), IgnoreElements[uint]()), uint(0), nil, true)
	})
}

func TestSkip(t *testing.T) {
	t.Run("Skip with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(EMPTY[uint](), Skip[uint](5)), []uint{}, nil, true)
	})

	t.Run("Skip with Range(1,10)", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 10), Skip[uint](5)),
			[]uint{6, 7, 8, 9, 10}, nil, true)
	})

	t.Run("Skip with ThrownError", func(t *testing.T) {
		var err = errors.New("stop")
		checkObservableResults(t, Pipe1(ThrownError[uint](func() error {
			return err
		}), Skip[uint](5)), []uint{}, err, false)
	})

	// t.Run("Skip with Scheduled", func(t *testing.T) {
	// 	checkObservableResults(t,
	// 		Pipe1(Scheduled[any](1, 2, errors.New("stop")), Skip[any](2)),
	// 		[]any{1, 2}, nil, true)
	// })
}

func TestSkipLast(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 10), SkipLast[uint](5)), []uint{1, 2, 3, 4, 5}, nil, true)
}

func TestSkipUntil(t *testing.T) {
	t.Run("SkipUntil with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[uint](),
			SkipUntil[uint](Scheduled("a")),
		), []uint{}, nil, true)
	})

	t.Run("SkipUntil with ThrownError", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResults(t, Pipe1(
			ThrownError[uint](func() error {
				return err
			}),
			SkipUntil[uint](Scheduled("a")),
		), []uint{}, err, false)
	})

	t.Run("SkipUntil with Range(1,10)", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint](1, 10),
			SkipUntil[uint](Interval(time.Millisecond*100)),
		), []uint{}, nil, true)
	})
}

func TestSkipWhile(t *testing.T) {
	t.Run("SkipWhile until condition meet", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled("Green Arrow", "SuperMan", "Flash", "SuperGirl", "Black Canary"),
			SkipWhile(func(v string, _ uint) bool {
				return v != "SuperGirl"
			})), []string{"SuperGirl", "Black Canary"}, nil, true)
	})

	t.Run("SkipWhile until index 5", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 10), SkipWhile(func(_ uint, idx uint) bool {
			return idx != 5
		})), []uint{6, 7, 8, 9, 10}, nil, true)
	})
}

func TestTake(t *testing.T) {
	checkObservableResults(t, Pipe1(Interval(time.Millisecond), Take[uint](3)), []uint{0, 1, 2}, nil, true)
	checkObservableResults(t, Pipe1(Range[uint](1, 100), Take[uint](3)), []uint{1, 2, 3}, nil, true)
}

func TestTakeLast(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 100), TakeLast[uint](3)), []uint{98, 99, 100}, nil, true)
}

func TestTakeUntil(t *testing.T) {
	t.Run("TakeUntil with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[uint](),
			TakeUntil[uint](Scheduled("a")),
		), []uint{}, nil, true)
	})

	t.Run("TakeUntil with Interval", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Interval(time.Millisecond),
			TakeUntil[uint](Interval(time.Millisecond*5)),
		), []uint{0, 1, 2, 3}, nil, true)
	})
}

func TestTakeWhile(t *testing.T) {
	t.Run("TakeWhile with Interval", func(t *testing.T) {
		result := make([]uint, 0)
		for i := uint(0); i <= 5; i++ {
			result = append(result, i)
		}
		checkObservableResults(t, Pipe1(Interval(time.Millisecond), TakeWhile(func(v uint, _ uint) bool {
			return v <= 5
		})), result, nil, true)
	})

	t.Run("TakeWhile with Range", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 100), TakeWhile(func(v uint, _ uint) bool {
			return v >= 50
		})), []uint{}, nil, true)
	})
}

func TestThrottleTime(t *testing.T) {
	t.Run("ThrottleTime with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			EMPTY[any](),
			ThrottleTime[any](time.Millisecond),
		), nil, nil, true)
	})

	t.Run("ThrottleTime with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResult(t, Pipe1(
			ThrownError[any](func() error {
				return err
			}), ThrottleTime[any](time.Millisecond),
		), nil, err, false)
	})
}
