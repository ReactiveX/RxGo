package bases

// Item is a higher-level alias for empty interface type
type Item interface{}

// Emitter can emits either an Item type or an error
type Emitter interface {
	Emit() (Item, error)
}

// Event is an alias for Emitter for consistency
type Event interface {
	Emitter
}
