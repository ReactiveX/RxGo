package bases

type EventHandler interface {
	Handle(interface{})
}
