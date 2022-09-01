package rxgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
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

func TestEmit(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Emit data", func(t *testing.T) {
		str := "hello world"
		ch := make(chan DataValuer[string], 1)
		emitData(str, ch)
		v := <-ch
		require.Equal(t, str, v.Value())
		require.Nil(t, v.Err())
	})

	t.Run("Emit error", func(t *testing.T) {
		err := fmt.Errorf("some error")
		ch := make(chan DataValuer[string], 1)
		emitError(err, ch)
		v := <-ch
		require.Equal(t, err, v.Err())
		require.Empty(t, v.Value())
	})
}
