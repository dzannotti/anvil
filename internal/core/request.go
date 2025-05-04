package core

type RequestOption struct {
	Value   any
	Label   string
	Default bool
}

type Request struct {
	Target   *Actor
	Text     string
	Options  []RequestOption
	Response chan RequestOption
	Resolved bool
	Chosen   RequestOption
}

func (r *Request) Answer(selected RequestOption) {
	if r.Resolved {
		r.Chosen = selected
		return
	}
	r.Resolved = true
	r.Response <- selected
	close(r.Response)
}

func (r *Request) DefaultOption() RequestOption {
	for _, o := range r.Options {
		if o.Default {
			return o
		}
	}
	return r.Options[0]
}

func (r *Request) AnswerWithDefault() {
	if r.Resolved {
		return
	}
	r.Answer(r.DefaultOption())
}
