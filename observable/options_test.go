package observable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithParallelism(t *testing.T) {
	var o options

	WithParallelism(2).apply(&o)

	assert.Equal(t, o.parallelism, 2)
}

func TestParseOptions(t *testing.T) {
	var o options

	o.parseOptions(WithParallelism(2))

	assert.Equal(t, o.parallelism, 2)
}
