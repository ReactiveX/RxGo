package rxgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	t.Run("Data with value", func(t *testing.T) {
		value := "hello world"
		data := streamData[string]{v: value}
		require.Equal(t, value, data.Value())
		require.Nil(t, data.Err())
	})

	err := fmt.Errorf("uncaught error")
	t.Run("Data with error[any]", func(t *testing.T) {
		data := streamData[any]{err: err}
		require.Nil(t, data.Value())
		require.Equal(t, err, data.Err())
	})

	t.Run("Data with error[string]", func(t *testing.T) {
		data := streamData[string]{err: err}
		require.Equal(t, "", data.Value())
		require.Equal(t, err, data.Err())
	})
}
