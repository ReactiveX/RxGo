package rxgo

import (
	"testing"

	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewSingleObserverWithConstructor(t *testing.T) {
	var got int64
	single := Just(1, 2, 3).Count().Subscribe(
		handlers.NextFunc(func(i interface{}) {
			got = i.(int64)
		}))

	v, err := single.Block()

	assert.Nil(t, err)
	assert.Equal(t, int64(3), got)
	assert.Equal(t, int64(3), v)
	assert.True(t, single.IsDisposed())
}

func TestCreateNewSingleObserverFromSingleObserver(t *testing.T) {
	var got int64
	singleObserver := CheckSingleEventHandler(handlers.NextFunc(func(i interface{}) {
		got = i.(int64)
	}))

	single := Just(1, 2, 3).Count().Subscribe(singleObserver)
	v, err := single.Block()

	assert.Nil(t, err)
	assert.Equal(t, int64(3), got)
	assert.Equal(t, int64(3), v)
	assert.True(t, single.IsDisposed())
}

func BenchmarkIsDisposed(b *testing.B) {
	for n := 0; n < b.N; n++ {
		o := NewSingleObserver()
		for i := 0; i < 10; i++ {
			o.IsDisposed()
		}
		o.Dispose()
		o.IsDisposed()
	}
}
