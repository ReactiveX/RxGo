package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/reactivex/rxgo/rx"

	"github.com/stretchr/testify/assert"
)

func TestImplementsEmitter(t *testing.T) {
	var subscription interface{} = New()
	_, isEmitter := subscription.(rx.Emitter)
	assert.True(t, isEmitter, "subscription should be implementation of emitter interface")
}

func TestCreateSubscription(t *testing.T) {
	sub := New()

	assert := assert.New(t)
	assert.Equal(time.Time{}, sub.SubscribeTime)
	assert.Equal(time.Time{}, sub.UnsubscribeTime)
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
	assert.WithinDuration(first, sub.SubscribeAt(), 5*time.Millisecond)
	assert.WithinDuration(first, sub.SubscribeAt(), 15*time.Millisecond)
}

func TestConstructEmitterFromContext(t *testing.T) {
	WithContext(context.Background())
}

func TestSubscriptionFromCancellableContext(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	emitter := WithContext(ctx)
	done := emitter.Done()
	cancelFunc()
	select {
	case <-done:
		assert.EqualError(t, emitter.Err(), context.Canceled.Error())
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
}

func TestSubscriptionFromTimeoutContext(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(0)) //nolint
	emitter := WithContext(ctx)
	done := emitter.Done()
	select {
	case <-done:
		assert.EqualError(t, emitter.Err(), context.DeadlineExceeded.Error())
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
	cancelFunc()
}

func TestSubscriptionSubscribeAt(t *testing.T) {
	ctx := context.Background()
	subscription := WithContext(ctx)
	beforeSubscription := time.Now()
	subscription = subscription.Subscribe()
	afterSubscription := time.Now()

	emitterSubscriptionTime := subscription.SubscribeAt()
	assert.True(t, emitterSubscriptionTime.After(beforeSubscription), "time less than a time evaluation before subscription")
	assert.True(t, emitterSubscriptionTime.Before(afterSubscription), "time greater than time evaluation after subscription")
}

func TestSubscriptionDeadline(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(0))
	emitter := WithContext(ctx)
	deadline, deadlineSet := emitter.Deadline()
	afterDeadlineTime := time.Now()
	done := emitter.Done()
	select {
	case <-done:
		assert.True(t, deadlineSet, "deadline should have been set")
		assert.True(t, deadline.Before(afterDeadlineTime), "deadline should have occurred already")
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
	cancelFunc()
}

func TestSubscriptionSettingError(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(0))
	emitter := WithContext(ctx)
	newFakeError := errors.New("some new error")
	emitter.Error = newFakeError
	done := emitter.Done()
	select {
	case <-done:
		assert.EqualError(t, emitter.Err(), newFakeError.Error())
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
	cancelFunc()
}

func TestSubscriptionValue(t *testing.T) {
	type test struct {
		ValuesToAssign map[interface{}]interface{}
		KeyToSearch    interface{}
		ExpectedResult interface{}
	}
	tests := []test{
		test{
			ValuesToAssign: map[interface{}]interface{}{"bob": 5},
			KeyToSearch:    "bob",
			ExpectedResult: 5},
		test{
			ValuesToAssign: map[interface{}]interface{}{},
			KeyToSearch:    "john",
			ExpectedResult: nil},
		test{
			ValuesToAssign: map[interface{}]interface{}{
				true:     true,
				false:    false,
				"string": false,
			},
			KeyToSearch:    true,
			ExpectedResult: true},
	}
	for _, test := range tests {
		ctx := context.Background()
		for k, v := range test.ValuesToAssign {
			ctx = context.WithValue(ctx, k, v)
		}
		subscription := WithContext(ctx)
		assert.Equal(t,
			test.ExpectedResult,
			subscription.Value(test.KeyToSearch),
			"value from search was not equivalent to expected")
	}

}
