package core

type EventDispatcher interface {
	Start(k string, e any)
	End()
	Add(k string, e any)
}
