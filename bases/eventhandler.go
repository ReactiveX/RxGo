package bases

type EventHandler interface {
	Handle(Emitter)
}
