package core

type EventDispatcher interface {
	Begin(event any)
	End()
	Emit(event any)
}
