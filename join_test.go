package rxgo

import (
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
			ThrownError[string](func() error {
				return createErr(1)
			}),
			ThrownError[string](func() error {
				return createErr(2)
			}),
			ThrownError[string](func() error {
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

func TestZip(t *testing.T) {
	t.Run("Zip with Scheduled", func(t *testing.T) {
		checkObservableResults(t, Zip(
			EMPTY[any](),
			Scheduled[any]("Foo", "Bar", "Beer"),
			Scheduled[any](true, true, false),
		), [][]any{
			{nil, "Foo", true},
			{nil, "Bar", true},
			{nil, "Beer", false},
			{nil, nil, nil},
		}, nil, true)
	})

	t.Run("Zip with Scheduled", func(t *testing.T) {
		checkObservableResults(t, Zip(
			Scheduled[any](27, 25, 29),
			Scheduled[any]("Foo", "Bar", "Beer"),
			Scheduled[any](true, true, false),
		), [][]any{
			{27, "Foo", true},
			{25, "Bar", true},
			{29, "Beer", false},
			{nil, nil, nil},
		}, nil, true)
	})
}
