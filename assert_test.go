package rxgo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type rxAssert interface {
	apply(*rxAssertImpl)
	hasItemsF() (bool, []interface{})
	hasRaisedError() (bool, error)
	hasNotRaisedError() bool
}

type rxAssertImpl struct {
	f                      func(*rxAssertImpl)
	checkHasItems          bool
	hasItems               []interface{}
	checkHasRaisedError    bool
	hasError               error
	checkHasNotRaisedError bool
}

func (ass *rxAssertImpl) apply(do *rxAssertImpl) {
	ass.f(do)
}

func (ass *rxAssertImpl) hasItemsF() (bool, []interface{}) {
	return ass.checkHasItems, ass.hasItems
}

func (ass *rxAssertImpl) hasRaisedError() (bool, error) {
	return ass.checkHasRaisedError, ass.hasError
}

func (ass *rxAssertImpl) hasNotRaisedError() bool {
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
		a.hasItems = items
	})
}

func hasRaisedError(err error) rxAssert {
	return newAssertion(func(a *rxAssertImpl) {
		a.checkHasRaisedError = true
		a.hasError = err
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

	if checkHasItems, expectedItems := ass.hasItemsF(); checkHasItems {
		<-done
		assert.Equal(t, expectedItems, got)
	}

	if checkHasRaisedError, expectedError := ass.hasRaisedError(); checkHasRaisedError {
		assert.Equal(t, expectedError, err)
	}

	if ass.hasNotRaisedError() {
		assert.Nil(t, err)
	}
}
