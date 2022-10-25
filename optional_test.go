package rxgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptional(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		str := "hello world"
		opt := Some(str)
		require.Equal(t, str, opt.MustGet())
		require.Equal(t, str, opt.OrElse(""))
		require.False(t, opt.IsNone())
		v, ok := opt.Get()
		require.Equal(t, str, v)
		require.True(t, ok)
	})

	t.Run("None", func(t *testing.T) {
		opt := None[string]()
		require.Panics(t, func() {
			opt.MustGet()
		})
		require.True(t, opt.IsNone())
		require.Equal(t, "Joker", opt.OrElse("Joker"))
		v, ok := opt.Get()
		require.Empty(t, v)
		require.False(t, ok)
	})
}
