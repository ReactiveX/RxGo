package rxgo

import (
	"fmt"
	"testing"
)

func TestCatchError(t *testing.T) {
	t.Run("CatchError with EMPTY", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			EMPTY[string](),
			CatchError(func(err error, caught Observable[string]) Observable[string] {
				return Of2("I", "II", "III", "IV", "V")
			}),
		), []string{}, nil, true)
	})

	t.Run("CatchError with values", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2("A", "I", "II", "III", "IV", "V", "Z"),
			CatchError(func(err error, caught Observable[string]) Observable[string] {
				return caught
			}),
		), []string{"A", "I", "II", "III", "IV", "V", "Z"}, nil, true)
	})

	t.Run("CatchError with ThrowError", func(t *testing.T) {
		var err = fmt.Errorf("throw")
		checkObservableResults(t, Pipe1(
			ThrowError[string](func() error {
				return err
			}),
			CatchError(func(err error, caught Observable[string]) Observable[string] {
				return Of2("I", "II", "III", "IV", "V")
			}),
		), []string{"I", "II", "III", "IV", "V"}, nil, true)
	})

	t.Run("CatchError with Map error", func(t *testing.T) {
		var err = fmt.Errorf("throw four")
		checkObservableResults(t, Pipe2(
			Of2[any](1, 2, 3, 4, 5),
			Map(func(v any, _ uint) (any, error) {
				if v == 4 {
					return 0, err
				}
				return v, nil
			}),
			CatchError(func(err error, caught Observable[any]) Observable[any] {
				return Of2[any]("I", "II", "III", "IV", "V")
			}),
		), []any{1, 2, 3, "I", "II", "III", "IV", "V"}, nil, true)
	})

	t.Run("CatchError with same observable", func(t *testing.T) {
		var err = fmt.Errorf("throw four")
		checkObservableResults(t, Pipe2(
			Of2[any](1, 2, 3, 4, 5),
			Map(func(v any, _ uint) (any, error) {
				if v == 4 {
					return 0, err
				}
				return v, nil
			}),
			CatchError(func(err error, caught Observable[any]) Observable[any] {
				return caught
			}),
		), []any{1, 2, 3, 1, 2, 3}, err, false)
	})
}
