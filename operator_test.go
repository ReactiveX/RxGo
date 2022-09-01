package rxgo

import (
	"fmt"
	"testing"
)

func TestTake(t *testing.T) {

	// checkObservableResults(t, Pipe1(Interval(time.Second), Take[uint, uint8](3)), []uint{0, 1, 2}, nil, true)
}

func TestTakeUntil(t *testing.T) {

}

func TestTakeWhile(t *testing.T) {
	result := make([]uint, 0)
	for i := uint(50); i <= 100; i++ {
		result = append(result, i)
	}
	checkObservableResults(t, Pipe1(Range[uint](1, 100), TakeWhile(func(v uint, _ uint) bool {
		return v >= 50
	})), result, nil, true)
}

func TestTakeLast(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 100), TakeLast[uint, uint](3)), []uint{98, 99, 100}, nil, true)
}

func TestSkip(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 10), Skip[uint](5)), []uint{6, 7, 8, 9, 10}, nil, true)
}

func TestElementAt(t *testing.T) {
	t.Run("ElementAt with Default Value", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), ElementAt[any](1, 10)), 10, nil, true)
	})

	t.Run("ElementAt position 2", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 100), ElementAt[uint](2)), 3, nil, true)
	})

	t.Run("ElementAt with Error(ErrArgumentOutOfRange)", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 10), ElementAt[uint](100)), 0, ErrArgumentOutOfRange, true)
	})
}

func TestFirst(t *testing.T) {

}

func TestLast(t *testing.T) {

}

func TestFind(t *testing.T) {

}

func TestFindIndex(t *testing.T) {

}

func TestMin(t *testing.T) {

}

func TestMax(t *testing.T) {

}

func TestCount(t *testing.T) {
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
		), []string{"Number(1)", "Number(2)"}, err, true)
	})
}

func TestTap(t *testing.T) {

}

func TestSingle(t *testing.T) {

}

func TestSkipWhile(t *testing.T) {

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

func TestWithLatestFrom(t *testing.T) {
	// timer := Interval(time.Millisecond * 1)
	// Pipe2(
	// 	rxgo.Pipe1(
	// 		rxgo.Interval(time.Second*2),
	// 		rxgo.Map(func(i uint, _ uint) (string, error) {
	// 			return string(rune('A' + i)), nil
	// 		}),
	// 	),
	// 	rxgo.WithLatestFrom[string](timer),
	// 	rxgo.Take[rxgo.Tuple[string, uint], uint](10),
	// ).SubscribeSync(func(f rxgo.Tuple[string, uint]) {
	// 	log.Println("[", f.First(), f.Second(), "]")
	// }, func(err error) {}, func() {})
}

func TestTimestamp(t *testing.T) {

}

func TestToArray(t *testing.T) {
	t.Run("ToArray with Numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 5), ToArray[uint]()), []uint{1, 2, 3, 4, 5}, nil, true)
	})

	t.Run("ToArray with Numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(newObservable(func(subscriber Subscriber[string]) {
			for i := 1; i <= 5; i++ {
				subscriber.Next(string(rune('A' - 1 + i)))
			}
			subscriber.Complete()
		}), ToArray[string]()), []string{"A", "B", "C", "D", "E"}, nil, true)
	})
}
