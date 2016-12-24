package bases

// Emitter can emits either an Item type or an error
type Emitter interface {
	Emit() (Item, error)
}

// Event is an alias for Emitter for consistency reasons
type Event interface {
	Emitter
}
