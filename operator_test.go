package rxgo

import (
	"testing"
)

func TestTake(t *testing.T) {

}

func TestTakeUntil(t *testing.T) {

}

func TestTakeWhile(t *testing.T) {

}

func TestTakeLast(t *testing.T) {

}

func TestElementAt(t *testing.T) {

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
