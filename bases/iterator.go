package bases

// Iterator is an interface with Next and HasNext methods
type Iterator interface {
	Next() (Emitter, error)
}
