package rxgo

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	// t.Run("Buffer with EMPTY", func(t *testing.T) {
	// 	checkObservableResult(t, Pipe1(
	// 		EMPTY[uint](),
	// 		Buffer[uint](Scheduled("a")),
	// 	), nil, nil, true)
	// })

	t.Run("Buffer with Scheduled", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Interval(time.Millisecond),
			Buffer[uint](Scheduled("a")),
		), []uint{}, nil, true)
	})

	// t.Run("Buffer with Interval", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe2(
	// 		Interval(time.Millisecond*10),
	// 		Buffer[uint](Interval(time.Millisecond*100)),
	// 		Take[[]uint](3),
	// 	), [][]uint{
	// 		{0, 1, 2, 3, 4, 5, 6, 7, 8},
	// 		{9, 10, 11, 12, 13, 14, 15, 16, 17},
	// 		{18, 19, 20, 21, 22, 23, 24, 25, 26},
	// 	}, nil, true)
	// })
}

func TestBufferCount(t *testing.T) {
	// t.Run("BufferCount with EMPTY", func(t *testing.T) {
	// 	checkObservableResult(t, Pipe1(
	// 		EMPTY[uint](),
	// 		BufferCount[uint](2),
	// 	), nil, nil, true)
	// })

	// t.Run("BufferCount with Range(1,7)", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Range[uint](1, 7),
	// 		BufferCount[uint](2),
	// 	), [][]uint{
	// 		{1, 2},
	// 		{3, 4},
	// 		{5, 6},
	// 		{7},
	// 	}, nil, true)
	// })

	// t.Run("BufferCount with Range(1,7)", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Range[uint](0, 7),
	// 		BufferCount[uint](3, 1),
	// 	), [][]uint{
	// 		{0, 1, 2},
	// 		{1, 2, 3},
	// 		{2, 3, 4},
	// 		{3, 4, 5},
	// 		{4, 5, 6},
	// 		{5, 6, 7},
	// 		{5, 6},
	// 		{7},
	// 	}, nil, true)
	// })
}

func TestConcatMap(t *testing.T) {
	t.Run("ConcatMap with error on upstream", func(t *testing.T) {
		var err = fmt.Errorf("throw")
		checkObservableResults(t, Pipe1(
			Scheduled[any]("z", err, "q"),
			ConcatMap(func(x any, i uint) Observable[string] {
				return Pipe2(
					Interval(time.Millisecond),
					Map(func(y, _ uint) (string, error) {
						return fmt.Sprintf("%v[%d]", x, y), nil
					}),
					Take[string](2),
				)
			}),
		), []string{"z[0]", "z[1]"}, err, false)
	})

	t.Run("ConcatMap with conditional ThrownError", func(t *testing.T) {
		var err = fmt.Errorf("throw")

		mapTo := func(v string, i uint) string {
			return fmt.Sprintf("%s[%d]", v, i)
		}

		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) Observable[string] {
				if i == 0 {
					return Scheduled(mapTo(x, i), mapTo(x, i), mapTo(x, i))
				}

				return ThrownError[string](func() error {
					return err
				})
			}),
		), []string{"z[0]", "z[0]", "z[0]"}, err, false)
	})

	t.Run("ConcatMap with ThrownError on return stream", func(t *testing.T) {
		var err = fmt.Errorf("throw")

		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) Observable[string] {
				return ThrownError[string](func() error {
					return err
				})
			}),
		), []string{}, err, false)
	})

	t.Run("ConcatMap with Interval + Map which return error", func(t *testing.T) {
		var err = errors.New("nopass")
		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) Observable[string] {
				return Pipe2(
					Interval(time.Second),
					Map(func(y, idx uint) (string, error) {
						if idx == 1 {
							return "", err
						}
						return fmt.Sprintf("%s[%d]", x, y), nil
					}),
					Take[string](2),
				)
			}),
		), []string{"z[0]"}, err, false)
	})

	t.Run("ConcatMap with Interval[string]", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled[uint](1, 2, 4),
			ConcatMap(func(x uint, i uint) Observable[string] {
				return Pipe2(
					Interval(time.Millisecond),
					Map(func(y, _ uint) (string, error) {
						return fmt.Sprintf("x -> %d, y -> %d", x, y), nil
					}),
					Take[string](3),
				)
			})), []string{
			"x -> 1, y -> 0",
			"x -> 1, y -> 1",
			"x -> 1, y -> 2",
			"x -> 2, y -> 0",
			"x -> 2, y -> 1",
			"x -> 2, y -> 2",
			"x -> 4, y -> 0",
			"x -> 4, y -> 1",
			"x -> 4, y -> 2",
		}, nil, true)
	})
}

func TestExhaustMap(t *testing.T) {
	t.Run("ExhaustMap with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[any](),
			ExhaustMap(func(x any, _ uint) Observable[string] {
				return Pipe1(
					Range[uint](88, 90),
					Map(func(y, _ uint) (string, error) {
						return fmt.Sprintf("%v:%d", x, y), nil
					}),
				)
			}),
		), []string{}, nil, true)
	})

	t.Run("ExhaustMap with error", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[any](),
			ExhaustMap(func(x any, _ uint) Observable[string] {
				return Pipe1(
					Range[uint](88, 90),
					Map(func(y, _ uint) (string, error) {
						return fmt.Sprintf("%v:%d", x, y), nil
					}),
				)
			}),
		), []string{}, nil, true)
	})

	// t.Run("ExhaustMap with Interval", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe2(
	// 		Interval(time.Millisecond*10),
	// 		ExhaustMap(func(x uint, _ uint) Observable[string] {
	// 			return Pipe1(
	// 				Range[uint](88, 90),
	// 				Map(func(y, _ uint) (string, error) {
	// 					return fmt.Sprintf("%02d:%d", x, y), nil
	// 				}),
	// 			)
	// 		}),
	// 		Take[string](3),
	// 	), []string{"00:88", "00:89", "00:90"}, nil, true)
	// })
}

func TestExhaustAll(t *testing.T) {
	// t.Run("ExhaustAll with Interval", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe3(
	// 		Interval(time.Millisecond*100),
	// 		Map(func(v uint, _ uint) (Observable[uint], error) {
	// 			return Range[uint](88, 10), nil
	// 		}),
	// 		ExhaustAll[uint](),
	// 		Take[uint](15),
	// 	), []uint{
	// 		88, 89, 90, 91, 92, 93, 94, 95, 96, 97,
	// 		88, 89, 90, 91, 92,
	// 	}, nil, true)
	// })
}

func TestGroupBy(t *testing.T) {
	// t.Run("GroupBy with EMPTY", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		EMPTY[any](),
	// 		GroupBy[any, any](),
	// 	), []any{}, nil, true)
	// })

	// type js struct {
	// 	id   uint
	// 	name string
	// }

	// t.Run("GroupBy with objects", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Scheduled(js{1, "JavaScript"}, js{2, "Parcel"}),
	// 		GroupBy[js](func(v js) string {
	// 			return v.name
	// 		}),
	// 	), []any{}, nil, true)
	// })
}

func TestMap(t *testing.T) {
	t.Run("Map with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[any](),
			Map(func(v any, _ uint) (any, error) {
				return v, nil
			}),
		), []any{}, nil, true)
	})

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

func TestMergeMap(t *testing.T) {
	t.Run("MergeMap with EMPTY", func(t *testing.T) {

	})

	// t.Run("MergeMap with complete", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Scheduled("a", "b", "v"),
	// 		MergeMap(func(x string, i uint) Observable[Tuple[string, uint]] {
	// 			return Pipe2(
	// 				Interval(time.Millisecond),
	// 				Map(func(y, _ uint) (Tuple[string, uint], error) {
	// 					return NewTuple(x, y), nil
	// 				}),
	// 				Take[Tuple[string, uint]](3),
	// 			)
	// 		}),
	// 	), []Tuple[string, uint]{
	// 		NewTuple[string, uint]("a", 0),
	// 		NewTuple[string, uint]("b", 0),
	// 		NewTuple[string, uint]("v", 0),
	// 		NewTuple[string, uint]("a", 1),
	// 		NewTuple[string, uint]("b", 1),
	// 		NewTuple[string, uint]("v", 1),
	// 		NewTuple[string, uint]("a", 2),
	// 		NewTuple[string, uint]("b", 2),
	// 		NewTuple[string, uint]("v", 2),
	// 	}, nil, true)
	// })

	// t.Run("MergeMap with error", func(t *testing.T) {
	// 	var (
	// 		result = make([]Tuple[string, uint], 0)
	// 		failed = errors.New("failed")
	// 		err    error
	// 		done   bool
	// 	)
	// 	Pipe1(
	// 		Scheduled("a", "b", "v"),
	// 		MergeMap(func(x string, i uint) Observable[Tuple[string, uint]] {
	// 			return Pipe2(
	// 				Interval(time.Millisecond),
	// 				Map(func(y, idx uint) (Tuple[string, uint], error) {
	// 					if idx > 3 {
	// 						return nil, failed
	// 					}
	// 					return NewTuple(x, y), nil
	// 				}),
	// 				Take[Tuple[string, uint]](5),
	// 			)
	// 		}),
	// 	).SubscribeSync(func(s Tuple[string, uint]) {
	// 		result = append(result, s)
	// 	}, func(e error) {
	// 		err = e
	// 	}, func() {
	// 		done = true
	// 	})
	// 	require.True(t, len(result) == 9)
	// 	require.Equal(t, failed, err)
	// 	require.False(t, done)
	// })
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

func TestPartition(t *testing.T) {

}

func TestPairWise(t *testing.T) {
	t.Run("PairWise with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(EMPTY[any](), PairWise[any]()),
			[]Tuple[any, any]{}, nil, true)
	})

	t.Run("PairWise with error", func(t *testing.T) {
		var err = errors.New("throw")
		checkObservableResults(t, Pipe1(Scheduled[any]("j", "k", err), PairWise[any]()),
			[]Tuple[any, any]{NewTuple[any, any]("j", "k")}, err, false)
	})

	t.Run("PairWise with numbers", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 5), PairWise[uint]()),
			[]Tuple[uint, uint]{
				NewTuple[uint, uint](1, 2),
				NewTuple[uint, uint](2, 3),
				NewTuple[uint, uint](3, 4),
				NewTuple[uint, uint](4, 5),
			}, nil, true)
	})

	t.Run("PairWise with alphaberts", func(t *testing.T) {
		checkObservableResults(t, Pipe1(newObservable(func(subscriber Subscriber[string]) {
			for i := 1; i <= 5; i++ {
				subscriber.Send() <- Next(string(rune('A' - 1 + i)))
			}
			subscriber.Send() <- Complete[string]()
		}), PairWise[string]()), []Tuple[string, string]{
			NewTuple("A", "B"),
			NewTuple("B", "C"),
			NewTuple("C", "D"),
			NewTuple("D", "E"),
		}, nil, true)
	})
}
