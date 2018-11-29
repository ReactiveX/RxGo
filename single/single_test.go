package single

import (
	"fmt"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
	"testing"
)

func TestSingleFilter(t *testing.T) {
	observable.Just(1, 2, 3).ElementAt(1).Filter(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			if i == 1 {
				return false
			}
		}
		return true
	}).Subscribe(handlers.NextFunc(func(i interface{}) {
		fmt.Printf("%v\n", i)
	})).Block()
}
