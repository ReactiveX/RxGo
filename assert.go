package rxgo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// RxAssertion lists the assertions which may be configured on an Observable.
type RxAssertion interface {
	apply(*assertion)
	hasItemsFunc() (bool, []interface{})
	hasSizeFunc() (bool, int)
	hasValueFunc() (bool, interface{})
	hasRaisedErrorFunc() (bool, error)
	hasRaisedAnErrorFunc() bool
	hasNotRaisedAnErrorFunc() bool
	isEmptyFunc() (bool, bool)
}

type assertion struct {
	f                        func(*assertion)
	checkHasItems            bool
	hasItems                 []interface{}
	checkHasSize             bool
	hasSize                  int
	checkHasValue            bool
	hasValue                 interface{}
	checkHasRaisedError      bool
	hasRaisedError           error
	checkHasRaisedAnError    bool
	checkHasNotRaisedAnError bool
	checkIsEmpty             bool
	isEmpty                  bool
}

func (ass *assertion) hasItemsFunc() (bool, []interface{}) {
	return ass.checkHasItems, ass.hasItems
}

func (ass *assertion) hasSizeFunc() (bool, int) {
	return ass.checkHasSize, ass.hasSize
}

func (ass *assertion) hasValueFunc() (bool, interface{}) {
	return ass.checkHasValue, ass.hasValue
}

func (ass *assertion) hasRaisedErrorFunc() (bool, error) {
	return ass.checkHasRaisedError, ass.hasRaisedError
}

func (ass *assertion) hasRaisedAnErrorFunc() bool {
	return ass.checkHasRaisedAnError
}

func (ass *assertion) hasNotRaisedAnErrorFunc() bool {
	return ass.checkHasNotRaisedAnError
}

func (ass *assertion) isEmptyFunc() (bool, bool) {
	return ass.checkIsEmpty, ass.isEmpty
}

func (ass *assertion) apply(do *assertion) {
	ass.f(do)
}

func newAssertion(f func(*assertion)) *assertion {
	return &assertion{
		f: f,
	}
}

func parseAssertions(assertions ...RxAssertion) RxAssertion {
	a := new(assertion)
	for _, assertion := range assertions {
		assertion.apply(a)
	}
	return a
}

// HasItems checks that an observable produces the corresponding items.
func HasItems(items ...interface{}) RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasItems = true
		a.hasItems = items
	})
}

// HasSize checks that an observable produces the corresponding number of items.
func HasSize(size int) RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasSize = true
		a.hasSize = size
	})
}

// IsEmpty checks that an observable produces zero items.
func IsEmpty() RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkIsEmpty = true
		a.isEmpty = true
	})
}

// IsNotEmpty checks that an observable produces items.
func IsNotEmpty() RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkIsEmpty = true
		a.isEmpty = false
	})
}

// HasValue checks that a single produces the corresponding value.
func HasValue(value interface{}) RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasValue = true
		a.hasValue = value
	})
}

// HasRaisedError checks that a single raises the corresponding error.
func HasRaisedError(err error) RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasRaisedError = true
		a.hasRaisedError = err
	})
}

// HasRaisedAnError checks that a single raises an error.
func HasRaisedAnError() RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasRaisedAnError = true
	})
}

// HasNotRaisedAnyError checks that a single does not raise an error.
func HasNotRaisedAnyError() RxAssertion {
	return newAssertion(func(a *assertion) {
		a.checkHasNotRaisedAnError = true
	})
}

func assertObservable(t *testing.T, ass RxAssertion, got []interface{}, err error) {
	checkHasItems, items := ass.hasItemsFunc()
	if checkHasItems {
		assert.Equal(t, items, got)
	}

	checkHasSize, size := ass.hasSizeFunc()
	if checkHasSize {
		assert.Equal(t, size, len(got))
	}

	checkIsEmpty, empty := ass.isEmptyFunc()
	if checkIsEmpty {
		if empty {
			assert.Equal(t, 0, len(got))
		} else {
			assert.NotEqual(t, 0, len(got))
		}
	}

	checkHasRaisedAnError := ass.hasRaisedAnErrorFunc()
	if checkHasRaisedAnError {
		assert.NotNil(t, err)
	}

	checkHasRaisedError, value := ass.hasRaisedErrorFunc()
	if checkHasRaisedError {
		assert.Equal(t, value, err)
	}

	checkHasNotRaisedError := ass.hasNotRaisedAnErrorFunc()
	if checkHasNotRaisedError {
		assert.Nil(t, err)
	}
}

// AssertObservable asserts the result of an Observable against a list of assertions.
func AssertObservable(t *testing.T, observable Observable, assertions ...RxAssertion) {
	ass := parseAssertions(assertions...)
	got := make([]interface{}, 0)
	var err error
	observable.ForEach(func(i interface{}) {
		got = append(got, i)
	}, func(e error) {
		err = e
	}, nil).Block()

	assertObservable(t, ass, got, err)
}

// AssertObservableEventually asserts eventually the result of an Observable against a list of assertions.
func AssertObservableEventually(t *testing.T, observable Observable, timeout time.Duration, assertions ...RxAssertion) {
	ass := parseAssertions(assertions...)

	chItem := make(chan interface{}, 1)
	defer close(chItem)
	chErr := make(chan error)
	defer close(chErr)

	got := make([]interface{}, 0)
	var err error

	observable.ForEach(func(i interface{}) {
		chItem <- i
	}, func(e error) {
		chErr <- e
	}, nil)
	ctxTimeout, ctxTimeoutF := context.WithTimeout(context.Background(), timeout)

mainLoop:
	for {
		select {
		case item, open := <-chItem:
			if open {
				got = append(got, item)
			} else {
				ctxTimeoutF()
				break mainLoop
			}
		case e := <-chErr:
			err = e
		case <-ctxTimeout.Done():
			break mainLoop
		}
	}

	assertObservable(t, ass, got, err)
}

// AssertSingle asserts the result of a Single against a list of assertions.
func AssertSingle(t *testing.T, single Single, assertions ...RxAssertion) {
	ass := parseAssertions(assertions...)

	v, err := single.Subscribe(nil).Block()

	checkHasValue, value := ass.hasValueFunc()
	if checkHasValue {
		if err != nil {
			assert.Fail(t, "error is not nil")
		} else {
			assert.Equal(t, value, v)
		}
	}

	checkHasRaisedAnError := ass.hasRaisedAnErrorFunc()
	if checkHasRaisedAnError {
		assert.NotNil(t, err)
	}

	checkHasRaisedError, value := ass.hasRaisedErrorFunc()
	if checkHasRaisedError {
		assert.Equal(t, value, err)
	}

	checkHasNotRaisedError := ass.hasNotRaisedAnErrorFunc()
	if checkHasNotRaisedError {
		assert.Nil(t, err)
	}
}

// AssertOptionalSingle asserts the result of an OptionalSingle against a list of assertions.
func AssertOptionalSingle(t *testing.T, optionalSingle OptionalSingle, assertions ...RxAssertion) {
	ass := parseAssertions(assertions...)

	v, err := optionalSingle.Subscribe(nil).Block()

	if err != nil {
		assert.Fail(t, "error while retrieving OptionalSingle results")
	}

	if optional, ok := v.(Optional); ok {
		checkIsEmpty, empty := ass.isEmptyFunc()
		if checkIsEmpty {
			if empty {
				assert.True(t, optional.IsEmpty())
			} else {
				assert.False(t, optional.IsEmpty())
			}
		}

		got, emptyError := optional.Get()

		checkHasRaisedAnError := ass.hasRaisedAnErrorFunc()
		if checkHasRaisedAnError {
			assert.Nil(t, emptyError)
			assert.IsType(t, errors.New(""), got)
		}

		checkHasRaisedError, emptyError := ass.hasRaisedErrorFunc()
		if checkHasRaisedError {
			assert.Equal(t, emptyError, got)
		}

		checkHasValue, value := ass.hasValueFunc()
		if checkHasValue {
			assert.Equal(t, value, got)
		}
	} else {
		assert.Fail(t, "OptionalSingle did not produce an Optional")
	}
}
