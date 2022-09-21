package rxgo

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepeat(t *testing.T) {
	t.Run("Repeat with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[any](),
			Repeat[any, uint](3),
		), []any{}, nil, true)
	})

	t.Run("Repeat with config", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2[any](12, 88),
			Repeat[any](RepeatConfig{
				Count: 5,
				Delay: time.Millisecond,
			}),
		), []any{12, 88, 12, 88, 12, 88, 12, 88, 12, 88}, nil, true)
	})

	t.Run("Repeat with error", func(t *testing.T) {
		var err = errors.New("throw")
		// Repeat with error will no repeat
		checkObservableResults(t, Pipe1(
			ThrowError[string](func() error {
				return err
			}),
			Repeat[string, uint](3),
		), []string{}, err, false)
	})

	t.Run("Repeat with error", func(t *testing.T) {
		var err = errors.New("throw at 3")
		// Repeat with error will no repeat
		checkObservableResults(t, Pipe2(
			Range[uint](1, 10),
			Map(func(v, _ uint) (uint, error) {
				if v >= 3 {
					return 0, err
				}
				return v, nil
			}),
			Repeat[uint, uint](3),
		), []uint{1, 2}, err, false)
	})

	t.Run("Repeat with alphaberts", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2("a", "bb", "cc", "ff", "gg"),
			Repeat[string, uint](3),
		), []string{
			"a", "bb", "cc", "ff", "gg",
			"a", "bb", "cc", "ff", "gg",
			"a", "bb", "cc", "ff", "gg",
		}, nil, true)
	})
}

func TestDo(t *testing.T) {
	t.Run("Do with Range(1, 5)", func(t *testing.T) {
		result := make([]string, 0)
		checkObservableResults(t, Pipe1(
			Range[uint](1, 5),
			Do(NewObserver(func(v uint) {
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

	t.Run("Do with Error", func(t *testing.T) {
		var (
			err    = fmt.Errorf("An error")
			result = make([]string, 0)
		)
		checkObservableResults(t, Pipe1(
			Scheduled[any](1, err),
			Do(NewObserver(func(v any) {
				result = append(result, fmt.Sprintf("Number(%v)", v))
			}, nil, nil)),
		), []any{1}, err, false)
		require.ElementsMatch(t, []string{"Number(1)"}, result)
	})
}

func TestDelay(t *testing.T) {

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

func TestTimeout(t *testing.T) {
	t.Run("Timeout with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			EMPTY[any](),
			Timeout[any](time.Second),
		), nil, nil, true)
	})

	t.Run("Timeout with error", func(t *testing.T) {
		var err = errors.New("failed")
		checkObservableResult(t, Pipe1(
			ThrowError[any](func() error {
				return err
			}),
			Timeout[any](time.Millisecond),
		), nil, err, false)
	})

	t.Run("Timeout with timeout error", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Interval(time.Millisecond*5),
			Timeout[uint](time.Millisecond),
		), uint(0), ErrTimeout, false)
	})

	t.Run("Timeout with Scheduled", func(t *testing.T) {
		checkObservableResult(t, Pipe1(
			Pipe1(Scheduled("a"), Delay[string](time.Millisecond*50)),
			Timeout[string](time.Millisecond*100),
		), "a", nil, true)
	})
}

func TestToArray(t *testing.T) {
	t.Run("ToArray with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), ToArray[any]()), []any{}, nil, true)
	})

	t.Run("ToArray with error", func(t *testing.T) {
		var err = errors.New("throw")
		checkObservableResult(t, Pipe1(Scheduled[any]("a", "z", err), ToArray[any]()),
			nil, err, false)
	})

	t.Run("ToArray with numbers", func(t *testing.T) {
		checkObservableResult(t, Pipe1(Range[uint](1, 5), ToArray[uint]()), []uint{1, 2, 3, 4, 5}, nil, true)
	})

	t.Run("ToArray with alphaberts", func(t *testing.T) {
		checkObservableResult(t, Pipe1(newObservable(func(subscriber Subscriber[string]) {
			for i := 1; i <= 5; i++ {
				subscriber.Send() <- Next(string(rune('A' - 1 + i)))
			}
			subscriber.Send() <- Complete[string]()
		}), ToArray[string]()), []string{"A", "B", "C", "D", "E"}, nil, true)
	})
}
