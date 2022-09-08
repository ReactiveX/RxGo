package rxgo

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSwitchMap(t *testing.T) {
	// checkObservableResults(t, Pipe1(
	// 	Scheduled[uint](1, 2),
	// 	SwitchMap(func(x uint, i uint) IObservable[string] {
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

func TestConcatMap(t *testing.T) {
	t.Run("ConcatMap with error on upstream", func(t *testing.T) {
		var err = fmt.Errorf("throw")
		checkObservableResults(t, Pipe1(
			Scheduled[any]("z", err, "q"),
			ConcatMap(func(x any, i uint) IObservable[string] {
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

	t.Run("ConcatMap with ThrownError on conditional return stream", func(t *testing.T) {
		var err = fmt.Errorf("throw")

		mapTo := func(v string, i uint) string {
			return fmt.Sprintf("%s[%d]", v, i)
		}

		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) IObservable[string] {
				if i == 0 {
					return Scheduled(mapTo(x, i), mapTo(x, i), mapTo(x, i))
				}

				return ThrownError[string](func() error {
					return err
				})
			}),
		), []string{"z[0]", "z[0]", "z[0]"}, err, false)
	})

	t.Run("ConcatMap with ThrownError on return stream", func(t *testing.T) {
		var err = fmt.Errorf("throw")

		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) IObservable[string] {
				return ThrownError[string](func() error {
					return err
				})
			}),
		), []string{}, err, false)
	})

	t.Run("ConcatMap with Interval + Map which return error", func(t *testing.T) {
		var err = errors.New("nopass")
		checkObservableResults(t, Pipe1(
			Scheduled("z", "q"),
			ConcatMap(func(x string, i uint) IObservable[string] {
				return Pipe2(
					Interval(time.Second),
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
			ConcatMap(func(x uint, i uint) IObservable[string] {
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

func TestMergeMap(t *testing.T) {

	t.Run("MergeMap with complete", func(t *testing.T) {
		var (
			result = make([]Tuple[string, uint], 0)
			err    error
			done   bool
		)
		Pipe1(
			Scheduled("a", "b", "v"),
			MergeMap(func(x string, i uint) IObservable[Tuple[string, uint]] {
				return Pipe2(
					Interval(time.Millisecond),
					Map(func(y, _ uint) (Tuple[string, uint], error) {
						return NewTuple(x, y), nil
					}),
					Take[Tuple[string, uint]](5),
				)
			}),
		).SubscribeSync(func(s Tuple[string, uint]) {
			result = append(result, s)
		}, func(e error) {
			err = e
		}, func() {
			done = true
		})
		require.True(t, len(result) == 15)
		require.Nil(t, err)
		require.True(t, done)
	})

	// t.Run("MergeMap with error", func(t *testing.T) {
	// 	var (
	// 		result = make([]Tuple[string, uint], 0)
	// 		failed = errors.New("failed")
	// 		err    error
	// 		done   bool
	// 	)
	// 	Pipe1(
	// 		Scheduled("a", "b", "v"),
	// 		MergeMap(func(x string, i uint) IObservable[Tuple[string, uint]] {
	// 			return Pipe2(
	// 				Interval(time.Millisecond),
	// 				Map(func(y, idx uint) (Tuple[string, uint], error) {
	// 					if idx > 3 {
	// 						return nil, failed
	// 					}
	// 					return NewTuple(x, y), nil
	// 				}),
	// 				Take[Tuple[string, uint]](5),
	// 			)
	// 		}),
	// 	).SubscribeSync(func(s Tuple[string, uint]) {
	// 		result = append(result, s)
	// 	}, func(e error) {
	// 		err = e
	// 	}, func() {
	// 		done = true
	// 	})
	// 	require.True(t, len(result) == 9)
	// 	require.Equal(t, failed, err)
	// 	require.False(t, done)
	// })
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
