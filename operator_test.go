package rxgo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestElementAt(t *testing.T) {
	t.Run("ElementAt with Default Value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), ElementAt[any](1, 10)), 10, nil, true)
	})

	t.Run("ElementAt position 2", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 100), ElementAt[uint](2)), 3, nil, true)
	})

	t.Run("ElementAt with Error(ErrArgumentOutOfRange)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 10), ElementAt[uint](100)), 0, ErrArgumentOutOfRange, false)
	})
}

func TestFirst(t *testing.T) {
	t.Run("First with empty value", func(t *testing.T) {
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

}

func TestFind(t *testing.T) {
	t.Run("First with empty value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Find(func(a any, u uint) bool {
			return a == nil
		})), None[any](), nil, true)
	})

	t.Run("First with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Scheduled("a", "b", "c", "d", "e"),
			Find(func(v string, u uint) bool {
				return v == "c"
			}),
		), Some("c"), nil, true)
	})
}

func TestFindIndex(t *testing.T) {
	t.Run("First with empty value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), FindIndex(func(a any, u uint) bool {
			return a == nil
		})), -1, nil, true)
	})

	t.Run("First with value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Scheduled("a", "b", "c", "d", "e"),
			FindIndex(func(v string, u uint) bool {
				return v == "c"
			}),
		), 2, nil, true)
	})
}

func TestMin(t *testing.T) {
	type human struct {
		age  int
		name string
	}
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
}

func TestMax(t *testing.T) {
	// checkObservableResult(t, Pipe1(Range[uint](1, 7), Max(func(a, b uint) int8 {
	// 	if a > b {
	// 		return 1
	// 	}
	// 	return -1
	// })), uint(7), nil, true)
}

func TestCount(t *testing.T) {
	checkObservableResult(t, Pipe1(Range[uint](0, 7), Count[uint]()), uint(7), nil, true)
	checkObservableResult(t, Pipe1(Range[uint](1, 7), Count(func(i uint, _ uint) bool {
		return i%2 == 1
	})), uint(4), nil, true)
}

func TestIgnoreElements(t *testing.T) {
	checkObservableResult(t, Pipe1(Range[uint](1, 7), IgnoreElements[uint]()), uint(0), nil, true)
}

func TestEvery(t *testing.T) {
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

func TestRepeat(t *testing.T) {

}

func TestIsEmpty(t *testing.T) {
	checkObservableResult(t, Pipe1(EMPTY[any](), IsEmpty[any]()), true, nil, true)
	checkObservableResult(t, Pipe1(Range[uint](1, 3), IsEmpty[uint]()), false, nil, true)
}

func TestDefaultIfEmpty(t *testing.T) {
	t.Run("DefaultIfEmpty with any", func(t *testing.T) {
		str := "hello world"
		checkObservableResult(t, Pipe1(EMPTY[any](), DefaultIfEmpty[any](str)), any(str), nil, true)
	})

	t.Run("DefaultIfEmpty with NonEmpty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 3), DefaultIfEmpty[uint](100)), []uint{1, 2, 3}, nil, true)
	})
}

func TestDistinctUntilChanged(t *testing.T) {
	t.Run("DistinctUntilChanged with empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), DistinctUntilChanged[any]()), nil, nil, true)
	})

	t.Run("DistinctUntilChanged with String", func(t *testing.T) {
		checkObservableResults(t,
			Pipe1(Scheduled("a", "a", "b", "a", "c", "c", "d"), DistinctUntilChanged[string]()),
			[]string{"a", "b", "a", "c", "d"}, nil, true)
	})

	t.Run("DistinctUntilChanged with Number", func(t *testing.T) {
		checkObservableResults(t,
			Pipe1(
				Scheduled(30, 31, 20, 34, 33, 29, 35, 20),
				DistinctUntilChanged(func(prev, current int) bool {
					return current <= prev
				}),
			), []int{30, 31, 34, 35}, nil, true)
	})

	t.Run("DistinctUntilChanged with Struct", func(t *testing.T) {
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

func TestFilter(t *testing.T) {

}

func TestMap(t *testing.T) {
	t.Run("Map with string", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint](1, 5),
			Map(func(v uint, _ uint) (string, error) {
				return fmt.Sprintf("Number(%d)", v), nil
			}),
		), []string{
			"Number(1)",
			"Number(2)",
			"Number(3)",
			"Number(4)",
			"Number(5)",
		}, nil, true)
	})

	t.Run("Map with Error", func(t *testing.T) {
		err := fmt.Errorf("omg")
		checkObservableResults(t, Pipe1(
			Range[uint](1, 5),
			Map(func(v uint, _ uint) (string, error) {
				if v == 3 {
					return "", err
				}
				return fmt.Sprintf("Number(%d)", v), nil
			}),
		), []string{"Number(1)", "Number(2)"}, err, false)
	})
}

func TestTap(t *testing.T) {
	t.Run("Tap with Range(1, 5)", func(t *testing.T) {
		result := make([]string, 0)
		checkObservableResults(t, Pipe1(
			Range[uint](1, 5),
			Tap(NewObserver(func(v uint) {
				result = append(result, fmt.Sprintf("Number(%v)", v))
			}, nil, nil)),
		), []uint{1, 2, 3, 4, 5}, nil, true)
		require.ElementsMatch(t, []string{
			"Number(1)",
			"Number(2)",
			"Number(3)",
			"Number(4)",
			"Number(5)",
		}, result)
	})

	t.Run("Tap with Error", func(t *testing.T) {
		var (
			err    = fmt.Errorf("An error")
			result = make([]string, 0)
		)
		checkObservableResults(t, Pipe1(
			Scheduled[any](1, err),
			Tap(NewObserver(func(v any) {
				result = append(result, fmt.Sprintf("Number(%v)", v))
			}, nil, nil)),
		), []any{1}, err, false)
		require.ElementsMatch(t, []string{"Number(1)"}, result)
	})
}

func TestSingle(t *testing.T) {

}

func TestSample(t *testing.T) {
	var (
		// result = make([]uint, 0)
		err  error
		done bool
	)
	Pipe1(
		Pipe1(Interval(time.Millisecond), Take[uint](10)),
		Sample[uint](Interval(time.Millisecond*2))).
		SubscribeSync(func(u uint) {}, func(e error) {
			err = e
		}, func() {
			done = true
		})
	require.Nil(t, err)
	require.True(t, done)
}

func TestScan(t *testing.T) {
	t.Run("Scan with initial value", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled[uint](1, 2, 3),
			Scan(func(acc, cur, _ uint) (uint, error) {
				return acc + cur, nil
			}, 10),
		), []uint{11, 13, 16}, nil, true)
	})

	t.Run("Scan with zero default value", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled[uint](1, 3, 5),
			Scan(func(acc, cur, _ uint) (uint, error) {
				return acc + cur, nil
			}, 0),
		), []uint{1, 4, 9}, nil, true)
	})
}

func TestReduce(t *testing.T) {
	t.Run("Reduce with zero default value", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled[uint](1, 3, 5),
			Reduce(func(acc, cur, _ uint) (uint, error) {
				return acc + cur, nil
			}, 0),
		), []uint{9}, nil, true)
	})
}

func TestDelay(t *testing.T) {

}

func TestThrottle(t *testing.T) {

}

func TestDebounceTime(t *testing.T) {

}

func TestRaceWith(t *testing.T) {
	// t.Run("RaceWith with Interval", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe2(
	// 		Pipe1(Interval(time.Millisecond*7), Map(func(v uint, _ uint) (string, error) {
	// 			return fmt.Sprintf("slowest -> %v", v), nil
	// 		})),
	// 		RaceWith(
	// 			Pipe1(Interval(time.Millisecond*3), Map(func(v uint, _ uint) (string, error) {
	// 				return fmt.Sprintf("fastest -> %v", v), nil
	// 			})),
	// 			Pipe1(Interval(time.Millisecond*5), Map(func(v uint, _ uint) (string, error) {
	// 				return fmt.Sprintf("average -> %v", v), nil
	// 			})),
	// 		),
	// 		Take[string](5),
	// 	),
	// 		[]string{"fastest -> 0"}, // "fastest -> 1", "fastest -> 2", "fastest -> 3", "fastest -> 4"
	// 		nil, true)
	// })
}

func TestWithLatestFrom(t *testing.T) {
	checkObservableResults(t, Pipe2(
		Interval(time.Second),
		WithLatestFrom[uint](Scheduled("a", "v")),
		Take[Tuple[uint, string]](3),
	), []Tuple[uint, string]{
		NewTuple[uint](0, "v"),
		NewTuple[uint](1, "v"),
		NewTuple[uint](2, "v"),
	}, nil, true)
}

// func TestOnErrorResumeNext(t *testing.T) {
// 	// t.Run("OnErrorResumeNext with error", func(t *testing.T) {
// 	// 	checkObservableResults(t, Pipe1(Scheduled[any](1, 2, 3, fmt.Errorf("error"), 5), OnErrorResumeNext[any]()),
// 	// 		[]any{1, 2, 3, 5}, nil, true)
// 	// })

// 	t.Run("OnErrorResumeNext with no error", func(t *testing.T) {
// 		checkObservableResults(t, Pipe1(Scheduled(1, 2, 3, 4, 5), OnErrorResumeNext[int]()),
// 			[]int{1, 2, 3, 4, 5}, nil, true)
// 	})
// }

func TestToArray(t *testing.T) {
	t.Run("ToArray with Numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 5), ToArray[uint]()), []uint{1, 2, 3, 4, 5}, nil, true)
	})

	t.Run("ToArray with Numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(newObservable(func(subscriber Subscriber[string]) {
			for i := 1; i <= 5; i++ {
				subscriber.Send() <- NextNotification(string(rune('A' - 1 + i)))
			}
			subscriber.Send() <- CompleteNotification[string]()
		}), ToArray[string]()), []string{"A", "B", "C", "D", "E"}, nil, true)
	})
}
