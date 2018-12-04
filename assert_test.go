package rxgo

import (
	"errors"
	"testing"

	"github.com/reactivex/rxgo/optional"
)

func TestAssertThatObservableHasItems(t *testing.T) {
	AssertThatObservable(t, Just(1, 2, 3), HasItems(1, 2, 3))
}

func TestAssertThatObservableHasSize(t *testing.T) {
	AssertThatObservable(t, Just(1, 2, 3), HasSize(3))
}

func TestAssertThatObservableIsEmpty(t *testing.T) {
	AssertThatObservable(t, Empty(), IsEmpty())
}

func TestAssertThatObservableIsNotEmpty(t *testing.T) {
	AssertThatObservable(t, Just(1), IsNotEmpty())
}

func TestAssertThatSingleHasValue(t *testing.T) {
	AssertThatSingle(t, newSingleFrom(1), HasValue(1))
}

func TestAssertThatSingleError(t *testing.T) {
	AssertThatSingle(t, newSingleFrom(errors.New("foo")),
		HasRaisedAnError(), HasRaisedError(errors.New("foo")))
}

func TestAssertThatSingleNotError(t *testing.T) {
	AssertThatSingle(t, newSingleFrom(1), HasNotRaisedAnError())
}

func TestAssertThatOptionalSingleIsEmpty(t *testing.T) {
	AssertThatOptionalSingle(t, newOptionalSingleFrom(optional.Empty()), IsEmpty())
}

func TestAssertThatOptionalSingleIsNotEmpty(t *testing.T) {
	AssertThatOptionalSingle(t, newOptionalSingleFrom(optional.Of(1)), IsNotEmpty())
}

func TestAssertThatOptionalSingleHasValue(t *testing.T) {
	AssertThatOptionalSingle(t, newOptionalSingleFrom(optional.Of(1)), HasValue(1))
}

func TestAssertThatOptionalSingleHasRaisedAnError(t *testing.T) {
	AssertThatOptionalSingle(t, newOptionalSingleFrom(optional.Of(errors.New("foo"))), HasRaisedAnError())
}

func TestAssertThatOptionalSingleHasRaisedError(t *testing.T) {
	AssertThatOptionalSingle(t, newOptionalSingleFrom(optional.Of(errors.New("foo"))), HasRaisedError(errors.New("foo")))
}
