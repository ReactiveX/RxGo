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
	someItemsToBeChecked() bool
	raisedErrorToBeChecked() (bool, error)
	raisedAnErrorToBeChecked() (bool, error)
	notRaisedErrorToBeChecked() bool
	itemToBeChecked() (bool, interface{})
	noItemToBeChecked() (bool, interface{})
}

type rxAssert struct {
	f                      func(*rxAssert)
	checkHasItems          bool
	checkHasNoItems        bool
	checkHasSomeItems      bool
	items                  []interface{}
	checkHasRaisedError    bool
	error                  error
	checkHasRaisedAnError  bool
	checkHasNotRaisedError bool
	checkHasItem           bool
	item                   interface{}
	checkHasNoItem         bool
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

func (ass *rxAssert) someItemsToBeChecked() bool {
	return ass.checkHasSomeItems
}

func (ass *rxAssert) raisedErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedError, ass.error
}

func (ass *rxAssert) raisedAnErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedAnError, ass.error
}

func (ass *rxAssert) notRaisedErrorToBeChecked() bool {
	return ass.checkHasNotRaisedError
}

func (ass *rxAssert) itemToBeChecked() (bool, interface{}) {
	return ass.checkHasItem, ass.item
}

func (ass *rxAssert) noItemToBeChecked() (bool, interface{}) {
	return ass.checkHasNoItem, ass.item
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

// HasSomeItems checks that the observable produces some items.
func HasSomeItems() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasSomeItems = true
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

// HasRaisedAnError checks that the observable has produce an error.
func HasRaisedAnError() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedAnError = true
	})
}

// HasNotRaisedError checks that the observable has not raised any error.
func HasNotRaisedError() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedError = true
	})
}

// HasItem checks if a single or optional single has a specific item.
func HasItem(i interface{}) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasItem = true
		a.item = i
	})
}

// HasNoItem checks if a single or optional single has no item.
func HasNoItem() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasNoItem = true
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

	if ass.noItemsToBeChecked() {
		<-done
		assert.Equal(t, 0, len(got))
	}

	if ass.someItemsToBeChecked() {
		<-done
		assert.NotEqual(t, 0, len(got))
	}

	if checkHasRaisedError, expectedError := ass.raisedErrorToBeChecked(); checkHasRaisedError {
		assert.Equal(t, expectedError, err)
	}

	if checkHasRaisedAnError, expectedError := ass.raisedAnErrorToBeChecked(); checkHasRaisedAnError {
		assert.Nil(t, expectedError)
	}

	if ass.notRaisedErrorToBeChecked() {
		assert.Nil(t, err)
	}
}

// AssertSingle asserts the result of an single against a list of assertions.
func AssertSingle(ctx context.Context, t *testing.T, single Single, assertions ...RxAssert) {
	// TODO Reuse logic
	ass := parseAssertions(assertions...)

	var got interface{}
	var err error

	item := <-single.Observe()
	if item.IsError() {
		err = item.Err
	} else {
		got = item.Value
	}

	if checkHasValue, value := ass.itemToBeChecked(); checkHasValue {
		assert.Equal(t, value, got)
	}

	if checkHasNoValue, value := ass.noItemToBeChecked(); checkHasNoValue {
		assert.Equal(t, value, got)
	}

	if checkHasRaisedError, expectedError := ass.raisedErrorToBeChecked(); checkHasRaisedError {
		assert.Equal(t, expectedError, err)
	}

	if ass.notRaisedErrorToBeChecked() {
		assert.Nil(t, err)
	}
}

// AssertOptionalSingle asserts the result of an optional single against a list of assertions.
func AssertOptionalSingle(ctx context.Context, t *testing.T, optionalSingle OptionalSingle, assertions ...RxAssert) {
	ass := parseAssertions(assertions...)

	var got interface{}
	var err error

	item := <-optionalSingle.Observe()
	if item.IsError() {
		err = item.Err
	} else {
		got = item.Value
	}

	if checkHasValue, value := ass.itemToBeChecked(); checkHasValue {
		assert.Equal(t, value, got)
	}

	if checkHasNoValue, value := ass.noItemToBeChecked(); checkHasNoValue {
		assert.Equal(t, value, got)
	}

	if checkHasRaisedError, expectedError := ass.raisedErrorToBeChecked(); checkHasRaisedError {
		assert.Equal(t, expectedError, err)
	}

	if ass.notRaisedErrorToBeChecked() {
		assert.Nil(t, err)
	}
}
