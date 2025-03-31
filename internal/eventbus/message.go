package eventbus

type Message struct {
	Data  any
	Depth int
	IsEnd bool
}
