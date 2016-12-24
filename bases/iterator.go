package bases

// Iterator is an interface with Next and HasNext methods
type Iterator interface {
	Next() (Item, error)
	HasNext() bool
}
