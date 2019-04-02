package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	option := ParseOptions(WithBufferedChannel(2))
	assert.Equal(t, option.Buffer(), 2)
}
