package server

import (
	"net/http"
	"net/url"
)

type Request struct {
	*http.Request
	Session interface{}
}

func NewRequest(request *http.Request) *Request {
	return &Request{Request: request}
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
