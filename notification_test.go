package rxgo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDematerialize(t *testing.T) {
	t.Run("Dematerialize with complete", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2[ObservableNotification[any]](
				Next[any]("a"),
				Next[any]("a"),
				Complete[any](),
			),
			Dematerialize[any](),
		), []any{"a", "a"}, nil, true)
	})

	t.Run("Dematerialize with error", func(t *testing.T) {
		var err = errors.New("failed on Of")
		checkObservableResults(t, Pipe1(
			Of2[ObservableNotification[any]](
				Next[any]("a"),
				Next[any]("joker"),
				Error[any](err),
				Next[any]("q"),
			),
			Dematerialize[any](),
		), []any{"a", "joker"}, err, false)
	})

	t.Run("Dematerialize with numbers (complete)", func(t *testing.T) {
		checkObservableResults(t, Pipe1(
			Of2[ObservableNotification[int64]](
				Next[int64](1088),
				Next[int64](18394004),
				Next[int64](-28292912),
				Next[int64](1),
				Next[int64](-1626),
				Complete[int64](),
			),
			Dematerialize[int64](),
		), []int64{1088, 18394004, -28292912, 1, -1626}, nil, true)
	})
}

func TestMaterialize(t *testing.T) {
	t.Run("Materialize with error", func(t *testing.T) {
		var err = fmt.Errorf("throw")
		checkObservableResults(t, Pipe1(Scheduled[any](1, "a", err, 100), Materialize[any]()),
			[]ObservableNotification[any]{
				Next[any](1),
				Next[any]("a"),
				Error[any](err),
			}, nil, true)
	})

	t.Run("Materialize with complete", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Scheduled[any](1, "a", struct{}{}), Materialize[any]()),
			[]ObservableNotification[any]{
				Next[any](1),
				Next[any]("a"),
				Next[any](struct{}{}),
				Complete[any](),
			}, nil, true)
	})
}

func TestNotification(t *testing.T) {
	t.Run("Next Notification", func(t *testing.T) {
		value := "hello world"
		data := Next(value)
		require.Equal(t, value, data.Value())
		require.Nil(t, data.Err())
	})

	err := fmt.Errorf("uncaught error")
	t.Run("Error Notification with any", func(t *testing.T) {
		data := Error[any](err)
		require.Nil(t, data.Value())
		require.Equal(t, err, data.Err())
	})

	t.Run("Error Notification with string", func(t *testing.T) {
		data := Error[string](err)
		require.Equal(t, "", data.Value())
		require.Equal(t, err, data.Err())
	})
}
