package grx

import (
        "testing"

        "github.com/stretchr/testify/assert"
)

func TestCreateNewObserverWithConstructor(t *testing.T) {
	//ob := NewObserver()
	ob := NewBaseObserver()
	
	assert := assert.New(t)	
	assert.IsType((*Observer)(nil), ob)
	assert.NotNil(ob._observable.getC())
	assert.NotNil(ob._observable.getInnerObserver())
	assert.Nil(ob.NextHandler)
	assert.Nil(ob.ErrHandler)
	assert.Nil(ob.DoneHandler)
	assert.Equal(ob, ob._observable.getInnerObserver())
}



