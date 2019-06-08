package rxgo

import (
	"testing"

	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
)

func TestSingleFilterNotMatching(t *testing.T) {
	got := 0

	Just(1, 2, 3).ElementAt(1).Filter(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			if i == 2 {
				return true
			}
		}
		return false
	}).Subscribe(handlers.NextFunc(func(i interface{}) {
		switch i := i.(type) {
		case Optional:
			if !i.IsEmpty() {
				g, _ := i.Get()
				got = g.(int)
			}
		}
	})).Block()

	assert.Equal(t, 2, got)
}

func TestSingleFilterMatching(t *testing.T) {
	got := 0

	Just(1, 2, 3).ElementAt(1).Filter(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			if i == 2 {
				return false
			}
		}
		return true
	}).Subscribe(handlers.NextFunc(func(i interface{}) {
		switch i := i.(type) {
		case Optional:
			if !i.IsEmpty() {
				g, _ := i.Get()
				got = g.(int)
			}
		}
	})).Block()

	assert.Equal(t, 0, got)
}

func TestSingleMap(t *testing.T) {
	got := 0

	Just(1, 2, 3).ElementAt(1).Map(func(i interface{}) interface{} {
		return i
	}).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(int) + 10
	})).Block()

	assert.Equal(t, 12, got)
}

func TestSingleMapWithTwoSubscription(t *testing.T) {
	just := newSingleFrom(1).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	}).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	})

	AssertSingle(t, just, HasValue(3))
	AssertSingle(t, just, HasValue(3))
}
