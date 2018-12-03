package rxgo

import (
	"testing"

	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
)

// ObservableAssertion lists the assertions which may be configured on an Observable.
type ObservableAssertion interface {
	apply(*assertion)
	HasItems() (bool, []interface{})
	HasSize() (bool, int)
}

// Single lists the assertions which may be configured on a Single.
type SingleAssertion interface {
	apply(*assertion)
	HasValue() (bool, interface{})
	HasRaisedError() (bool, error)
	HasRaisedAnError() bool
}

type assertion struct {
	f                     func(*assertion)
	checkHasItems         bool
	hasItems              []interface{}
	checkHasSize          bool
	hasSize               int
	checkHasValue         bool
	hasValue              interface{}
	checkHasRaisedError   bool
	hasRaisedError        error
	checkHasRaisedAnError bool
}

func (ass *assertion) HasItems() (bool, []interface{}) {
	return ass.checkHasItems, ass.hasItems
}

func (ass *assertion) HasSize() (bool, int) {
	return ass.checkHasSize, ass.hasSize
}

func (ass *assertion) HasValue() (bool, interface{}) {
	return ass.checkHasValue, ass.hasValue
}

func (ass *assertion) HasRaisedError() (bool, error) {
	return ass.checkHasRaisedError, ass.hasRaisedError
}

func (ass *assertion) HasRaisedAnError() bool {
	return ass.checkHasRaisedAnError
}

func (ass *assertion) apply(do *assertion) {
	ass.f(do)
}

func newAssertion(f func(*assertion)) *assertion {
	return &assertion{
		f: f,
	}
}

func parseObservableAssertions(assertions ...ObservableAssertion) ObservableAssertion {
	a := new(assertion)
	for _, assertion := range assertions {
		assertion.apply(a)
	}
	return a
}

func parseSingleAssertions(assertions ...SingleAssertion) SingleAssertion {
	a := new(assertion)
	for _, assertion := range assertions {
		assertion.apply(a)
	}
	return a
}

// HasItems checks that an observable produces the corresponding items.
func HasItems(items ...interface{}) ObservableAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasItems = true
		a.hasItems = items
	})
}

// HasItems checks that an observable produces the corresponding number of items.
func HasSize(size int) ObservableAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasSize = true
		a.hasSize = size
	})
}

// HasValue checks that a single produces the corresponding value.
func HasValue(value interface{}) SingleAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasValue = true
		a.hasValue = value
	})
}

// HasRaisedError checks that a single raises the corresponding error.
func HasRaisedError(err error) SingleAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasRaisedError = true
		a.hasRaisedError = err
	})
}

// HasRaisedAnError checks that a single raises an error.
func HasRaisedAnError() SingleAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasRaisedAnError = true
	})
}

// AssertThatObservable asserts the result of an observable against a list of assertions.
func AssertThatObservable(t *testing.T, observable Observable, assertions ...ObservableAssertion) {
	ass := parseObservableAssertions(assertions...)
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

func AssertThatSingle(t *testing.T, single Single, assertions ...SingleAssertion) {
	ass := parseSingleAssertions(assertions...)

	v, err := single.Subscribe(nil).Block()

	checkHasValue, value := ass.HasValue()
	if checkHasValue {
		if err != nil {
			assert.Fail(t, "error is not nil")
		} else {
			assert.Equal(t, value, v)
		}
	}

	checkHasRaisedAnError := ass.HasRaisedAnError()
	if checkHasRaisedAnError {
		assert.NotNil(t, err)
	}

	checkHasRaisedError, value := ass.HasRaisedError()
	if checkHasRaisedError {
		assert.Equal(t, value, err)
	}
}
