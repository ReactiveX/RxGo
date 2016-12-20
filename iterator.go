package grx

type Iterator interface {
	Next() (interface{}, error)
	HasNext() bool
}
