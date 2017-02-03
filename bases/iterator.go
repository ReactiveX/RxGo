package bases

// Iterator is a type with Next method
type Iterator interface {
	Next() (interface{}, error)
}
