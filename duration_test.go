package rxgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithFrequency(t *testing.T) {
	frequency := WithDuration(100 * time.Millisecond)
	assert.Equal(t, 100*time.Millisecond, frequency.duration())
}
