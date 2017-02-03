package observer

import (
	"testing"

	"github.com/jochasinga/rx/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewObserverWithConstructor(t *testing.T) {
	assert := assert.New(t)
	ob := New()

	assert.IsType(Observer{}, ob)
	assert.NotNil(ob.NextHandler)
	assert.NotNil(ob.ErrHandler)
	assert.NotNil(ob.DoneHandler)

	onNext := handlers.NextFunc(func(item interface{}) {})
	onError := handlers.ErrFunc(func(err error) {})

	ob2 := New(onNext, onError)
	assert.NotNil(ob2.NextHandler)
	assert.NotNil(ob2.ErrHandler)
	assert.NotNil(ob2.DoneHandler)

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
