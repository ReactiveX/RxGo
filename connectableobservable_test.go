package rxgo

import (
	"errors"
	"testing"

	"github.com/reactivex/rxgo/v2/handlers"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	just := Just(1, 2, 3).Publish()
	got1 := make([]interface{}, 0)
	got2 := make([]interface{}, 0)

	just.Subscribe(handlers.NextFunc(func(i interface{}) {
		got1 = append(got1, i)
	}))

	just.Subscribe(handlers.NextFunc(func(i interface{}) {
		got2 = append(got2, i)
	}))

	just.Connect().Block()
	assert.Equal(t, []interface{}{1, 2, 3}, got1)
	assert.Equal(t, []interface{}{1, 2, 3}, got2)
}

func TestConnectOnError(t *testing.T) {
	just := Just(1, 2, 3, errors.New("foo"), 4).Publish()
	got1 := make([]interface{}, 0)
	got2 := make([]interface{}, 0)

	just.Subscribe(handlers.NextFunc(func(i interface{}) {
		got1 = append(got1, i)
	}))

	just.Subscribe(handlers.NextFunc(func(i interface{}) {
		got2 = append(got2, i)
	}))

	just.Connect().Block()
	assert.Equal(t, []interface{}{1, 2, 3}, got1)
	assert.Equal(t, []interface{}{1, 2, 3}, got2)
}
