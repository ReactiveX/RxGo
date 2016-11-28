package grx

import (
        "testing"

        "github.com/stretchr/testify/assert"
)

func TestHandlersImplementHandlerFunc(t *testing.T) {
        assert.Implements(t, (*Handler)(nil), new(NextFunc))
        assert.Implements(t, (*Handler)(nil), new(ErrFunc))
        assert.Implements(t, (*Handler)(nil), new(DoneFunc))
}

func TestObservableImplementStream(t *testing.T) {
	assert.Implements(t, (*Stream)(nil), new(Observable))
}

func TestObserverImplementSentinel(t *testing.T) {
	assert.Implements(t, (*Sentinel)(nil), new(Observer))
}



