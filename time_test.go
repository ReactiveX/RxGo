package rxgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithTimestamp(t *testing.T) {
	t.Run("WithTimestamp with Numbers", func(t *testing.T) {
		var (
			now    = time.Now()
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
			require.True(t, v.Time().After(now))
			require.Equal(t, uint(i+1), v.Value())
		}
		require.True(t, done)
	})
}
