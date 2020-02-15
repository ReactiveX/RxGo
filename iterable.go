package rxgo

type Iterable interface {
	Next() <-chan Item
}
