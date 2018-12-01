package rxgo

import (
	"github.com/reactivex/rxgo/handlers"
	"github.com/stvp/assert"
	"testing"
)

// Assertion lists the assertions which may be configured.
type Assertion interface {
	apply(*assertion)
	HasItems() (bool, []interface{})
	HasSize() (bool, int)
}

type assertion struct {
	f             func(*assertion)
	checkHasItems bool
	hasItems      []interface{}
	checkHasSize  bool
	hasSize       int
}

func (ass *assertion) HasItems() (bool, []interface{}) {
	return ass.checkHasItems, ass.hasItems
}

func (ass *assertion) HasSize() (bool, int) {
	return ass.checkHasSize, ass.hasSize
}

func (ass *assertion) apply(do *assertion) {
	ass.f(do)
}

func newAssertion(f func(*assertion)) *assertion {
	return &assertion{
		f: f,
	}
}

func parseAssertions(assertions ...Assertion) Assertion {
	a := new(assertion)
	for _, assertion := range assertions {
		assertion.apply(a)
	}
	return a
}

// HasItems checks that an observable produces the corresponding items.
func HasItems(items ...interface{}) Assertion {
	return newAssertion(func(a *assertion) {
		a.checkHasItems = true
		a.hasItems = items
	})
}

// HasItems checks that an observable produces the corresponding number of items.
func HasSize(size int) Assertion {
	return newAssertion(func(a *assertion) {
		a.checkHasSize = true
		a.hasSize = size
	})
}

// AssertThatObservable asserts the result of an observable against a list of assertions.
func AssertThatObservable(t *testing.T, observable Observable, assertion ...Assertion) {
	ass := parseAssertions(assertion...)
	got := make([]interface{}, 0)
	observable.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = append(got, i)
	})).Block()

	checkHasItems, items := ass.HasItems()
	if checkHasItems {
		assert.Equal(t, items, got)
	}

	checkHasSize, size := ass.HasSize()
	if checkHasSize {
		assert.Equal(t, size, len(got))
	}
}
