package subscription

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateSubscription(t *testing.T) {
	sub := New()

	assert := assert.New(t)
	assert.Equal(time.Time{}, sub.SubscribeAt)
	assert.Equal(time.Time{}, sub.UnsubscribeAt)
	assert.Nil(sub.Err())
}

func TestSubscription(t *testing.T) {
	var first time.Time
	go func() {
		first = time.Now()
	}()

	sub := New().Subscribe()
	<-time.After(10 * time.Millisecond)
	sub = sub.Unsubscribe()

	assert := assert.New(t)
	assert.WithinDuration(first, sub.SubscribeAt, 5*time.Millisecond)
	assert.WithinDuration(first, sub.SubscribeAt, 15*time.Millisecond)
}
