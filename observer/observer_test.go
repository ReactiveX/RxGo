package observer

import (
	"testing"

	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewObserverWithConstructor(t *testing.T) {
	assert := assert.New(t)
	ob := New()

	assert.Implements((*Observer)(nil), ob)
	assert.NotNil(ob.OnNext)
	assert.NotNil(ob.OnError)
	assert.NotNil(ob.OnDone)

	onNext := handlers.NextFunc(func(item interface{}) {})
	onError := handlers.ErrFunc(func(err error) {})

	ob2 := New(onNext, onError)
	assert.NotNil(ob2.OnNext)
	assert.NotNil(ob2.OnError)
	assert.NotNil(ob2.OnDone)

}

func TestCreateNewObserverWithObserver(t *testing.T) {
	nexttext := ""
	donetext := ""

	nextf := handlers.NextFunc(func(item interface{}) {
		if text, ok := item.(string); ok {
			nexttext = text
		}
	})

	donef := handlers.DoneFunc(func() {
		donetext = "Hello"
	})

	ob := New(donef, nextf)

	ob.OnNext("Next")
	ob.OnDone()

	assert.Equal(t, "Next", nexttext)
	assert.Equal(t, "Hello", donetext)
}

func TestChangesOnNextIntoObserver(t *testing.T) {
	// given custom OnNext function
	element := "ready to be changed"
	expectedChange := "changed"
	customOnNext := func(item interface{}) {
		element = expectedChange
	}

	// and gets observer
	gotObserver := New(customOnNext)

	// when reacts on next element
	gotObserver.OnNext("irrevelant")

	// then element is changed
	assert.Equal(t, expectedChange, element)
}

func TestChangesOnErrorIntoObserver(t *testing.T) {
	// given custom OnError function
	element := "ready to be changed"
	expectedChange := "changed"
	customOnError := func(err error) { element = expectedChange }

	// and observer using customOnError
	gotObserver := New(customOnError)

	// when calls observer on error
	gotObserver.OnError(nil)

	// then got observer uses given custom OnError function
	assert.Equal(t, expectedChange, element)
}

func TestChangesOnDoneIntoObserver(t *testing.T) {
	// given custom OnDone function
	element := "ready to be changed"
	expectedChange := "changed"
	customOnDone := func() { element = expectedChange }

	// and observer using custom on done function
	gotObserver := New(customOnDone)

	// when calls observer on done
	gotObserver.OnDone()

	// then got observer uses given custom OnDone function
	assert.Equal(t, expectedChange, element)
}

func TestCustomObserverIsAbleToObserveSequence(t *testing.T) {
	// given custom observer implementing Observer interface
	type customObserver struct {
		Observer
	}
	// and sequence observer created backed by custom observer
	sequenceObserver := customObserver{}

	// when checks event handler
	eventHandler := New(sequenceObserver)

	// then sequence observer is recognized as event handler
	assert.Equal(t, sequenceObserver, eventHandler)
}

