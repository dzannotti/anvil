package core

type EventDispatcher interface {
	Begin(k string, e any)
	End()
	Emit(k string, e any)
}
