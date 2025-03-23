package log

type Event struct {
	Data  any
	Depth int
	IsEnd bool
}

func NewEvent(data any) Event {
	return Event{Data: data, Depth: 0, IsEnd: false}
}
