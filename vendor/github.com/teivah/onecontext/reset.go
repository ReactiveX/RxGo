package onecontext

import (
	"context"
	"time"
)

// ResetValuesContext holds the logic reset the values of a context.
type ResetValuesContext struct {
	ctx context.Context
}

// ResetValues reset the values of a context.
func ResetValues(ctx context.Context) *ResetValuesContext {
	return &ResetValuesContext{
		ctx: ctx,
	}
}

// Deadline returns the original context deadline.
func (c *ResetValuesContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

// Done returns the original done channel.
func (c *ResetValuesContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err returns the original context error.
func (c *ResetValuesContext) Err() error {
	return c.ctx.Err()
}

// Value returns nil regardless of the key.
func (c *ResetValuesContext) Value(_ interface{}) interface{} {
	return nil
}
