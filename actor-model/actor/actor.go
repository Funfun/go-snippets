package actor

type Event struct {
}

type Request struct {
}

type Response struct {
}

type Actor struct {
	eventCh   chan Event
	requestCh chan reqRes
	quitCh    chan struct{}
}

func (a *Actor) consumeEvent(e Event) {
	// do stuff with an event
}

func (a *Actor) handleRequest(r *Request) (res *Response, err error) {
	// do handle request in here

	return &Response{}, nil
}

func (a *Actor) Process() {
	for {
		select {
		case e := <-a.eventCh:
			a.consumeEvent(e)
		case r := <-a.requestCh:
			res, err := a.handleRequest(r.req)
			r.resc <- resErr{res, err}
		case <-a.quitCh:
			return
		}
	}
}

func (a *Actor) MakeRequest(r *Request) (*Response, error) {
	resc := make(chan resErr)
	a.requestCh <- reqRes{req: r, resc: resc}
	res := <-resc

	return res.res, res.err
}

func (a *Actor) SendEvent(e Event) {
	a.eventCh <- e
}

func (a *Actor) Stop() {
	close(a.quitCh)
}

type reqRes struct {
	req  *Request
	resc chan resErr
}

type resErr struct {
	res *Response
	err error
}
