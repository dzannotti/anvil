package core

type LogWriter interface {
	Start(k string, e any)
	End()
	Add(k string, e any)
}
