package optional

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOf(t *testing.T) {
	some1 := Of("foo")
	assert.False(t, some1.IsEmpty())
	assert.Exactly(t, some1.Get(), "foo")

	some2 := Of(nil)
	assert.False(t, some2.IsEmpty())
	assert.Nil(t, some2.Get())

	empty := Empty()
	assert.True(t, empty.IsEmpty())
}
