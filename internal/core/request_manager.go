package core

import "errors"

type RequestManager struct {
	activeRequest *Request
}

func NewRequestManager() *RequestManager {
	return &RequestManager{}
}

func (rm *RequestManager) Ask(actor *Actor, text string, options []RequestOption) (RequestOption, error) {
	if rm.activeRequest != nil {
		return RequestOption{}, errors.New("there is already a pending request, please wait until it is resolved")
	}

	rm.activeRequest = &Request{
		Target:   actor,
		Text:     text,
		Options:  options,
		Response: make(chan RequestOption),
	}

	selectedOption := <-rm.activeRequest.Response
	rm.activeRequest = nil
	return selectedOption, nil
}

func (rm *RequestManager) HasPendingRequest() bool {
	return rm.activeRequest != nil
}

func (rm *RequestManager) GetPendingRequest() *Request {
	return rm.activeRequest
}

func (rm *RequestManager) AnswerDefault() error {
	if rm.activeRequest == nil {
		return errors.New("no pending request to answer")
	}

	rm.activeRequest.AnswerWithDefault()
	return nil
}
