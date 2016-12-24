package bases

type Sentinel interface {
	OnNext(Item)
	OnError(error)
	OnDone()
}
