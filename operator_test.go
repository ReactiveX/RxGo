package rxgo

import (
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

}

func TestDistinctUntilChanged(t *testing.T) {

}

func TestFilter(t *testing.T) {

}

func TestMap(t *testing.T) {

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

func TestToArray(t *testing.T) {

}
