package tooth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStream(t *testing.T) {
	testStream := NewObservable("myStream")
	assert.IsType(t, &Observable{}, testStream)
}

func TestStreamSubscription(t *testing.T) {
	testStream := NewObservable("myStream")
	err := testStream.Subscribe(Observer(func(e Event, c Context) interface{} {
		// empty handler
		return nil
	}))
	assert.Error(t, err)
}

