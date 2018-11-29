package optional

var emptyOptional = new(empty)

// Optional defines a container for empty values
type Optional interface {
	Get() interface{}
	IsEmpty() bool
}

type some struct {
	content interface{}
}

type empty struct {
}

func (s *some) Get() interface{} {
	return s.content
}

func (s *some) IsEmpty() bool {
	return false
}

func (e *empty) Get() interface{} {
	return nil
}

func (e *empty) IsEmpty() bool {
	return true
}

// Of returns a non empty optional
func Of(data interface{}) Optional {
	return &some{
		content: data,
	}
}

// Empty returns an empty optional
func Empty() Optional {
	return emptyOptional
}
