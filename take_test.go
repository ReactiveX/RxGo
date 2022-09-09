package rxgo

import (
	"testing"
	"time"
)

func TestTake(t *testing.T) {
	checkObservableResults(t, Pipe1(Interval(time.Millisecond), Take[uint](3)), []uint{0, 1, 2}, nil, true)
	checkObservableResults(t, Pipe1(Range[uint](1, 100), Take[uint](3)), []uint{1, 2, 3}, nil, true)
}

func TestTakeLast(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 100), TakeLast[uint](3)), []uint{98, 99, 100}, nil, true)
}

func TestTakeUntil(t *testing.T) {
	t.Run("TakeUntil with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[uint](),
			TakeUntil[uint](Scheduled("a")),
		), []uint{}, nil, true)
	})

	t.Run("TakeUntil with Interval", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Interval(time.Millisecond),
			TakeUntil[uint](Interval(time.Millisecond*5)),
		), []uint{0, 1, 2, 3}, nil, true)
	})
}

func TestTakeWhile(t *testing.T) {
	t.Run("TakeWhile with Interval", func(t *testing.T) {
		result := make([]uint, 0)
		for i := uint(0); i <= 5; i++ {
			result = append(result, i)
		}
		checkObservableResults(t, Pipe1(Interval(time.Millisecond), TakeWhile(func(v uint, _ uint) bool {
			return v <= 5
		})), result, nil, true)
	})

	t.Run("TakeWhile with Range", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 100), TakeWhile(func(v uint, _ uint) bool {
			return v >= 50
		})), []uint{}, nil, true)
	})
}
