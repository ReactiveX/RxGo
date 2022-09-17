package rxgo

import (
	"testing"
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
