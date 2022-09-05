package rxgo

import "testing"

func TestSkip(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 10), Skip[uint](5)), []uint{6, 7, 8, 9, 10}, nil, true)
}

func TestSkipLast(t *testing.T) {
	checkObservableResults(t, Pipe1(Range[uint](1, 10), SkipLast[uint](5)), []uint{1, 2, 3, 4, 5}, nil, true)
}

// func TestSkipUntil(t *testing.T) {
// 	// checkObservableResults(t, Pipe1(Range[uint](1, 10), SkipUntil[uint](Scheduled[uint](2, 2, 3))), []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true)
// }

func TestSkipWhile(t *testing.T) {
	t.Run("SkipWhile until condition meet", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Scheduled("Green Arrow", "SuperMan", "Flash", "SuperGirl", "Black Canary"),
			SkipWhile(func(v string, _ uint) bool {
				return v != "SuperGirl"
			})), []string{"SuperGirl", "Black Canary"}, nil, true)
	})

	t.Run("SkipWhile until index 5", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Range[uint](1, 10), SkipWhile(func(_ uint, idx uint) bool {
			return idx != 5
		})), []uint{6, 7, 8, 9, 10}, nil, true)
	})
}
