package rxgo

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	// t.Run("Buffer with Empty", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Empty[uint](),
	// 		Buffer[uint](Of2("a")),
	// 	), [][]uint{}, nil, true)
	// })

	// t.Run("Buffer with error", func(t *testing.T) {
	// 	var err = fmt.Errorf("failed")
	// 	checkObservableResult(t, Pipe1(
	// 		Throw[string](func() error {
	// 			return err
	// 		}),
	// 		Buffer[string](Of2("a")),
	// 	), []string(nil), err, false)
	// })

	// t.Run("Buffer with Empty should throw ErrEmpty", func(t *testing.T) {
	// 	checkObservableResult(t, Pipe2(
	// 		Empty[string](),
	// 		ThrowIfEmpty[string](),
	// 		Buffer[string](Interval(time.Millisecond)),
	// 	), nil, ErrEmpty, false)
	// })

	// t.Run("Buffer with Interval", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Of2("a", "b", "c", "d", "e"),
	// 		Buffer[string](Interval(time.Millisecond*500)),
	// 	), [][]string{
	// 		{"a", "b", "c", "d", "e"},
	// 	}, nil, true)
	// })
}

func TestBufferCount(t *testing.T) {
	t.Run("BufferCount with Empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Empty[uint](),
			BufferCount[uint](2),
		), nil, nil, true)
	})

	t.Run("BufferCount with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResult(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			BufferCount[any](2),
		), nil, err, false)
	})

	t.Run("BufferCount with error", func(t *testing.T) {
		var (
			of  = Of2("a", "h", "j", "o", "k", "e", "r", "!")
			err = errors.New("failed")
		)

		checkObservableResults(t, Pipe2(
			of,
			Map(func(v string, _ uint) (string, error) {
				if strings.EqualFold(v, "e") {
					return "", err
				}
				return v, nil
			}),
			BufferCount[string](2),
		), [][]string{{"a", "h"}, {"j", "o"}}, err, false)

		checkObservableResults(t, Pipe2(
			of,
			Map(func(v string, _ uint) (string, error) {
				if strings.EqualFold(v, "r") {
					return "", err
				}
				return v, nil
			}),
			BufferCount[string](2),
		), [][]string{{"a", "h"}, {"j", "o"}, {"k", "e"}}, err, false)
	})

	t.Run("BufferCount with Range(1,7)", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint](1, 7),
			BufferCount[uint](2),
		), [][]uint{{1, 2}, {3, 4}, {5, 6}, {7}}, nil, true)
	})

	t.Run("BufferCount with Range(1,7)", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Range[uint](0, 7),
			BufferCount[uint](3, 1),
		), [][]uint{
			{0, 1, 2},
			{1, 2, 3},
			{2, 3, 4},
			{3, 4, 5},
			{4, 5, 6},
			{5, 6},
			{6},
		}, nil, true)
	})
}

func TestBufferTime(t *testing.T) {
	t.Run("BufferTime with Empty", func(t *testing.T) {
		checkObservableHasResults(t, Pipe1(
			Empty[string](),
			BufferTime[string](time.Millisecond*500),
		), true, nil, true)
	})

	t.Run("BufferTime with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableHasResults(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			BufferTime[any](time.Millisecond*500),
		), false, err, false)
	})

	t.Run("BufferTime with Of", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2("a", "z", "j", "p"),
			BufferTime[string](time.Millisecond*500),
		), [][]string{{"a", "z", "j", "p"}}, nil, true)
	})

	t.Run("BufferTime with Interval", func(t *testing.T) {
		checkObservableHasResults(t, Pipe2(
			Interval(time.Millisecond*100),
			BufferTime[uint](time.Millisecond*300),
			Take[[]uint](3),
		), true, nil, true)
	})
}

func TestBufferToggle(t *testing.T) {
	toggleFunc := BufferToggle[uint](Interval(time.Second), func(v uint) Observable[uint] {
		if v%2 == 0 {
			return Interval(time.Millisecond * 500)
		}
		return Empty[uint]()
	})

	t.Run("BufferToggle with Empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Empty[uint](),
			toggleFunc,
		), [][]uint{}, nil, true)
	})

	t.Run("BufferToggle with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResults(t, Pipe1(
			Throw[uint](func() error {
				return err
			}),
			toggleFunc,
		), [][]uint{}, err, false)
	})
}

func TestBufferWhen(t *testing.T) {
	// t.Run("BufferWhen with Empty", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Empty[string](),
	// 		BufferWhen[string](func() Observable[string] {
	// 			return Of2("a")
	// 		}),
	// 	), [][]string{}, nil, true)
	// })

	t.Run("BufferWhen with Of", func(t *testing.T) {
		values := []string{"I", "#@$%^&*", "XD", "Z"}
		checkObservableResults(t, Pipe1(
			Of2(values[0], values[1:]...),
			BufferWhen[string](func() Observable[uint] {
				return Interval(time.Millisecond * 10)
			}),
		), [][]string{values}, nil, true)
	})

	// t.Run("BufferWhen with values", func(t *testing.T) {
	// 	const interval = 500
	// 	checkObservableHasResults(t, Pipe2(
	// 		Interval(time.Millisecond*interval),
	// 		BufferWhen[uint](func() Observable[uint] {
	// 			return Interval(time.Millisecond * interval * 2)
	// 		}),
	// 		Take[[]uint](4),
	// 	), nil, true)
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

	t.Run("ConcatMap with conditional Throw", func(t *testing.T) {
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

				return Throw[string](func() error {
					return err
				})
			}),
		), []string{"z[0]", "z[0]", "z[0]"}, err, false)
	})

	t.Run("ConcatMap with Throw on return stream", func(t *testing.T) {
		var err = fmt.Errorf("throw")

		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) Observable[string] {
				return Throw[string](func() error {
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
					Interval(time.Millisecond*10),
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
	t.Run("ExhaustMap with Empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Empty[any](),
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
			Empty[any](),
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

func TestExpand(t *testing.T) {
	// t.Run("Expand with Empty", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe2(
	// 		Range[uint8](1, 5),
	// 		Expand(func(v uint8, _ uint) Observable[string] {
	// 			return Of2(fmt.Sprintf("Number(%d)", v))
	// 		}),
	// 		Take[Either[uint8, string]](5),
	// 	), []Either[uint8, string]{
	// 		Left[uint8, string](1),
	// 		Left[uint8, string](2),
	// 		Left[uint8, string](3),
	// 		Left[uint8, string](4),
	// 		Left[uint8, string](5),
	// 	}, nil, true)
	// })
}

func TestGroupBy(t *testing.T) {
	// t.Run("GroupBy with Empty", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe1(
	// 		Empty[any](),
	// 		GroupBy[any, any](),
	// 	), []any{}, nil, true)
	// })

	// type js struct {
	// 	id   uint
	// 	name string
	// }

	// t.Run("GroupBy with objects", func(t *testing.T) {
	// 	checkObservableResults(t, Pipe2(
	// 		Of2(js{1, "JavaScript"}, js{2, "Parcel"}, js{2, "Webpack"}, js{1, "TypeScript"}, js{3, "TSLint"}),
	// 		GroupBy(func(v js) string {
	// 			return v.name
	// 		}),
	// 		MergeMap(func(group GroupedObservable[string, js], index uint) Observable[[]js] {
	// 			return Pipe1(
	// 				group.(Observable[js]),
	// 				Reduce(func(acc []js, value js, _ uint) ([]js, error) {
	// 					acc = append(acc, value)
	// 					return acc, nil
	// 				}, []js{}),
	// 			)
	// 		}),
	// 	), [][]js{
	// 		{{1, "JavaScript"}, {1, "TypeScript"}},
	// 		{{2, "Parcel"}, {2, "Webpack"}},
	// 		{{3, "TSLint"}},
	// 	}, nil, true)
	// })
}

func TestMap(t *testing.T) {
	t.Run("Map with Empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Empty[any](),
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
	// t.Run("MergeMap with Empty", func(t *testing.T) {
	// 	checkObservableHasResults(t, Pipe1(
	// 		Empty[string](),
	// 		MergeMap(func(x string, i uint) Observable[Tuple[string, uint]] {
	// 			return Pipe2(
	// 				Interval(time.Millisecond),
	// 				Map(func(y, _ uint) (Tuple[string, uint], error) {
	// 					return NewTuple(x, y), nil
	// 				}),
	// 				Take[Tuple[string, uint]](3),
	// 			)
	// 		}),
	// 	), false, nil, true)
	// })

	// t.Run("MergeMap with inner Empty", func(t *testing.T) {
	// 	checkObservableHasResults(t, Pipe1(
	// 		Of2("a", "b", "v"),
	// 		MergeMap(func(x string, i uint) Observable[any] {
	// 			return Empty[any]()
	// 		}),
	// 	), false, nil, true)
	// })

	// t.Run("MergeMap with complete", func(t *testing.T) {
	// 	checkObservableHasResults(t, Pipe1(
	// 		Of2("a", "b", "v"),
	// 		MergeMap(func(x string, i uint) Observable[Tuple[string, uint]] {
	// 			return Pipe2(
	// 				Interval(time.Millisecond),
	// 				Map(func(y, _ uint) (Tuple[string, uint], error) {
	// 					return NewTuple(x, y), nil
	// 				}),
	// 				Take[Tuple[string, uint]](3),
	// 			)
	// 		}),
	// 	), true, nil, true)
	// })

	// t.Run("MergeMap with error", func(t *testing.T) {
	// 	var (
	// 		err = errors.New("failed")
	// 	)
	// 	checkObservableHasResults(t, Pipe1(
	// 		Of2("a", "b", "v"),
	// 		MergeMap(func(x string, i uint) Observable[Tuple[string, uint]] {
	// 			return Pipe2(
	// 				Interval(time.Millisecond),
	// 				Map(func(y, idx uint) (Tuple[string, uint], error) {
	// 					if idx > 3 {
	// 						return nil, err
	// 					}
	// 					return NewTuple(x, y), nil
	// 				}),
	// 				Take[Tuple[string, uint]](5),
	// 			)
	// 		}),
	// 	), true, err, true)
	// 	// 	var (
	// 	// 		result = make([]Tuple[string, uint], 0)
	// 	//
	// 	// 		err    error
	// 	// 		done   bool
	// 	// 	)
	// 	// 	Pipe1(
	// 	// 		Scheduled("a", "b", "v"),

	// 	// 	).SubscribeSync(func(s Tuple[string, uint]) {
	// 	// 		result = append(result, s)
	// 	// 	}, func(e error) {
	// 	// 		err = e
	// 	// 	}, func() {
	// 	// 		done = true
	// 	// 	})
	// 	// 	require.True(t, len(result) == 9)
	// 	// 	require.Equal(t, failed, err)
	// 	// 	require.False(t, done)
	// })
}

func TestMergeScan(t *testing.T) {
	t.Run("MergeScan with Empty", func(t *testing.T) {
		checkObservableHasResults(t, Pipe1(
			Empty[any](),
			MergeScan(func(acc string, v any, _ uint) Observable[string] {
				return Of2(fmt.Sprintf("%v %s", v, acc))
			}, "hello ->"),
		), false, nil, true)
	})

	t.Run("MergeScan with outer error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableHasResults(t, Pipe1(
			Throw[any](func() error {
				return err
			}),
			MergeScan(func(acc string, v any, _ uint) Observable[string] {
				return Of2(fmt.Sprintf("%v %s", v, acc))
			}, "hello ->"),
		), false, err, false)
	})

	t.Run("MergeScan with inner error", func(t *testing.T) {
		var err = errors.New("cannot more than 5")
		checkObservableHasResults(t, Pipe1(
			Interval(time.Millisecond*5),
			MergeScan(func(acc string, v uint, _ uint) Observable[string] {
				if v > 5 {
					return Throw[string](func() error {
						return err
					})
				}
				return Of2(fmt.Sprintf("%v %s", v, acc))
			}, "hello ->"),
		), true, err, false)
	})

	t.Run("MergeScan with Range(1, 5)", func(t *testing.T) {
		checkObservableHasResults(t, Pipe1(
			Range[uint8](1, 5),
			MergeScan(func(acc uint, v uint8, _ uint) Observable[uint] {
				return Of2(acc + uint(v))
			}, uint(88)),
		), true, nil, true)
	})

	t.Run("MergeScan with Range(1, 5) and Interval", func(t *testing.T) {
		checkObservableHasResults(t, Pipe2(
			Range[uint8](1, 5),
			MergeScan(func(acc uint, v uint8, _ uint) Observable[uint] {
				return Interval(time.Millisecond)
			}, uint(88)),
			Take[uint](2),
		), true, nil, true)
	})
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

func TestSwitchMap(t *testing.T) {
	// t.Run("SwitchMap with Empty", func(t *testing.T) {
	// 	checkObservableHasResults(t, Pipe1(
	// 		Empty[uint](),
	// 		SwitchMap(func(_, _ uint) Observable[string] {
	// 			return Of2("hello", "world", "!!")
	// 		}),
	// 	), false, nil, true)
	// })

	// t.Run("SwitchMap with Empty and inner error", func(t *testing.T) {
	// 	var err = errors.New("throw")
	// 	checkObservableHasResults(t, Pipe1(
	// 		Empty[uint](),
	// 		SwitchMap(func(_, _ uint) Observable[any] {
	// 			return Throw[any](func() error {
	// 				return err
	// 			})
	// 		}),
	// 	), false, nil, true)
	// })

	// t.Run("SwitchMap with error", func(t *testing.T) {
	// 	var err = errors.New("throw")
	// 	checkObservableHasResults(t, Pipe1(
	// 		Throw[any](func() error {
	// 			return err
	// 		}),
	// 		SwitchMap(func(v any, _ uint) Observable[any] {
	// 			return Of2(v)
	// 		}),
	// 	), false, err, false)
	// })

	// t.Run("SwitchMap with inner error", func(t *testing.T) {
	// 	var err = errors.New("throw")
	// 	checkObservableHasResults(t, Pipe1(
	// 		Range[uint](1, 5),
	// 		SwitchMap(func(_, _ uint) Observable[any] {
	// 			return Throw[any](func() error {
	// 				return err
	// 			})
	// 		}),
	// 	), false, err, false)
	// })

	// t.Run("SwitchMap with inner error", func(t *testing.T) {
	// 	var err = errors.New("throw")
	// 	checkObservableHasResults(t, Pipe1(
	// 		Throw[any](func() error {
	// 			return err
	// 		}),
	// 		SwitchMap(func(v any, _ uint) Observable[any] {
	// 			return Of2(v)
	// 		}),
	// 	), false, err, false)
	// })

	// t.Run("SwitchMap with values", func(t *testing.T) {
	// 	checkObservableHasResults(t, Pipe1(
	// 		Range[uint](1, 100),
	// 		SwitchMap(func(v, _ uint) Observable[string] {
	// 			arr := make([]string, 0)
	// 			for i := uint(0); i < v; i++ {
	// 				arr = append(arr, fmt.Sprintf("%d{%d}", v, i))
	// 			}
	// 			return Of2(arr[0], arr[1:]...)
	// 		}),
	// 	), true, nil, true)
	// })
}

func TestPairWise(t *testing.T) {
	t.Run("PairWise with Empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Empty[any](), PairWise[any]()),
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
