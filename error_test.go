package rxgo

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCatch(t *testing.T) {
	t.Run("Catch with Empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Empty[string](),
			Catch(func(err error, caught Observable[string]) Observable[string] {
				return Of2("I", "II", "III", "IV", "V")
			}),
		), []string{}, nil, true)
	})

	t.Run("Catch with values", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2("A", "I", "II", "III", "IV", "V", "Z"),
			Catch(func(err error, caught Observable[string]) Observable[string] {
				return caught
			}),
		), []string{"A", "I", "II", "III", "IV", "V", "Z"}, nil, true)
	})

	t.Run("Catch with Throw", func(t *testing.T) {
		var err = fmt.Errorf("throw")
		checkObservableResults(t, Pipe1(
			Throw[string](func() error {
				return err
			}),
			Catch(func(err error, caught Observable[string]) Observable[string] {
				return Of2("I", "II", "III", "IV", "V")
			}),
		), []string{"I", "II", "III", "IV", "V"}, nil, true)
	})

	t.Run("Catch with Map error", func(t *testing.T) {
		var err = fmt.Errorf("throw four")
		checkObservableResults(t, Pipe2(
			Of2[any](1, 2, 3, 4, 5),
			Map(func(v any, _ uint) (any, error) {
				if v == 4 {
					return 0, err
				}
				return v, nil
			}),
			Catch(func(err error, caught Observable[any]) Observable[any] {
				return Of2[any]("I", "II", "III", "IV", "V")
			}),
		), []any{1, 2, 3, "I", "II", "III", "IV", "V"}, nil, true)
	})

	t.Run("Catch with same observable", func(t *testing.T) {
		var err = fmt.Errorf("throw four")

		checkObservableResults(t, Pipe2(
			Of2[any](1, 2, 3, 4, 5),
			Map(func(v any, _ uint) (any, error) {
				if v == 4 {
					return 0, err
				}
				return v, nil
			}),
			Catch(func(err error, caught Observable[any]) Observable[any] {
				return Empty[any]()
			}),
		), []any{1, 2, 3}, nil, true)
	})
}

func TestRetry(t *testing.T) {
	t.Run("Retry with Empty", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Empty[any](),
			Retry[any, uint8](2),
		), []any{}, nil, true)
	})

	t.Run("Retry with Throw", func(t *testing.T) {
		var err = fmt.Errorf("throwing")
		checkObservableResults(t, Pipe1(
			Throw[string](func() error {
				return err
			}),
			Retry[string, uint8](2),
		), []string{}, err, false)
	})

	t.Run("Retry with count 2", func(t *testing.T) {
		var err = fmt.Errorf("throw five")
		checkObservableResults(t, Pipe2(
			Interval(time.Millisecond*100),
			Map(func(v, _ uint) (uint, error) {
				if v > 5 {
					return 0, err
				}
				return v, nil
			}),
			Retry[uint, uint8](2),
		), []uint{0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5}, err, false)
	})

	t.Run("Retry with count 2", func(t *testing.T) {
		var err = fmt.Errorf("throw five")
		checkObservableResults(t, Pipe2(
			Interval(time.Millisecond*100),
			Map(func(v, _ uint) (uint, error) {
				if v > 5 {
					return 0, err
				}
				return v, nil
			}),
			Retry[uint](RetryConfig{
				Count: 2,
				Delay: time.Millisecond,
			}),
		), []uint{0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5}, err, false)
	})

	t.Run("Retry with Iif", func(t *testing.T) {
		var ok = false
		checkObservableResults(t, Pipe1(
			Iif(func() bool {
				if ok {
					ok = false
					return ok
				}
				ok = true
				return ok
			},
				Of2("x", "^", "@", "#"),
				Throw[string](func() error {
					return errors.New("retry")
				})),
			Retry[string, uint](3),
		), []string{"x", "^", "@", "#"}, nil, true)
	})

	t.Run("Retry with Defer", func(t *testing.T) {
		var count = 0
		checkObservableResults(t, Pipe1(
			Defer(func() Observable[string] {
				count++
				if count < 2 {
					return Throw[string](func() error {
						return errors.New("retry")
					})
				}
				return Of2("hello", "world", "!!!")
			}),
			Retry[string, uint](3),
		), []string{"hello", "world", "!!!"}, nil, true)
	})

	t.Run("Retry with Defer", func(t *testing.T) {
		var count = 0
		checkObservableResults(t, Pipe1(
			Defer(func() Observable[string] {
				count++
				if count <= 3 {
					return Throw[string](func() error {
						return errors.New("retry")
					})
				}
				return Of2("hello", "world", "!!!")
			}),
			Retry[string, uint](3),
		), []string{"hello", "world", "!!!"}, nil, true)
	})

	t.Run("Retry with Defer but failed", func(t *testing.T) {
		var (
			count = 0
			err   = errors.New("retry")
		)
		checkObservableResults(t, Pipe1(
			Defer(func() Observable[string] {
				count++
				if count < 5 {
					return Throw[string](func() error {
						return err
					})
				}
				return Of2("hello", "world", "!!!")
			}),
			Retry[string, uint](2),
		), []string{}, err, false)
	})
}
