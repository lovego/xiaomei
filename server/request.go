package server

import (
	"net/http"
	"net/url"

	"github.com/bughou-go/xiaomei/server/session"
)

type Request struct {
	*http.Request
	sess       session.Session
	sessParsed bool
	sessData   interface{}
}

func NewRequest(request *http.Request, sess session.Session) *Request {
	return &Request{Request: request, sess: sess}
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

func (req *Request) Session(p *interface{}) {
	if req.sessParsed {
		*p = req.sessData
	}
	req.sess.Get(req.Request, p)
	req.sessData = *p
	req.sessParsed = true
}
