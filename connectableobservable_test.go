package rxgo

import (
	"testing"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/options"
	"github.com/stretchr/testify/assert"
)

func TestConnectableObservable(t *testing.T) {
	in := make(chan interface{}, 2)
	out1 := make(chan interface{}, 2)
	out2 := make(chan interface{}, 2)
	obs := FromChannel(in).Publish()
	obs.Subscribe(handlers.NextFunc(func(i interface{}) {
		out1 <- i
	}), options.WithBufferBackpressureStrategy(2))
	obs.Subscribe(handlers.NextFunc(func(i interface{}) {
		out2 <- i
	}), options.WithBufferBackpressureStrategy(2))
	in <- 1
	in <- 2
	_, _, cancelled := channel(out1, wait)
	assert.True(t, cancelled)
	obs.Connect()
	item, _, _ := channel(out1, wait)
	assert.Equal(t, 1, item)
	item, _, _ = channel(out1, wait)
	assert.Equal(t, 2, item)
	item, _, _ = channel(out2, wait)
	assert.Equal(t, 1, item)
	item, _, _ = channel(out2, wait)
	assert.Equal(t, 2, item)
}

func TestConnectableObservable_Map(t *testing.T) {
	obs := FromSlice([]interface{}{1, 2, 3, 5}).Publish().Map(func(i interface{}) interface{} {
		return i.(int) + 1
	})
	AssertObservable(t, obs, HasItems(2, 3, 4, 6))
}
