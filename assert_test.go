package rxgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasItems(t *testing.T) {
	ass := parseAssertions(HasItems(1, 2, 3))

	configured, items := ass.HasItems()
	assert.True(t, configured)
	assert.Equal(t, []interface{}{1, 2, 3}, items)

	configured, _ = ass.HasSize()
	assert.False(t, configured)
}

func TestHasSize(t *testing.T) {
	ass := parseAssertions(HasSize(3))

	configured, size := ass.HasSize()
	assert.True(t, configured)
	assert.Equal(t, 3, size)

	configured, _ = ass.HasItems()
	assert.False(t, configured)
}

func TestAssertHasItems(t *testing.T) {
	AssertThatObservable(t, Just(1, 2, 3), HasItems(1, 2, 3))
}

func TestAssertHasSize(t *testing.T) {
	AssertThatObservable(t, Just(1, 2, 3), HasSize(3))
}
