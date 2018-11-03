package observable

import (
	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/observer"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestEmitsNoElements(t *testing.T) {
	// given
	mockedObserver := observer.NewObserverMock()

	// and
	sequence := Create(func(emitter *observer.Observer, disposed bool) {
		emitter.OnDone()
	})

	// when
	<-sequence.Subscribe(mockedObserver.Capture())

	// then emits no elements
	mockedObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertCalled(t, "OnDone")
}

func TestEmitsElements(t *testing.T) {
	// given
	mockedObserver := observer.NewObserverMock()

	// and
	elementsToEmit := []int{1, 2, 3, 4, 5}

	// and
	sequence := Create(func(emitter *observer.Observer, disposed bool) {
		for _, el := range elementsToEmit {
			emitter.OnNext(el)
		}
		emitter.OnDone()
	})

	// when
	<-sequence.Subscribe(mockedObserver.Capture())

	// then emits elements
	for _, emitted := range elementsToEmit {
		mockedObserver.AssertCalled(t, "OnNext", emitted)
	}
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertCalled(t, "OnDone")
}

func TestOnlyFirstDoneCounts(t *testing.T) {
	// given
	mockedObserver := observer.NewObserverMock()

	// and
	sequence := Create(func(emitter *observer.Observer, disposed bool) {
		emitter.OnDone()
		emitter.OnDone()
	})

	// when
	<-sequence.Subscribe(mockedObserver.Capture())

	// then emits first done
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	mockedObserver.AssertNumberOfCalls(t, "OnDone", 1)
}

func TestDoesntEmitElementsAfterDone(t *testing.T) {
	// given
	mockedObserver := observer.NewObserverMock()

	// and
	sequence := Create(func(emitter *observer.Observer, disposed bool) {
		emitter.OnDone()
		emitter.OnNext("it cannot be emitted")
	})

	// when
	<-sequence.Subscribe(mockedObserver.Capture())

	// then stops emission after done
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	mockedObserver.AssertCalled(t, "OnDone")
}


// to clear out error emission
func testEmitsError(t *testing.T) {
	// given
	mockedObserver := observer.NewObserverMock()

	// and
	expectedError := errors.New(errors.UndefinedError, "expected")

	// and
	sequence := Create(func(emitter *observer.Observer, disposed bool) {
		emitter.OnError(expectedError)
	})

	// when
	<-sequence.Subscribe(mockedObserver.Capture())

	// then emits error
	mockedObserver.AssertCalled(t, "OnError", expectedError)
	mockedObserver.AssertNotCalled(t, "OnNext")
	mockedObserver.AssertNotCalled(t, "OnDone")
}

// to clear out error emission
func testFinishEmissionOnError(t *testing.T) {
	// given
	mockedObserver := observer.NewObserverMock()

	// and
	expectedError := errors.New(errors.UndefinedError, "expected")

	// and
	sequence := Create(func(emitter *observer.Observer, disposed bool) {
		emitter.OnError(expectedError)
		emitter.OnNext("some element which cannot be emitted")
		emitter.OnDone()
	})

	// when
	<-sequence.Subscribe(mockedObserver.Capture())

	// then emits error
	mockedObserver.AssertCalled(t, "OnError", expectedError)
	mockedObserver.AssertNotCalled(t, "OnNext")
	mockedObserver.AssertNotCalled(t, "OnDone")
}
