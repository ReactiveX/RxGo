package grx

import (
	"testing"
	"time"
	"net/http"

	"github.com/stretchr/testify/assert"

	"github.com/jochasinga/grx/event"
	"github.com/jochasinga/grx/observable"
	"github.com/jochasinga/grx/observer"
)

// TestCreateObservableWithCOnstructor tests if the constructor method returns an Observable
func TestCreateObservableWithConstructor(t *testing.T) {
	testStream := observable.New()
	assert.IsType(t, &observable.Observable{}, testStream)
}

// TestCreateObservableWithEmpty tests if Empty creates an empty Observable that terminates right away.
func TestCreateObservableWithEmpty(t *testing.T) {
	msg := "Sumpin's"
	testStream := observable.Empty()
	testObserver := &observer.Observer{
		OnCompleted: func(e *event.Event) {
			msg += " brewin'"
		},
	}
	testStream.Subscribe(testObserver)
	assert.Equal(t, "Sumpin's brewin'", msg)
}

// TestCreateObservableWithJust tests if Just method returns an <*Observable>
func TestCreateObservableWithJust(t *testing.T) {
	
	// Provided a URL string
	url := "http://api.com/api/v1.0/user"

	// and an Observable created with Just method
	testStream := observable.Just(url)

	// Observable must have the type <Observable>
	assert.IsType(t, &observable.Observable{}, testStream)
}

// TestCreateObservableWithFrom tests if From methods returns an <*Observable>
func TestCreateObservableWithFrom(t *testing.T) {
	// Provided a slice of URL strings
	urls := []interface{}{
		"http://api.com/api/v1.0/user",
		"https://dropbox.com/api/v2.1/get/storage",
		"http://googleapi.com/map",
	}

	// and an Observable created with From method
	testStream := observable.From(urls)

	// Observable must have the type <*Observable>
	assert.IsType(t, &observable.Observable{}, testStream)
}

func TestCreateObservableWithStart(t *testing.T) {

	// Provided a "directive" function which return an event that emits
	// a slice of integer numbers from 0 to 20.
	directive := func() *event.Event {
		ev := &event.Event{ Value: []int{} }
		for i:=0; i<=20; i++ {
			ev.Value = append(ev.Value.([]int), i)
		}
		return ev
	}

	testStream := observable.Start(directive)

	// Make sure it's the right type
	assert.IsType(t, &observable.Observable{}, testStream)
	nums := []int{}
	testObserver := &observer.Observer{
		OnNext: func(ev *event.Event) {
			for _, num := range ev.Value.([]int) {
				nums = append(nums, num * 2)
			}
		},
		OnCompleted: func(ev *event.Event) {
			nums = append(nums, 666)
		},
	}
	
	testStream.Subscribe(testObserver)
	expected := []int{}
	for i:=0; i<=20; i++ {
		expected = append(expected, i * 2)
	}
	expected = append(expected, 666)
	assert.Exactly(t, expected, nums)
}

func TestSubscribeToJustObservable1(t *testing.T) {
	urlWithUserID := ""
	
	// Provided an Observable created with Just method
	url := "http://api.com/api/v1.0/user"
	expected := url + "?id=999"
	testStream := observable.Just(url)
	
	testObserver := &observer.Observer{
		OnNext: func(e *event.Event) { urlWithUserID = e.Value.(string) },
		OnCompleted: func(e *event.Event) { urlWithUserID = urlWithUserID + "?id=999" },
	}
	testStream.Subscribe(testObserver)
	assert.Exactly(t, expected, urlWithUserID)
}

func TestSubscribeToJustObservable2(t *testing.T) {
	nums := []int{}

	// Create an Observable with an integer.
	testStream := observable.Just(1)

	// Create an Observer object.
	testObserver := &observer.Observer{
		// While there is more events down the stream (in this case, just one), add 1 to the value emitted.
		OnNext: func(e *event.Event) {
			num := e.Value.(int) + 1
			nums = append(nums, num)
		},
			
		// If an error is encountered at any time, panic.
		OnError: func(e *event.Event) { panic(e.Error) },
			
		// When the stream comes to an end, append 0 to the slice.
		OnCompleted: func(e *event.Event) { nums = append(nums, 0) },
	}
	
	// Start listening to stream
	testStream.Subscribe(testObserver)
	assert.Exactly(t, []int{2, 0}, nums)
}

func TestSubscribeToFromObservable(t *testing.T) {
	nums := []interface{}{1, 2, 3, 4, 5, 6}
	numCopy := []int{}

	// Create an Observable with an integer.
	testObservable := observable.From(nums)
	testObserver := &observer.Observer{
		// While there is more events down the stream (in this case, just one), add 1 to the value emitted.
		OnNext: func(e *event.Event) {
			numCopy = append(numCopy, e.Value.(int) + 1)
		},
			
		// When the stream comes to an end, append 0 to the slice.
		OnCompleted: func(e *event.Event) {
			numCopy = append(numCopy, 0)
		},

	}

	// Start listening to stream
	testObservable.Subscribe(testObserver)
	assert.Exactly(t, []int{2, 3, 4, 5, 6, 7, 0}, numCopy)
}

func TestStartMethodWithFakeExternalCalls(t *testing.T) {

	fakeResponses := []*http.Response{}

	// Fake directives that returns an Event containing an HTTP response.
	directive1 := func() *event.Event {
		res := &http.Response{
			Status: "404 NOT FOUND",
			StatusCode: 404,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
		}
		time.Sleep(time.Millisecond * 20)
		return &event.Event{ Value: res }
	}

	directive2 := func() *event.Event {
		res := &http.Response{
			Status: "200 OK",
			StatusCode: 200,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
		}
		time.Sleep(10 * time.Millisecond)
		return &event.Event{ Value: res }
	}

	directive3 := func() *event.Event {
		res := &http.Response{
			Status: "500 SERVER ERROR",
			StatusCode: 500,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
		}
		time.Sleep(30 * time.Millisecond)
		return &event.Event{ Value: res }
	}

	testObserver := &observer.Observer{
		OnNext: func(ev *event.Event) {
			fakeResponses = append(fakeResponses, ev.Value.(*http.Response))
		},
		OnCompleted: func(ev *event.Event) {
			fakeResponses = append(fakeResponses, &http.Response{
					Status: "999 End",
					StatusCode: 999,
			})
		},
	}

	testObservable := observable.Start(directive1, directive2, directive3).Subscribe(testObserver)

	assert := assert.New(t)

	assert.IsType(&observable.Observable{}, testObservable)
	assert.Equal(4, len(fakeResponses))
	assert.Equal(200, fakeResponses[0].StatusCode)
	assert.Equal(404, fakeResponses[1].StatusCode)
	assert.Equal(500, fakeResponses[2].StatusCode)
	assert.Equal(999, fakeResponses[3].StatusCode)
}

func TestCreateObservableWithInterval(t *testing.T) {
	numch := make(chan int, 1)
	go func() {
		_ = observable.Interval(time.Millisecond).Subscribe(&observer.Observer{
			OnNext: func(e *event.Event) {
				numch <- e.Value.(int)
			},
		})
	}()

	for i := 0; i <= 10; i++ {
		<-time.After(1 * time.Millisecond)
		assert.Equal(t, i, <-numch)
	}
}

func TestSubscriptionDoesNotBlock(t *testing.T) {
	testStream := observable.Just("Hello")
	
	testObserver := &observer.Observer{
		OnNext: func(ev *event.Event) {
			time.Sleep(1 * time.Second)
			return
		},
		OnCompleted: func(ev *event.Event) {
			t.Log("DONE")
		},
	}


	first := time.Now()
	testStream.Subscribe(testObserver)
	elapsed := time.Since(first)

	comp := assert.Comparison(func() bool {
		return elapsed < 1 * time.Second
	})

	assert.Condition(t, comp)
}





