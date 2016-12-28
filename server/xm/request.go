package xm

import (
	"net/http"
	"net/url"
	"reflect"

	"github.com/bughou-go/xiaomei/server/xm/session"
)

type Request struct {
	*http.Request
	sess       session.Session
	sessParsed bool
	sessData   reflect.Value
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

// req.Session retains the value the param pointer points to.
// the next call to req.Session will return the value instead of parsing from req again.
// so modification to the session value will remain.
func (req *Request) Session(p interface{}) {
	if req.sessParsed {
		if req.sessData.IsValid() {
			reflect.ValueOf(p).Elem().Set(req.sessData)
		}
		return
	}
	req.sess.Get(req.Request, p)
	if p != nil {
		v := reflect.ValueOf(p).Elem()
		if v.IsValid() {
			req.sessData = v
		}
	}
	req.sessParsed = true
}
