package optional

import (
	"testing"

	"github.com/reactivex/rxgo/errors"
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
	empty := Empty()
	got, err := empty.Get()
	assert.True(t, empty.IsEmpty())
	if err != nil {
		assert.Exactly(t, errors.New(errors.NoSuchElementError), err)
	} else {
		assert.Fail(t, "error is not nil")
	}
	assert.Exactly(t, got, nil)
}
