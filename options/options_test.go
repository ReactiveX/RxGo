package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions(t *testing.T) {
	option := ParseOptions(WithParallelism(1), WithBufferedChannel(2))

	assert.Equal(t, option.Parallelism(), 1)
	assert.Equal(t, option.BufferedChannelCapacity(), 2)
}
