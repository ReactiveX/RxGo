package rxgo

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCombineLatestWith(t *testing.T) {
	t.Run("CombineLatestWith EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			EMPTY[any](),
			CombineLatestWith(
				Scheduled[any]("end"),
				Pipe2(
					Interval(time.Millisecond*100),
					Map(func(v, _ uint) (any, error) {
						return v, nil
					}),
					Take[any](10),
				),
			),
		), nil, nil, true)
	})

	t.Run("CombineLatestWith with values", func(t *testing.T) {
		checkObservableResults(t, Pipe2(
			Interval(time.Millisecond*500),
			CombineLatestWith(
				Range[uint](1, 10),
				Scheduled[uint](88),
			),
			Take[[]uint](1),
		), [][]uint{{0, 10, 88}}, nil, true)
	})
}

func TestForkJoin(t *testing.T) {
	t.Run("ForkJoin with one EMPTY", func(t *testing.T) {
		// ForkJoin only capture all latest value from every stream
		checkObservableResult(t, ForkJoin(
			EMPTY[any](),
			Scheduled[any]("j", "k", "end"),
			Pipe1(Range[uint](1, 10), Map(func(v, _ uint) (any, error) {
				return v, nil
			})),
		), []any{nil, "end", uint(10)}, nil, true)
	})

	t.Run("ForkJoin with all EMPTY", func(t *testing.T) {
		checkObservableResult(t, ForkJoin(
			EMPTY[uint](),
			EMPTY[uint](),
			EMPTY[uint](),
		), []uint{0, 0, 0}, nil, true)
	})

	t.Run("ForkJoin with error observable", func(t *testing.T) {
		var err = fmt.Errorf("failed")
		checkObservableResult(t, ForkJoin(
			Scheduled[uint](1, 88, 2, 7215251),
			Pipe1(Interval(time.Millisecond*10), Map(func(v, _ uint) (uint, error) {
				return v, err
			})),
			Interval(time.Millisecond*100),
		), nil, err, false)
	})

	t.Run("ForkJoin with multiple error", func(t *testing.T) {
		createErr := func(index uint) error {
			return fmt.Errorf("failed at %d", index)
		}
		checkObservableResultWithAnyError(t, ForkJoin(
			ThrowError[string](func() error {
				return createErr(1)
			}),
			ThrowError[string](func() error {
				return createErr(2)
			}),
			ThrowError[string](func() error {
				return createErr(3)
			}),
			Scheduled("a"),
		), nil, []error{createErr(1), createErr(2), createErr(3)}, false)
	})

	t.Run("ForkJoin with complete", func(t *testing.T) {
		checkObservableResult(t, ForkJoin(
			Scheduled[uint](1, 88, 2, 7215251),
			Pipe1(Interval(time.Millisecond*10), Take[uint](3)),
		), []uint{7215251, 2}, nil, true)
	})
}

func TestMergeWith(t *testing.T) {
	t.Run("MergeWith all EMTPY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[any](),
			MergeWith(
				EMPTY[any](),
				EMPTY[any](),
			),
		), []any{}, nil, true)
	})

	t.Run("MergeWith multiple EMTPY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2[any]("a", "b", "q", "j", "z"),
			MergeWith(
				EMPTY[any](),
				EMPTY[any](),
			),
		), []any{"a", "b", "q", "j", "z"}, nil, true)
	})

	t.Run("MergeWith Interval", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Pipe2(
				Interval(time.Millisecond),
				Take[uint](3),
				Map(func(v uint, _ uint) (string, error) {
					return fmt.Sprintf("a -> %v", v), nil
				}),
			),
			MergeWith(
				Pipe2(
					Interval(time.Millisecond*300),
					Take[uint](5),
					Map(func(v uint, _ uint) (string, error) {
						return fmt.Sprintf("b -> %v", v), nil
					}),
				),
				EMPTY[string](),
			),
		), []string{
			"a -> 0", "a -> 1", "a -> 2",
			"b -> 0", "b -> 1", "b -> 2", "b -> 3", "b -> 4",
		}, nil, true)
	})

	t.Run("MergeWith Of", func(t *testing.T) {
		checkObservableHasResults(t, Pipe1(
			Of2[any]("a", "b", "q", "j", "z"),
			MergeWith(Pipe1(
				Range[uint](1, 10),
				Map(func(v, _ uint) (any, error) {
					return any(v), nil
				}),
			)),
		), true, nil, true)
	})

	t.Run("MergeWith error", func(t *testing.T) {
		var err = errors.New("cannot more than 5")
		checkObservableHasResults(t, Pipe1(
			Of2[any]("a", "b", "q", "j", "z"),
			MergeWith(Pipe1(
				Range[uint](1, 10),
				Map(func(v, _ uint) (any, error) {
					if v > 5 {
						return nil, err
					}
					return any(v), nil
				}),
			)),
		), true, err, false)
	})

	t.Run("MergeWith all errors", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableHasResults(t, Pipe1(
			ThrowError[any](func() error {
				return err
			}),
			MergeWith(
				ThrowError[any](func() error {
					return err
				}),
				ThrowError[any](func() error {
					return err
				}),
			),
		), false, err, false)
	})
}

func TestPartition(t *testing.T) {
	t.Run("Partition with EMPTY", func(t *testing.T) {})

	t.Run("Partition with error", func(t *testing.T) {})

	t.Run("Partition", func(t *testing.T) {})
}

func TestRaceWith(t *testing.T) {
	t.Run("RaceWith with EMPTY", func(t *testing.T) {})

	t.Run("RaceWith with error", func(t *testing.T) {})

	t.Run("RaceWith with Interval", func(t *testing.T) {
		//		checkObservableResults(t, Pipe2(
		//			Pipe1(Interval(time.Millisecond*7), Map(func(v uint, _ uint) (string, error) {
		//				return fmt.Sprintf("slowest -> %v", v), nil
		//			})),
		//			RaceWith(
		//				Pipe1(Interval(time.Millisecond*3), Map(func(v uint, _ uint) (string, error) {
		//					return fmt.Sprintf("fastest -> %v", v), nil
		//				})),
		//				Pipe1(Interval(time.Millisecond*5), Map(func(v uint, _ uint) (string, error) {
		//					return fmt.Sprintf("average -> %v", v), nil
		//				})),
		//			),
		//			Take[string](5),
		//		),
		//			[]string{"fastest -> 0"}, // "fastest -> 1", "fastest -> 2", "fastest -> 3", "fastest -> 4"
		//			nil, true)
	})
}

func TestZip(t *testing.T) {
	t.Run("Zip with all EMPTY", func(t *testing.T) {
		checkObservableResults(t, Zip(
			EMPTY[any](),
			EMPTY[any](),
			EMPTY[any](),
		), [][]any{}, nil, true)
	})

	t.Run("Zip with error", func(t *testing.T) {
		var err = errors.New("stop")
		checkObservableResults(t, Pipe1(
			Zip(
				Of2[any](27, 25, 29),
				Of2[any]("Foo", "Bar", "Beer"),
				Of2[any](true, true, false),
			),
			Map(func(v []any, i uint) ([]any, error) {
				if i >= 2 {
					return nil, err
				}
				return v, nil
			}),
		), [][]any{
			{27, "Foo", true},
			{25, "Bar", true},
		}, err, false)
	})

	t.Run("Zip with EMPTY and Of", func(t *testing.T) {
		checkObservableResults(t, Zip(
			EMPTY[any](),
			Of2[any]("Foo", "Bar", "Beer"),
			Of2[any](true, true, false),
		), [][]any{}, nil, true)
	})

	t.Run("Zip with Of (not tally)", func(t *testing.T) {
		checkObservableResults(t, Zip(
			Of2[any](27, 25, 29),
			Of2[any]("Foo", "Beer"),
			Of2[any](true, true, false),
		), [][]any{
			{27, "Foo", true},
			{25, "Beer", true},
		}, nil, true)
	})

	t.Run("Zip with Of (tally)", func(t *testing.T) {
		checkObservableResults(t, Zip(
			Scheduled[any](27, 25, 29),
			Scheduled[any]("Foo", "Bar", "Beer"),
			Scheduled[any](true, true, false),
		), [][]any{
			{27, "Foo", true},
			{25, "Bar", true},
			{29, "Beer", false},
		}, nil, true)
	})
}
