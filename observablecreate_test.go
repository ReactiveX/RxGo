package rxgo

import (
	"testing"

	"errors"
	rxerrors "github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
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
		items := []interface{}{test}
		it, err := iterable.New(items)
		if err != nil {
			t.Fail()
		}
		return From(it)
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

func TestRepeatInfinityOperator(t *testing.T) {
	myStream := Repeat("mystring")

	item, err := myStream.Next()

	if err != nil {
		assert.Fail(t, "fail to emit next item", err)
	}

	if value, ok := item.(string); ok {
		assert.Equal(t, value, "mystring")
	} else {
		assert.Fail(t, "fail to emit next item", err)
	}
}

func TestRepeatNtimeOperator(t *testing.T) {
	myStream := Repeat("mystring", 2)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"mystring", "mystring", "end"}, stringarray)
}

func TestRepeatNtimeMultiVariadicOperator(t *testing.T) {
	myStream := Repeat("mystring", 2, 2, 3, 4, 5, 6, 7)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"mystring", "mystring", "end"}, stringarray)
}

func TestRepeatWithZeroNtimeOperator(t *testing.T) {
	myStream := Repeat("mystring", 0)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"end"}, stringarray)
}

func TestRepeatWithNegativeTimesOperator(t *testing.T) {
	myStream := Repeat("mystring", -10)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"end"}, stringarray)
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
