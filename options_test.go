package rxgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithParallelism(t *testing.T) {
	var observableOptions options

	option := WithParallelism(2)
	option.apply(&observableOptions)

	assert.Equal(t, observableOptions.parallelism, 2)
}

func TestWithBufferedChannel(t *testing.T) {
	var observableOptions options

	option := WithBufferedChannel(2)
	option.apply(&observableOptions)

	assert.Equal(t, observableOptions.channelBufferCapacity, 2)
}
