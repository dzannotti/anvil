package eventbus

type Message struct {
	Data  any
	Kind  string
	Depth int
	End   bool
}
