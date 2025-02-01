package onecontext

import (
	"context"
	"time"
)

// DetachedContext holds the logic to detach a cancellation signal from a context.
type DetachedContext struct {
	ctx    context.Context
	ch     chan struct{}
	cancel func()
}

// Detach detaches the cancellation signal from a context.
func Detach(ctx context.Context) (*DetachedContext, func()) {
	ch := make(chan struct{})
	cancel := func() {
		close(ch)
	}
	return &DetachedContext{
		ctx:    ctx,
		ch:     ch,
		cancel: cancel,
	}, cancel
}

// Deadline returns a nil deadline.
func (c *DetachedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// Done returns a cancellation signal that expires only when the context is canceled from the cancel function returned in Detach.
func (c *DetachedContext) Done() <-chan struct{} {
	return c.ch
}

// Err returns an error if the context is canceled from the cancel function returned in Detach.
func (c *DetachedContext) Err() error {
	select {
	case <-c.Done():
		return ErrCanceled
	default:
		return nil
	}
}

// Value returns the value associated with the key from the original context.
func (c *DetachedContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}
