package rxgo

import (
	"testing"

	"errors"

	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewObserverWithConstructor(t *testing.T) {
	ob := NewObserver()
	ob.OnDone()
	ob.OnError(errors.New(""))
	ob.OnNext("")
	ob.OnNext(errors.New(""))
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

	ob := NewObserver(donef, nextf)

	ob.OnNext("Next")
	ob.OnDone()

	assert.Equal(t, "Next", nexttext)
	assert.Equal(t, "Hello", donetext)
}

func TestHandle(t *testing.T) {
	i := 0

	nextf := handlers.NextFunc(func(item interface{}) {
		i += 5
	})

	errorf := handlers.ErrFunc(func(error) {
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
