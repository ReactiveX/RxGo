package rxgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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
