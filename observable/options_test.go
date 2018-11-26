package observable

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
