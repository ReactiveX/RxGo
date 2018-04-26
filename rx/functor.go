package rx

type (
	MappableFunc func(interface{}) interface{}
	FilterableFunc func(interface{}) bool
	IterableFunc func(interface{})
	ScannableFunc func(interface{}, interface{}) interface{}
)

type Functor interface {
	Map(MappableFunc) Observable
	Filter(FilterableFunc) Observable
	ForEach(IterableFunc)
	Scan(ScannableFunc) Observable
}
	

