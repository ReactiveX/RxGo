package rxgo

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	// Of something test
	some1 := Of("foo")
	got, err := some1.Get()

	assert.False(t, some1.IsEmpty())
	assert.Nil(t, err)
	assert.Exactly(t, got, "foo")
}

func TestOfEmpty(t *testing.T) {
	some := Of(nil)
	got, err := some.Get()

	assert.False(t, some.IsEmpty())
	assert.Nil(t, err)
	assert.Exactly(t, got, nil)
}

func TestEmpty(t *testing.T) {
	empty := EmptyOptional()
	got, err := empty.Get()
	assert.True(t, empty.IsEmpty())
	assert.Error(t, err)
	assert.IsType(t, &NoSuchElementError{}, errors.Cause(err))
	assert.Exactly(t, got, nil)
}
