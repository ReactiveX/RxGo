package rx

import (
	"testing"
	"github.com/reactivex/rxgo/observable"
	"github.com/stretchr/testify/assert"
)

func TestObservableDoesImplementRxObservable(t *testing.T) {
	source := observable.New(1)
	assert.IsType(t, observable.Observable(nil), source)
	// 	assert.Implements(t, (Observable)(nil), source)
	// }
}
