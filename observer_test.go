package grx

import (
        "testing"

        "github.com/stretchr/testify/assert"
)

func TestCreateNewObserverWithConstructor(t *testing.T) {
	ob := NewObserver()
	
	assert := assert.New(t)	
	assert.IsType((*Observer)(nil), ob)
	assert.NotNil(ob.observable.C)
	assert.NotNil(ob.observable.observer)
	assert.Nil(ob.NextHandler)
	assert.Nil(ob.ErrHandler)
	assert.Nil(ob.DoneHandler)
	assert.Equal(ob, ob.observable.observer)
}



