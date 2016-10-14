package grx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCreateObservableWithCOnstructor tests if the constructor method returns an Observable
func TestCreateObservableWithConstructor(t *testing.T) {
	testStream := NewObservable("myStream")
	assert.IsType(t, &Observable{}, testStream)
}

// TestCreateObservableWithEmpty tests if Empty creates an empty Observable that terminates right away.
func TestCreateObservableWithEmpty(t *testing.T) {
	msg := "Sumpin's"
	observable := Empty()
	observable.Subscribe(&Observer{
		OnCompleted: func(e *Event) { msg = msg + " brewin'" },
	})
	assert.Equal(t, "Sumpin's brewin'", msg)
}

// TestCreateObservableWithJust tests if Just method returns an <*Observable>
func TestCreateObservableWithJust(t *testing.T) {
	
	// Provided a URL string
	url := "http://api.com/api/v1.0/user"

	// and an Observable created with Just method
	observable := Just(url)

	// Observable must have the type <Observable>
	assert.IsType(t, &Observable{}, observable)
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
	observable := From(urls)

	// Observable must have the type <*Observable>
	assert.IsType(t, &Observable{}, observable)
}

func TestCreateObservableWithStart(t *testing.T) {

	// Provided a "directive" function which return an event that emits
	// a slice of integer numbers from 0 to 20.
	directive := func() *Event {
		ev := &Event{ Value: []int{} }
		for i:=0; i<=20; i++ {
			<-time.After(1000)
			ev.Value = append(ev.Value.([]int), i)
		}
		return ev
	}

	observable := Start(directive)

	// Make sure it's the right type
	assert.IsType(t, &Observable{}, observable)
	nums := []int{}
	obs := &Observer{
		OnNext: func(ev *Event) {
			for _, num := range ev.Value.([]int) {
				nums = append(nums, num * 2)
			}
		},
		OnCompleted: func(ev *Event) {
			nums = append(nums, 666)
		},
	}
	
	observable.Subscribe(obs)
	expected := []int{}
	for i:=0; i<=20; i++ {
		expected = append(expected, i*2)
	}
	expected = append(expected, 666)
	assert.Exactly(t, expected, nums)
}

func TestSubscribeToJustObservable1(t *testing.T) {
	urlWithUserID := ""
	
	// Provided an Observable created with Just method
	url := "http://api.com/api/v1.0/user"
	expected := url + "?id=999"

	observable := Just(url)
	obs := &Observer{
		OnNext: func(e *Event) { urlWithUserID = e.Value.(string) },
		OnCompleted: func(e *Event) { urlWithUserID = urlWithUserID + "?id=999" },
	}
	observable.Subscribe(obs)
	assert.Exactly(t, expected, urlWithUserID)
}

func TestSubscribeToJustObservable2(t *testing.T) {
	nums := []int{}

	// Create an Observable with an integer.
	numObservable := Just(1)

	// Create an Observer object.
	observer := &Observer{
		// While there is more events down the stream (in this case, just one), add 1 to the value emitted.
		OnNext: func(e *Event) {
			num := e.Value.(int) + 1
			nums = append(nums, num)
		},
			
		// If an error is encountered at any time, panic.
		OnError: func(e *Event) { panic(e.Error) },
			
		// When the stream comes to an end, append 0 to the slice.
		OnCompleted: func(e *Event) { nums = append(nums, 0) },
	}
	
	// Start listening to stream
	numObservable.Subscribe(observer)
	assert.Exactly(t, []int{2, 0}, nums)
}

func TestSubscribeToFromObservable(t *testing.T) {
	nums := []interface{}{1, 2, 3, 4, 5, 6}
	numCopy := []int{}

	// Create an Observable with an integer.
	numObservable := From(nums)
	observer := &Observer{
		// While there is more events down the stream (in this case, just one), add 1 to the value emitted.
		OnNext: func(e *Event) {
			numCopy = append(numCopy, e.Value.(int) + 1)
		},
			
		// When the stream comes to an end, append 0 to the slice.
		OnCompleted: func(e *Event) {
			numCopy = append(numCopy, 0)
		},

	}

	// Start listening to stream
	numObservable.Subscribe(observer)
	assert.Exactly(t, []int{2, 3, 4, 5, 6, 7, 0}, numCopy)
}



