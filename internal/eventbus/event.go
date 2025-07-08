package eventbus

type Event struct {
	Data  any
	Kind  string
	Depth int
	End   bool
}
