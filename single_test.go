package rx

import (
	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
	"testing"
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
		got = i.(int)
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
		got = i.(int)
	})).Block()

	assert.Equal(t, 2, got)
}
