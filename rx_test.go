package rx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateObservable(t *testing.T) {
	testStream := NewObservable("myStream")
	assert.IsType(t, &Observable{}, testStream)
}

func TestCreateObservableWithJust(t *testing.T) {
	// Provided a URL string
	url := "http://api.com/api/v1.0/user"

	// and an Observable created with Just method
	observable := Just(url)

	// Observable must have the type <Observable>
	assert.IsType(t, &Observable{}, observable)
}

func TestCreateObservableWithFrom(t *testing.T) {
	// Provided a slice of URL strings
	urls := []interface{}{
		"http://api.com/api/v1.0/user",
		"https://dropbox.com/api/v2.1/get/storage",
		"http://googleapi.com/map",
	}

	// and an Observable created with From method
	observable := From(urls)

	// Observable must have the type <Observable>
	assert.IsType(t, &Observable{}, observable)
}

func TestCreateObservableWithStart(t *testing.T) {
	myFunc := func() Event {
		ev := Event{ Value: []int{} }

		for i:=0; i<=100; i++ {
			<-time.After(1000)
			ev.Value = append(ev.Value.([]int), i)
		}

		return ev
	}

	observable := Start(myFunc)

	// Make sure it's the right type
	assert.IsType(t, &Observable{}, observable)

	obs := &Observer{
		OnNext: func(e Event) {
			t.Logf("This event has value: %v", e.Value)
		},
	}

	observable.Subscribe(obs)
}

func TestSubscribingtoJustObservable(t *testing.T) {

	// Provided an Observable created with Just method
	url := "http://api.com/api/v1.0/user"
	observable := Just(url)
	/*
	obs := func(e Event) {
		t.Log(e.Value)
	}
        */
	obs := &Observer{
		OnNext: func(e Event) { t.Logf("This event has value: %v", e.Value) },
		OnError: func(e Event) { t.Errorf("This event has error: %v", e.Error) },
		OnCompleted: func(e Event) { t.Logf("This event has last value: %v", e.Value) },
	}

	observable.Subscribe(obs)
}

func TestSubscribingtoFromObservable(t *testing.T) {

	urls := []interface{}{
		"http://api.com/api/v1.0/user",
		"https://dropbox.com/api/v2.1/get/storage",
		"http://googleapi.com/map",
	}
	toBeFilled := []interface{}{}
	
	observable := From(urls)

	obs := &Observer{
		OnNext: func(e Event) {
			t.Logf("This event has value: %v", e.Value)
			toBeFilled = append(toBeFilled, e.Value)
		},
		OnError: func(e Event) { t.Errorf("This event has error: %v", e.Error) },
		OnCompleted: func(e Event) { t.Logf("This event has last value: %v", e.Value) },
	}

	observable.Subscribe(obs)

	for _, url := range toBeFilled {
		if len(url.(string)) <= 0 {
			t.Error("fack")
		}
	}

	/*
	go func() {
		result, err := func() {
			// very slow response
			<-time.After(1000)
			observable.Add(Event{
				Value: "Hello result",
				Error: nil,
			})
		}()
	}()
        */
}



