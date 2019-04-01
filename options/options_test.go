package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	option := ParseOptions(WithParallelism(1), WithBufferedChannel(2))

	assert.Equal(t, option.Parallelism(), 1)
	assert.Equal(t, option.Buffer(), 2)
}
