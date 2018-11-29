package observable

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestFlatMapCompletesWhenSequenceIsEmpty(t *testing.T) {
	// given
	emissionObserver := NewObserverMock()

	// and empty sequence
	sequence := Empty()

	// and flattens the sequence with identity
	sequence = sequence.FlatMap(identity, 1)

	// when subscribes to the sequence
	sequence.Subscribe(emissionObserver.Capture()).Block()

	// then completes without any emission
	emissionObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	emissionObserver.AssertNotCalled(t, "OnError", mock.Anything)
	emissionObserver.AssertCalled(t, "OnDone")
}

func TestFlatMapReturnsSameElementBecauseIdentifyApplied(t *testing.T) {
	// given
	emissionObserver := NewObserverMock()

	// and sequence containing one element
	element := 1
	sequence := Just(element)

	// and flattens the sequence with identity
	sequence = sequence.FlatMap(identity, 1)

	// when subscribes to the sequence
	sequence.Subscribe(emissionObserver.Capture()).Block()

	// then completes with emission of the same element
	emissionObserver.AssertNotCalled(t, "OnError", mock.Anything)
	emissionObserver.AssertCalled(t, "OnNext", element)
	emissionObserver.AssertCalled(t, "OnDone")
}

func TestFlatMapReturnsSliceElements(t *testing.T) {
	// given
	emissionObserver := NewObserverMock()

	// and sequence containing slice with few elements
	element1 := "element1"
	element2 := "element2"
	element3 := "element3"
	slice := &([]string{element1, element2, element3})
	sequence := Just(slice)

	// and flattens the sequence with identity
	sequence = sequence.FlatMap(flattenThreeElementSlice, 1)

	// when subscribes to the sequence
	sequence.Subscribe(emissionObserver.Capture()).Block()

	// then completes with emission of flatten elements
	emissionObserver.AssertNotCalled(t, "OnError", mock.Anything)
	emissionObserver.AssertNotCalled(t, "OnNext", slice)
	emissionObserver.AssertCalled(t, "OnNext", element1)
	emissionObserver.AssertCalled(t, "OnNext", element2)
	emissionObserver.AssertCalled(t, "OnNext", element3)
	emissionObserver.AssertCalled(t, "OnDone")
}

// TODO To be reimplemented
//func TestFlatMapUsesForParallelProcessingAtLeast1Process(t *testing.T) {
//	// given
//	emissionObserver := observer.NewObserverMock()
//
//	// and
//	var maxInParallel uint = 0
//
//	// and
//	var requestedMaxInParallel uint = 0
//	flatteningFuncMock := func(out chan interface{}, o Observable, apply func(interface{}) Observable, maxInParallel uint) {
//		requestedMaxInParallel = maxInParallel
//		flatObservedSequence(out, o, apply, maxInParallel)
//	}
//
//	// and flattens the sequence with identity
//	sequence := someSequence.FlatMap(identity, maxInParallel, flatteningFuncMock)
//
//	// when subscribes to the sequence
//	<-sequence.Subscribe(emissionObserver.Capture())
//
//	// then completes with emission of the same element
//	assert.Equal(t, uint(1), requestedMaxInParallel)
//}

var (
	someElement  = "some element"
	someSequence = Just(someElement)
)

func identity(el interface{}) Observable {
	return Just(el)
}

func flattenThreeElementSlice(el interface{}) Observable {
	slice := *(el.(*[]string))
	return Just(slice[0], slice[1], slice[2])
}
