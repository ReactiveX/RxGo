package observer

import (
	"testing"

	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/handlers"
	"github.com/stretchr/testify/assert"
)

func TestObserverImplementsBaseObserver(t *testing.T) {
	assert.Implements(t, (*bases.Observer)(nil), Observer{})
}

func TestCreateNewObserverWithConstructor(t *testing.T) {
	assert := assert.New(t)
	ob := New()

	assert.IsType(Observer{}, ob)
	assert.NotNil(ob.NextHandler)
	assert.NotNil(ob.ErrHandler)
	assert.NotNil(ob.DoneHandler)
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
