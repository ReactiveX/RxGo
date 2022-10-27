package rxgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithTimestamp(t *testing.T) {
	t.Run("WithTimestamp with Numbers", func(t *testing.T) {
		var (
			result = make([]Timestamp[uint], 0)
			done   = true
		)
		Pipe1(Range[uint](1, 5), WithTimestamp[uint]()).
			SubscribeSync(func(t Timestamp[uint]) {
				result = append(result, t)
			}, func(err error) {}, func() {
				done = true
			})

		for i, v := range result {
			require.False(t, v.Time().IsZero())
			require.Equal(t, uint(i+1), v.Value())
		}
		require.True(t, done)
	})
}
