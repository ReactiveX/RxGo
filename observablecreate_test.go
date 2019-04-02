package rxgo

import (
	"testing"

	"github.com/reactivex/rxgo/options"

	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

func (it *statefulIterable) Next() (interface{}, error) {
	it.count++
	if it.count < 3 {
		return it.count, nil
	}
	return nil, rxerrors.New(rxerrors.EndOfIteratorError)
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

func (it *statelessIterable) Next() (interface{}, error) {
	it.count++
	if it.count < 3 {
		return it.count, nil
	}
	return nil, rxerrors.New(rxerrors.EndOfIteratorError)
}

//func TestFromStatelessIterable(t *testing.T) {
//	obs := FromIterable(&statelessIterable{
//		count: -1,
//	})
//
//	AssertThatObservable(t, obs, HasItems(0, 1, 2))
//	AssertThatObservable(t, obs, HasItems(0, 1, 2))
//}

func TestRange(t *testing.T) {
	obs, err := Range(5, 3)
	if err != nil {
		t.Fail()
	}
	AssertThatObservable(t, obs, HasItems(5, 6, 7, 8))
	AssertThatObservable(t, obs, HasItems(5, 6, 7, 8))
}

func TestRangeWithNegativeCount(t *testing.T) {
	r, err := Range(1, -5)
	assert.NotNil(t, err)
	assert.Nil(t, r)
}

func TestRangeWithMaximumExceeded(t *testing.T) {
	r, err := Range(1<<31, 1)
	assert.NotNil(t, err)
	assert.Nil(t, r)
}

func TestTimer(t *testing.T) {
	d := new(mockDuration)
	d.On("duration").Return(1 * time.Millisecond)

	obs := Timer(d)

	AssertThatObservable(t, obs, HasItems(float64(0)))
	d.AssertCalled(t, "duration")
}

func TestTimerWithNilDuration(t *testing.T) {
	obs := Timer(nil)

	AssertThatObservable(t, obs, HasItems(float64(0)))
}

var _ = Describe("Observable types", func() {
	Context("when creating a cold observable with Just operator", func() {
		observable := Just(1, 2, 3)
		It("should allow multiple subscription to receive the same items", func() {
			outNext1 := make(chan interface{}, 1)
			observable.Subscribe(nextHandler(outNext1))
			Expect(pollItem(outNext1, timeout)).Should(Equal(1))
			Expect(pollItem(outNext1, timeout)).Should(Equal(2))
			Expect(pollItem(outNext1, timeout)).Should(Equal(3))
			Expect(pollItem(outNext1, timeout)).Should(Equal(noData))

			outNext2 := make(chan interface{}, 1)
			observable.Subscribe(nextHandler(outNext2))
			Expect(pollItem(outNext2, timeout)).Should(Equal(1))
			Expect(pollItem(outNext2, timeout)).Should(Equal(2))
			Expect(pollItem(outNext2, timeout)).Should(Equal(3))
			Expect(pollItem(outNext2, timeout)).Should(Equal(noData))
		})
	})

	Context("when creating a cold observable with FromChannel operator", func() {
		ch := make(chan interface{}, 10)
		ch <- 1
		observable := FromChannel(ch)
		It("should allow observer to receive items regardless of the moment it subscribes", func() {
			outNext1 := make(chan interface{}, 1)
			observable.Subscribe(nextHandler(outNext1))
			ch <- 2
			ch <- 3
			Expect(pollItem(outNext1, timeout)).Should(Equal(1))
			Expect(pollItem(outNext1, timeout)).Should(Equal(2))
			Expect(pollItem(outNext1, timeout)).Should(Equal(3))
			Expect(pollItem(outNext1, timeout)).Should(Equal(noData))
		})
	})

	Context("when creating a hot observable with FromEventSource operator without back-pressure strategy", func() {
		ch := make(chan interface{}, 10)
		ch <- 1
		observable := FromEventSource(ch, options.WithoutBackpressureStrategy())
		outNext1 := make(chan interface{}, 1)
		It("should drop an item if there is no subscriber", func() {
			Eventually(len(ch), timeout, pollingInterval).Should(Equal(0))
		})
		It("an observer should receive items depending on the moment it subscribed", func() {
			observable.Subscribe(nextHandler(outNext1))
			ch <- 2
			ch <- 3
			Expect(pollItem(outNext1, timeout)).Should(Equal(2))
			Expect(pollItem(outNext1, timeout)).Should(Equal(3))
			Expect(pollItem(outNext1, timeout)).Should(Equal(noData))
		})
		It("another observer should receive items depending on the moment it subscribed", func() {
			outNext2 := make(chan interface{}, 1)
			observable.Subscribe(nextHandler(outNext2))
			ch <- 4
			ch <- 5
			Expect(pollItem(outNext1, timeout)).Should(Equal(4))
			Expect(pollItem(outNext1, timeout)).Should(Equal(5))
			Expect(pollItem(outNext2, timeout)).Should(Equal(4))
			Expect(pollItem(outNext2, timeout)).Should(Equal(5))
			Expect(pollItem(outNext2, timeout)).Should(Equal(noData))
		})
	})

	Context("when creating a hot observable with FromEventSource operator and a buffer back-pressure strategy", func() {
		ch := make(chan interface{}, 10)
		ch <- 1
		observable := FromEventSource(ch, options.WithBufferBackpressureStrategy(2))
		outNext1 := make(chan interface{})
		It("should drop an item if there is no subscriber", func() {
			Eventually(len(ch), timeout, pollingInterval).Should(Equal(0))
		})
		Context("an observer subscribes", func() {
			observable.Subscribe(nextHandler(outNext1))
			ch <- 2
			ch <- 3
			ch <- 4
			ch <- 5
			It("should consume the messages from the channel", func() {
				Eventually(len(ch), timeout, pollingInterval).Should(Equal(0))
			})
			It("should receive only the buffered items", func() {
				Expect(len(pollItems(outNext1, timeout))).Should(Equal(3))
			})
		})
	})

	Context("when creating a hot observable with FromEventSource operator and a buffer back-pressure strategy", func() {
		ch := make(chan interface{}, 10)
		ch <- 1
		observable := FromEventSource(ch, options.WithBufferBackpressureStrategy(2))
		outNext1 := make(chan interface{})
		outNext2 := make(chan interface{})
		It("should drop an item if there is no subscriber", func() {
			Eventually(len(ch), timeout, pollingInterval).Should(Equal(0))
		})
		Context("two observer subscribe", func() {
			observable.Subscribe(nextHandler(outNext1))
			ch <- 2
			ch <- 3
			ch <- 4
			ch <- 5
			It("should consume the messages from the channel", func() {
				Eventually(len(ch), timeout, pollingInterval).Should(Equal(0))
			})
			observable.Subscribe(nextHandler(outNext2))
			ch <- 6
			ch <- 7
			It("should consume the messages from the channel", func() {
				Eventually(len(ch), timeout, pollingInterval).Should(Equal(0))
			})
			It("the two observer should receive only the buffered items", func() {
				Expect(len(pollItems(outNext1, timeout))).Should(Equal(3))
				Expect(len(pollItems(outNext2, timeout))).Should(Equal(2))
			})
		})
	})

	Context("when creating three observables", func() {
		ch1 := make(chan interface{}, 10)
		ch2 := make(chan interface{}, 10)
		ch3 := make(chan interface{}, 10)
		observable1 := FromChannel(ch1)
		observable2 := FromChannel(ch2)
		observable3 := FromChannel(ch3)
		Context("when merging them using Merge operator", func() {
			ch3 <- 1
			ch2 <- 2
			ch1 <- 3
			ch1 <- 4
			ch3 <- 5
			It("it should produce all the items merged", func() {
				mergedObservable := Merge(observable1, observable2, observable3)
				outNext, _, _ := subscribe(mergedObservable)
				Expect(len(pollItems(outNext, timeout))).Should(Equal(5))
			})
		})
	})
})
