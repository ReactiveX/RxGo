package rxgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertPredicate is a custom predicate based on the items.
type AssertPredicate func(items []interface{}) error

// RxAssert lists the Observable assertions.
type RxAssert interface {
	apply(*rxAssert)
	itemsToBeChecked() (bool, []interface{})
	itemsNoOrderedToBeChecked() (bool, []interface{})
	noItemsToBeChecked() bool
	someItemsToBeChecked() bool
	raisedErrorToBeChecked() (bool, error)
	raisedErrorsToBeChecked() (bool, []error)
	raisedAnErrorToBeChecked() (bool, error)
	notRaisedErrorToBeChecked() bool
	itemToBeChecked() (bool, interface{})
	noItemToBeChecked() (bool, interface{})
	customPredicatesToBeChecked() (bool, []AssertPredicate)
}

type rxAssert struct {
	f                       func(*rxAssert)
	checkHasItems           bool
	checkHasNoItems         bool
	checkHasSomeItems       bool
	items                   []interface{}
	checkHasItemsNoOrder    bool
	itemsNoOrder            []interface{}
	checkHasRaisedError     bool
	err                     error
	checkHasRaisedErrors    bool
	errs                    []error
	checkHasRaisedAnError   bool
	checkHasNotRaisedError  bool
	checkHasItem            bool
	item                    interface{}
	checkHasNoItem          bool
	checkHasCustomPredicate bool
	customPredicates        []AssertPredicate
}

func (ass *rxAssert) apply(do *rxAssert) {
	ass.f(do)
}

func (ass *rxAssert) itemsToBeChecked() (bool, []interface{}) {
	return ass.checkHasItems, ass.items
}

func (ass *rxAssert) itemsNoOrderedToBeChecked() (bool, []interface{}) {
	return ass.checkHasItemsNoOrder, ass.itemsNoOrder
}

func (ass *rxAssert) noItemsToBeChecked() bool {
	return ass.checkHasNoItems
}

func (ass *rxAssert) someItemsToBeChecked() bool {
	return ass.checkHasSomeItems
}

func (ass *rxAssert) raisedErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedError, ass.err
}

func (ass *rxAssert) raisedErrorsToBeChecked() (bool, []error) {
	return ass.checkHasRaisedErrors, ass.errs
}

func (ass *rxAssert) raisedAnErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedAnError, ass.err
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

func (ass *rxAssert) customPredicatesToBeChecked() (bool, []AssertPredicate) {
	return ass.checkHasCustomPredicate, ass.customPredicates
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

// HasItem checks if a single or optional single has a specific item.
func HasItem(i interface{}) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasItem = true
		a.item = i
	})
}

// HasItemsNoOrder checks that an observable produces the corresponding items regardless of the order.
func HasItemsNoOrder(items ...interface{}) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasItemsNoOrder = true
		a.itemsNoOrder = items
	})
}

// IsNotEmpty checks that the observable produces some items.
func IsNotEmpty() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasSomeItems = true
	})
}

// IsEmpty checks that the observable has not produce any item.
func IsEmpty() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasNoItems = true
	})
}

// HasError checks that the observable has produce a specific error.
func HasError(err error) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedError = true
		a.err = err
	})
}

// HasAnError checks that the observable has produce an error.
func HasAnError() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedAnError = true
	})
}

// HasErrors checks that the observable has produce a set of errors.
func HasErrors(errs ...error) RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedErrors = true
		a.errs = errs
	})
}

// HasNoError checks that the observable has not raised any error.
func HasNoError() RxAssert {
	return newAssertion(func(a *rxAssert) {
		a.checkHasRaisedError = true
	})
}

// CustomPredicate checks a custom predicate.
func CustomPredicate(predicate AssertPredicate) RxAssert {
	return newAssertion(func(a *rxAssert) {
		if !a.checkHasCustomPredicate {
			a.checkHasCustomPredicate = true
			a.customPredicates = make([]AssertPredicate, 0)
		}
		a.customPredicates = append(a.customPredicates, predicate)
	})
}

func parseAssertions(assertions ...RxAssert) RxAssert {
	ass := new(rxAssert)
	for _, assertion := range assertions {
		assertion.apply(ass)
	}
	return ass
}

// Assert asserts the result of an iterable against a list of assertions.
func Assert(ctx context.Context, t *testing.T, iterable Iterable, assertions ...RxAssert) {
	ass := parseAssertions(assertions...)

	got := make([]interface{}, 0)
	errs := make([]error, 0)

	observe := iterable.Observe()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case item, ok := <-observe:
			if !ok {
				break loop
			}
			if item.Error() {
				errs = append(errs, item.E)
			} else {
				got = append(got, item.V)
			}
		}
	}

	if checked, predicates := ass.customPredicatesToBeChecked(); checked {
		for _, predicate := range predicates {
			err := predicate(got)
			if err != nil {
				assert.Fail(t, err.Error())
			}
		}
	}
	if checkHasItems, expectedItems := ass.itemsToBeChecked(); checkHasItems {
		assert.Equal(t, expectedItems, got)
	}
	if checkHasItemsNoOrder, itemsNoOrder := ass.itemsNoOrderedToBeChecked(); checkHasItemsNoOrder {
		m := make(map[interface{}]interface{})
		for _, v := range itemsNoOrder {
			m[v] = nil
		}

		for _, v := range got {
			delete(m, v)
		}
		if len(m) != 0 {
			assert.Fail(t, "missing elements", "%v", got)
		}
	}
	if checkHasItem, value := ass.itemToBeChecked(); checkHasItem {
		length := len(got)
		if length != 1 {
			assert.FailNow(t, "wrong number of items", "expected 1, got %d", length)
		}
		assert.Equal(t, value, got[0])
	}
	if ass.noItemsToBeChecked() {
		assert.Equal(t, 0, len(got))
	}
	if ass.someItemsToBeChecked() {
		assert.NotEqual(t, 0, len(got))
	}
	if checkHasRaisedError, expectedError := ass.raisedErrorToBeChecked(); checkHasRaisedError {
		if expectedError == nil {
			assert.Equal(t, 0, len(errs))
		} else {
			if len(errs) == 0 {
				assert.FailNow(t, "no error raised", "expected %v", expectedError)
			}
			assert.Equal(t, expectedError, errs[0])
		}
	}
	if checkHasRaisedErrors, expectedErrors := ass.raisedErrorsToBeChecked(); checkHasRaisedErrors {
		assert.Equal(t, expectedErrors, errs)
	}
	if checkHasRaisedAnError, expectedError := ass.raisedAnErrorToBeChecked(); checkHasRaisedAnError {
		assert.Nil(t, expectedError)
	}
	if ass.notRaisedErrorToBeChecked() {
		assert.Equal(t, 0, len(errs))
	}
}
