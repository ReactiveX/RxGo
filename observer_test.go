package rxgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewObserverWithConstructor(t *testing.T) {
	ob := NewObserver()
	err := ob.OnDone()
	assert.NoError(t, err)

	err = ob.OnNext("a")
	assert.Error(t, err)
	assert.IsType(t, &ClosedObserverError{}, err)
}

func TestCreateNewObserverWithObserver(t *testing.T) {
	nexttext := ""
	donetext := ""

	nextf := NextFunc(func(item interface{}) {
		if text, ok := item.(string); ok {
			nexttext = text
		}
	})

	donef := DoneFunc(func() {
		donetext = "Hello"
	})

	ob := NewObserver(donef, nextf)

	err := ob.OnNext("Next")
	assert.NoError(t, err)
	err = ob.OnDone()
	assert.NoError(t, err)

	assert.Equal(t, "Next", nexttext)
	assert.Equal(t, "Hello", donetext)
}

func TestHandle(t *testing.T) {
	i := 0

	nextf := NextFunc(func(item interface{}) {
		i += 5
	})

	errorf := ErrFunc(func(error) {
		i += 2
	})

	ob := NewObserver(nextf, errorf)
	ob.Handle("")
	ob.Handle(errors.New(""))
	assert.Equal(t, 7, i)
}

func BenchmarkObserver_IsDisposed(b *testing.B) {
	for n := 0; n < b.N; n++ {
		o := NewObserver()
		for i := 0; i < 10; i++ {
			o.IsDisposed()
		}
		o.Dispose()
		o.IsDisposed()
	}
}
