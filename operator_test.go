package rxgo

import (
	"fmt"
	"testing"
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
	checkObservableResult(t, Pipe1(
		Range[uint](1, 7),
		IgnoreElements[uint](),
	), uint(0), nil, true)
}

func TestEvery(t *testing.T) {

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
	t.Run("Tap with", func(t *testing.T) {
		// err := fmt.Errorf("omg")
		// checkObservableResults(t, Pipe1(
		// 	Range[uint](1, 5),
		// 	Tap(func(v uint, _ uint) (string, error) {
		// 		if v == 3 {
		// 			return "", err
		// 		}
		// 		return fmt.Sprintf("Number(%d)", v), nil
		// 	}),
		// ), []string{"Number(1)", "Number(2)"}, err, false)
	})
}

func TestSingle(t *testing.T) {

}

func TestConcatMap(t *testing.T) {

}

func TestExhaustMap(t *testing.T) {

}

func TestScan(t *testing.T) {

}

func TestReduce(t *testing.T) {

}

func TestDelay(t *testing.T) {

}

func TestThrottle(t *testing.T) {

}

// func TestWithLatestFrom(t *testing.T) {
// 	// timer := Interval(time.Millisecond * 1)
// 	// Pipe2(
// 	// 	rxgo.Pipe1(
// 	// 		rxgo.Interval(time.Second*2),
// 	// 		rxgo.Map(func(i uint, _ uint) (string, error) {
// 	// 			return string(rune('A' + i)), nil
// 	// 		}),
// 	// 	),
// 	// 	rxgo.WithLatestFrom[string](timer),
// 	// 	rxgo.Take[rxgo.Tuple[string, uint], uint](10),
// 	// ).SubscribeSync(func(f rxgo.Tuple[string, uint]) {
// 	// 	log.Println("[", f.First(), f.Second(), "]")
// 	// }, func(err error) {}, func() {})
// }

func TestTimestamp(t *testing.T) {

}

func TestToArray(t *testing.T) {
	t.Run("ToArray with Numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 5), ToArray[uint]()), []uint{1, 2, 3, 4, 5}, nil, true)
	})

	t.Run("ToArray with Numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(newObservable(func(subscriber Subscriber[string]) {
			for i := 1; i <= 5; i++ {
				subscriber.Send() <- newData(string(rune('A' - 1 + i)))
			}
			subscriber.Send() <- newComplete[string]()
		}), ToArray[string]()), []string{"A", "B", "C", "D", "E"}, nil, true)
	})
}
