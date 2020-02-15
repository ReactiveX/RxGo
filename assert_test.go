package rxgo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type rxAssert interface {
	apply(*rxAssertImpl)
	itemsToBeChecked() (bool, []interface{})
	noItemsToBeChecked() bool
	raisedErrorToBeChecked() (bool, error)
	notRaisedErrorToBeChecked() bool
}

type rxAssertImpl struct {
	f                      func(*rxAssertImpl)
	checkHasItems          bool
	checkHasNoItems        bool
	items                  []interface{}
	checkHasRaisedError    bool
	error                  error
	checkHasNotRaisedError bool
}

func (ass *rxAssertImpl) apply(do *rxAssertImpl) {
	ass.f(do)
}

func (ass *rxAssertImpl) itemsToBeChecked() (bool, []interface{}) {
	return ass.checkHasItems, ass.items
}

func (ass *rxAssertImpl) noItemsToBeChecked() bool {
	return ass.checkHasNoItems
}

func (ass *rxAssertImpl) raisedErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedError, ass.error
}

func (ass *rxAssertImpl) notRaisedErrorToBeChecked() bool {
	return ass.checkHasNotRaisedError
}

func newAssertion(f func(*rxAssertImpl)) *rxAssertImpl {
	return &rxAssertImpl{
		f: f,
	}
}

func hasItems(items ...interface{}) rxAssert {
	return newAssertion(func(a *rxAssertImpl) {
		a.checkHasItems = true
		a.items = items
	})
}

func hasNoItems() rxAssert {
	return newAssertion(func(a *rxAssertImpl) {
		a.checkHasNoItems = true
	})
}

func hasRaisedError(err error) rxAssert {
	return newAssertion(func(a *rxAssertImpl) {
		a.checkHasRaisedError = true
		a.error = err
	})
}

func hasNotRaisedError() rxAssert {
	return newAssertion(func(a *rxAssertImpl) {
		a.checkHasRaisedError = true
	})
}

func parseAssertions(assertions ...rxAssert) rxAssert {
	ass := new(rxAssertImpl)
	for _, assertion := range assertions {
		assertion.apply(ass)
	}
	return ass
}

func assertObservable(t *testing.T, ctx context.Context, observable Observable, assertions ...rxAssert) {
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
