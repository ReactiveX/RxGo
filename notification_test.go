package rxgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaterialize(t *testing.T) {
	t.Run("Materialize with error", func(t *testing.T) {
		var err = fmt.Errorf("throw")
		checkObservableResults(t, Pipe1(Scheduled[any](1, "a", err), Materialize[any]()),
			[]ObservableNotification[any]{
				NextNotification[any](1),
				NextNotification[any]("a"),
				ErrorNotification[any](err),
			}, nil, true)
	})

	t.Run("Materialize with complete", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Scheduled[any](1, "a", struct{}{}), Materialize[any]()),
			[]ObservableNotification[any]{
				NextNotification[any](1),
				NextNotification[any]("a"),
				NextNotification[any](struct{}{}),
				CompleteNotification[any](),
			}, nil, true)
	})
}

func TestNotification(t *testing.T) {
	t.Run("Next Notification", func(t *testing.T) {
		value := "hello world"
		data := NextNotification(value)
		require.Equal(t, value, data.Value())
		require.Nil(t, data.Err())
	})

	err := fmt.Errorf("uncaught error")
	t.Run("Error Notification with any", func(t *testing.T) {
		data := ErrorNotification[any](err)
		require.Nil(t, data.Value())
		require.Equal(t, err, data.Err())
	})

	t.Run("Error Notification with string", func(t *testing.T) {
		data := ErrorNotification[string](err)
		require.Equal(t, "", data.Value())
		require.Equal(t, err, data.Err())
	})
}
