package rxgo

import (
	"testing"
	"time"
)

func TestSwitchMap(t *testing.T) {
	// checkObservableResults(t, Pipe1(
	// 	Scheduled[uint](1, 2),
	// 	SwitchMap(func(x uint, i uint) Observable[string] {
	// 		return Pipe2(
	// 			Interval(time.Second),
	// 			Map(func(y, _ uint) (string, error) {
	// 				return fmt.Sprintf("x -> %d, y -> %d", x, y), nil
	// 			}),
	// 			Take[string](3),
	// 		)
	// 	}),
	// ), []string{"x -> 2, y -> 0", "x -> 2, y -> 1",
	// 	"x -> 2, y -> 2"}, nil, true)
}

func TestExhaustMap(t *testing.T) {

}

func TestMerge(t *testing.T) {
	// err := fmt.Errorf("some error")
	// checkObservableResults(t, Pipe1(
	// 	Scheduled[any](1, err),
	// 	Merge(Scheduled[any](1)),
	// ), []any{1, 1}, err, false)

	t.Run("Merge with Interval", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Pipe1(Interval(time.Millisecond), Take[uint](3)),
			Merge(Scheduled[uint](1)),
		), []uint{1, 0, 1, 2}, nil, true)
	})
}
