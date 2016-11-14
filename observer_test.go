package grx

import (
        "testing"

        "github.com/stretchr/testify/assert"
)

func TestObserverImplementSentinel(t *testing.T) {
	assert.Implements(t, (*Sentinel)(nil), new(Observer))
}
