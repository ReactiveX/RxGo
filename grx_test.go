package grx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventHandlerImplementation(t *testing.T) {
	assert := assert.New(t)
	assert.Implements((*EventHandler)(nil), new(NextFunc), "NextFunc SHOULD implement EventHandler")
	assert.Implements((*EventHandler)(nil), new(ErrFunc), "ErrFunc SHOULD implement EventHandler")
	assert.Implements((*EventHandler)(nil), new(DoneFunc), "DoneFunc SHOULD implement EventHandler")
	assert.Implements((*EventHandler)(nil), new(BaseObserver), "BaseObserver SHOULD implement EventHandler")
}

func TestObservableIteratorImplementation(t *testing.T) {
	assert.Implements(t, (*Iterator)(nil), new(BaseObservable), "BaseObservable SHOULD implement Iterator")
}

func TestStreamImplementation(t *testing.T) {
	assert.Implements(t, (*Stream)(nil), new(BaseObservable), "BaseObservable SHOULD implement Stream")
}

func TestBaseObservableImplementation(t *testing.T) {
	assert.Implements(t, (*Observable)(nil), new(BaseObservable))
}

func TestSentinelImplementation(t *testing.T) {
	assert.Implements(t, (*Sentinel)(nil), new(BaseObserver), "BaseObserver SHOULD implement Sentinel")
}

func TestObserverImplementation(t *testing.T) {
	assert.Implements(t, (*Observer)(nil), new(BaseObserver), "BaseObserver SHOULD implement Observer")
}
