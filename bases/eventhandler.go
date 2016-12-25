package bases

type EventHandler interface {
	Apply(Emitter)
}
