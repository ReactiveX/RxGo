package rxgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RxAssert lists the Observable assertions.
type RxAssert interface {
	apply(*rxAssert)
	itemsToBeChecked() (bool, []interface{})
	noItemsToBeChecked() bool
	raisedErrorToBeChecked() (bool, error)
	notRaisedErrorToBeChecked() bool
}

type rxAssert struct {
	f                      func(*rxAssert)
	checkHasItems          bool
	checkHasNoItems        bool
	items                  []interface{}
	checkHasRaisedError    bool
	error                  error
	checkHasNotRaisedError bool
}

func (ass *rxAssert) apply(do *rxAssert) {
	ass.f(do)
}

func (ass *rxAssert) itemsToBeChecked() (bool, []interface{}) {
	return ass.checkHasItems, ass.items
}

func (ass *rxAssert) noItemsToBeChecked() bool {
	return ass.checkHasNoItems
}

func (ass *rxAssert) raisedErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedError, ass.error
}

func (ass *rxAssert) notRaisedErrorToBeChecked() bool {
	return ass.checkHasNotRaisedError
}

func newAssertion(f func(*rxAssert)) *rxAssert {
	return &rxAssert{
		f: f,
	}
}

// HasItems checks that the observable produces the corresponding items.
func HasItems(items ...interface{}) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasItems = true
		a.items = items
	})
}

// HasNoItems checks that the observable has not produce any item.
func HasNoItems() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasNoItems = true
	})
}

// HasRaisedError checks that the observable has produce a specific error.
func HasRaisedError(err error) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedError = true
		a.error = err
	})
}

// HasNotRaisedError checks that the observable has not raised any error.
func HasNotRaisedError() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedError = true
	})
}

func parseAssertions(assertions ...RxAssert) RxAssert {
	ass := new(rxAssert)
	for _, assertion := range assertions {
		assertion.apply(ass)
	}
	return ass
}

// AssertObservable asserts the result of an Observable against a list of assertions.
func AssertObservable(ctx context.Context, t *testing.T, observable Observable, assertions ...RxAssert) {
	ass := parseAssertions(assertions...)

	got := make([]interface{}, 0)
	var err error
	done := make(chan struct{})

	observable.ForEach(ctx, func(i interface{}) {
		got = append(got, i)
	}, func(e error) {
		err = e
		close(done)
	}, func() {
		close(done)
	})

	if checkHasItems, expectedItems := ass.itemsToBeChecked(); checkHasItems {
		<-done
		assert.Equal(t, expectedItems, got)
	}

	if checkHasNoItems := ass.noItemsToBeChecked(); checkHasNoItems {
		<-done
		assert.Equal(t, 0, len(got))
	}

	if checkHasRaisedError, expectedError := ass.raisedErrorToBeChecked(); checkHasRaisedError {
		assert.Equal(t, expectedError, err)
	}

	if ass.notRaisedErrorToBeChecked() {
		assert.Nil(t, err)
	}
}
