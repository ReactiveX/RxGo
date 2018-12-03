package rxgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasItems(t *testing.T) {
	ass := parseObservableAssertions(HasItems(1, 2, 3))

	configured, items := ass.hasItemsFunc()
	assert.True(t, configured)
	assert.Equal(t, []interface{}{1, 2, 3}, items)

	configured, _ = ass.hasSizeFunc()
	assert.False(t, configured)
}

func TestHasSize(t *testing.T) {
	ass := parseObservableAssertions(HasSize(3))

	configured, size := ass.hasSizeFunc()
	assert.True(t, configured)
	assert.Equal(t, 3, size)

	configured, _ = ass.hasItemsFunc()
	assert.False(t, configured)
}

func TestIsEmpty(t *testing.T) {
	ass := parseObservableAssertions(IsEmpty())

	configured, size := ass.hasSizeFunc()
	assert.True(t, configured)
	assert.Equal(t, 0, size)
}

func TestHasValue(t *testing.T) {
	ass := parseSingleAssertions(HasValue(1))

	configured, value := ass.hasValueFunc()
	assert.True(t, configured)
	assert.Equal(t, 1, value)
}

func TestHasRaisedError(t *testing.T) {
	ass := parseSingleAssertions(HasRaisedError(errors.New("foo")))

	configured, error := ass.hasRaisedErrorFunc()
	assert.True(t, configured)
	assert.Equal(t, errors.New("foo"), error)
}

func TestHasRaisedAnError(t *testing.T) {
	ass := parseSingleAssertions(HasRaisedAnError())

	configured := ass.hasRaisedAnErrorFunc()
	assert.True(t, configured)
}

func TestAssertHasItems(t *testing.T) {
	AssertThatObservable(t, Just(1, 2, 3), HasItems(1, 2, 3))
}

func TestAssertHasSize(t *testing.T) {
	AssertThatObservable(t, Just(1, 2, 3), HasSize(3))
}

func TestAssertHasValue(t *testing.T) {
	AssertThatSingle(t, newSingleFrom(1), HasValue(1))
}

func TestAssertSingleError(t *testing.T) {
	AssertThatSingle(t, newSingleFrom(errors.New("foo")),
		HasRaisedAnError(), HasRaisedError(errors.New("foo")))
}
