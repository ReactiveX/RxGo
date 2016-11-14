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
