package rxgo

import (
	"errors"
	"testing"
)

func TestAssertThatObservableHasItems(t *testing.T) {
	AssertObservable(t, Just(1, 2, 3), HasItems(1, 2, 3))
}

func TestAssertThatObservableHasSize(t *testing.T) {
	AssertObservable(t, Just(1, 2, 3), HasSize(3))
}

func TestAssertThatObservableIsEmpty(t *testing.T) {
	AssertObservable(t, Empty(), IsEmpty())
}

func TestAssertThatObservableIsNotEmpty(t *testing.T) {
	AssertObservable(t, Just(1), IsNotEmpty())
}

func TestAssertThatSingleHasValue(t *testing.T) {
	AssertSingle(t, newSingleFrom(1), HasValue(1))
}

func TestAssertThatSingleError(t *testing.T) {
	AssertSingle(t, newSingleFrom(errors.New("foo")),
		HasRaisedAnError(), HasRaisedError(errors.New("foo")))
}

func TestAssertThatSingleNotError(t *testing.T) {
	AssertSingle(t, newSingleFrom(1), HasNotRaisedAnyError())
}

func TestAssertThatOptionalSingleIsEmpty(t *testing.T) {
	AssertOptionalSingle(t, newOptionalSingleFrom(EmptyOptional()), IsEmpty())
}

func TestAssertThatOptionalSingleIsNotEmpty(t *testing.T) {
	AssertOptionalSingle(t, newOptionalSingleFrom(Of(1)), IsNotEmpty())
}

func TestAssertThatOptionalSingleHasValue(t *testing.T) {
	AssertOptionalSingle(t, newOptionalSingleFrom(Of(1)), HasValue(1))
}

func TestAssertThatOptionalSingleHasRaisedAnError(t *testing.T) {
	AssertOptionalSingle(t, newOptionalSingleFrom(Of(errors.New("foo"))), HasRaisedAnError())
}

func TestAssertThatOptionalSingleHasRaisedError(t *testing.T) {
	AssertOptionalSingle(t, newOptionalSingleFrom(Of(errors.New("foo"))), HasRaisedError(errors.New("foo")))
}
