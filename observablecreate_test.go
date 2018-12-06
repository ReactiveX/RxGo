package rxgo

import (
	"testing"

	"errors"
	"time"

	rxerrors "github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmitsNoElements(t *testing.T) {
	// given
	mockedObserver := NewObserverMock()

	// and
	sequence := Create(func(emitter Observer, disposed bool) {
		emitter.OnDone()
	})

	// when
	sequence.Subscribe(mockedObserver.Capture()).Block()

	// then emits no elements
	mockedObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertCalled(t, "OnDone")
}

func TestEmitsElements(t *testing.T) {
	// given
	mockedObserver := NewObserverMock()

	// and
	elementsToEmit := []int{1, 2, 3, 4, 5}

	// and
	sequence := Create(func(emitter Observer, disposed bool) {
		for _, el := range elementsToEmit {
			emitter.OnNext(el)
		}
		emitter.OnDone()
	})

	// when
	sequence.Subscribe(mockedObserver.Capture()).Block()

	// then emits elements
	for _, emitted := range elementsToEmit {
		mockedObserver.AssertCalled(t, "OnNext", emitted)
	}
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertCalled(t, "OnDone")
}

func TestOnlyFirstDoneCounts(t *testing.T) {
	// given
	mockedObserver := NewObserverMock()

	// and
	sequence := Create(func(emitter Observer, disposed bool) {
		emitter.OnDone()
		emitter.OnDone()
	})

	// when
	sequence.Subscribe(mockedObserver.Capture()).Block()

	// then emits first done
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	mockedObserver.AssertNumberOfCalls(t, "OnDone", 1)
}

func TestDoesntEmitElementsAfterDone(t *testing.T) {
	// given
	mockedObserver := NewObserverMock()

	// and
	sequence := Create(func(emitter Observer, disposed bool) {
		emitter.OnDone()
		emitter.OnNext("it cannot be emitted")
	})

	// when
	sequence.Subscribe(mockedObserver.Capture()).Block()

	// then stops emission after done
	mockedObserver.AssertNotCalled(t, "OnError", mock.Anything)
	mockedObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	mockedObserver.AssertCalled(t, "OnDone")
}

// to clear out error emission
func testEmitsError(t *testing.T) {
	// given
	mockedObserver := NewObserverMock()

	// and
	expectedError := rxerrors.New(rxerrors.UndefinedError, "expected")

	// and
	sequence := Create(func(emitter Observer, disposed bool) {
		emitter.OnError(expectedError)
	})

	// when
	sequence.Subscribe(mockedObserver.Capture()).Block()

	// then emits error
	mockedObserver.AssertCalled(t, "OnError", expectedError)
	mockedObserver.AssertNotCalled(t, "OnNext")
	mockedObserver.AssertNotCalled(t, "OnDone")
}

// to clear out error emission
func testFinishEmissionOnError(t *testing.T) {
	// given
	mockedObserver := NewObserverMock()

	// and
	expectedError := rxerrors.New(rxerrors.UndefinedError, "expected")

	// and
	sequence := Create(func(emitter Observer, disposed bool) {
		emitter.OnError(expectedError)
		emitter.OnNext("some element which cannot be emitted")
		emitter.OnDone()
	})

	// when
	sequence.Subscribe(mockedObserver.Capture()).Block()

	// then emits error
	mockedObserver.AssertCalled(t, "OnError", expectedError)
	mockedObserver.AssertNotCalled(t, "OnNext")
	mockedObserver.AssertNotCalled(t, "OnDone")
}

func TestDefer(t *testing.T) {
	test := 5
	var value int
	onNext := handlers.NextFunc(func(item interface{}) {
		switch item := item.(type) {
		case int:
			value = item
		}
	})
	// First subscriber
	stream1 := Defer(func() Observable {
		return Just(test)
	})
	test = 3
	stream2 := stream1.Map(func(i interface{}) interface{} {
		return i
	})
	stream2.Subscribe(onNext).Block()
	assert.Exactly(t, 3, value)
	// Second subscriber
	test = 8
	stream2 = stream1.Map(func(i interface{}) interface{} {
		return i
	})
	stream2.Subscribe(onNext).Block()
	assert.Exactly(t, 8, value)
}

func TestError(t *testing.T) {
	var got error
	err := errors.New("foo")
	stream := Error(err)
	stream.Subscribe(handlers.ErrFunc(func(e error) {
		got = e
	})).Block()

	assert.Equal(t, err, got)
}

func TestIntervalOperator(t *testing.T) {
	fin := make(chan struct{})
	myStream := Interval(fin, 10*time.Millisecond)
	nums := []int{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			if num >= 5 {
				fin <- struct{}{}
				close(fin)
			}
			nums = append(nums, num)
		}
	})

	myStream.Subscribe(onNext).Block()

	assert.Exactly(t, []int{0, 1, 2, 3, 4, 5}, nums)
}

func TestEmptyCompletesSequence(t *testing.T) {
	// given
	emissionObserver := NewObserverMock()

	// and empty sequence
	sequence := Empty()

	// when subscribes to the sequence
	sequence.Subscribe(emissionObserver.Capture()).Block()

	// then completes without any emission
	emissionObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	emissionObserver.AssertNotCalled(t, "OnError", mock.Anything)
	emissionObserver.AssertCalled(t, "OnDone")
}

func TestNever(t *testing.T) {
	never := Never()
	assert.NotNil(t, never)
}

func TestConcatWithOneObservable(t *testing.T) {
	obs := Concat(Just(1, 2, 3))
	AssertThatObservable(t, obs, HasItems(1, 2, 3))
}

func TestConcatWithTwoObservables(t *testing.T) {
	obs := Concat(Just(1, 2, 3), Just(4, 5, 6))
	AssertThatObservable(t, obs, HasItems(1, 2, 3, 4, 5, 6))
}

func TestConcatWithMoreThanTwoObservables(t *testing.T) {
	obs := Concat(Just(1, 2, 3), Just(4, 5, 6), Just(7, 8, 9))
	AssertThatObservable(t, obs, HasItems(1, 2, 3, 4, 5, 6, 7, 8, 9))
}

func TestConcatWithEmptyObservables(t *testing.T) {
	obs := Concat(Empty(), Empty(), Empty())
	AssertThatObservable(t, obs, IsEmpty())
}

func TestConcatWithAnEmptyObservable(t *testing.T) {
	obs := Concat(Empty(), Just(1, 2, 3))
	AssertThatObservable(t, obs, HasItems(1, 2, 3))

	obs = Concat(Just(1, 2, 3), Empty())
	AssertThatObservable(t, obs, HasItems(1, 2, 3))
}

func TestFromSlice(t *testing.T) {
	obs := FromSlice([]interface{}{1, 2, 3})
	AssertThatObservable(t, obs, HasItems(1, 2, 3))
	AssertThatObservable(t, obs, HasItems(1, 2, 3))
}

func TestFromChannel(t *testing.T) {
	ch := make(chan interface{}, 3)
	obs := FromChannel(ch)

	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	AssertThatObservable(t, obs, HasItems(1, 2, 3))
	AssertThatObservable(t, obs, IsEmpty())
}

func TestJust(t *testing.T) {
	obs := Just(1, 2, 3)
	AssertThatObservable(t, obs, HasItems(1, 2, 3))
	AssertThatObservable(t, obs, HasItems(1, 2, 3))
}

type statefulIterable struct {
	count int
}

func (it *statefulIterable) Next() bool {
	it.count = it.count + 1
	return it.count < 3
}

func (it *statefulIterable) Value() interface{} {
	return it.count
}

func (it *statefulIterable) Iterator() Iterator {
	return it
}

func TestFromStatefulIterable(t *testing.T) {
	obs := FromIterable(&statefulIterable{
		count: -1,
	})

	AssertThatObservable(t, obs, HasItems(0, 1, 2))
	AssertThatObservable(t, obs, IsEmpty())
}

type statelessIterable struct {
	count int
}

func (it *statelessIterable) Next() bool {
	it.count = it.count + 1
	return it.count < 3
}

func (it *statelessIterable) Value() interface{} {
	return it.count
}

func (it *statelessIterable) Iterator() Iterator {
	return &statelessIterable{
		count: -1,
	}
}

func TestFromStatelessIterable(t *testing.T) {
	obs := FromIterable(&statelessIterable{
		count: -1,
	})

	AssertThatObservable(t, obs, HasItems(0, 1, 2))
	AssertThatObservable(t, obs, HasItems(0, 1, 2))
}
