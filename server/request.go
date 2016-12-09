package server

import (
	"net/http"
	"net/url"

	"github.com/bughou-go/xiaomei/server/session"
)

type Request struct {
	*http.Request
	session requestSession
}
type requestSession struct {
	store  session.Store
	parsed bool
	data   interface{}
}

func NewRequest(request *http.Request, store session.Store) *Request {
	return &Request{Request: request, session: requestSession{store: store}}
}

func (req *Request) ClientAddr() string {
	addr := req.Header.Get("X-Real-IP")
	if addr != `` {
		return addr
	}
	return req.RemoteAddr
}

func (req *Request) Query() url.Values {
	return req.URL.Query()
}

func (req *Request) Session(p interface{}) {
	if req.session.parsed {
		p = req.session.data
		return
	}
	p = req.session.data
}
