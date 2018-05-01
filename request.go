package xiaomei

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei/session"
)

type Request struct {
	*http.Request
	sess         session.Session
	sessParsed   bool
	sessData     reflect.Value
	Span         *tracer.Span
	Error, Stack string
}

func NewRequest(request *http.Request, sess session.Session, at time.Time) *Request {
	return &Request{Request: request, sess: sess, Span: &tracer.Span{At: at}}
}

func (req *Request) ClientAddr() string {
	if addrs := req.Header.Get("X-Forwarded-For"); addrs != `` {
		addr := strings.SplitN(addrs, `, `, 2)[0]
		if addr != `unknown` {
			return addr
		}
	}
	if addr := req.Header.Get("X-Real-IP"); addr != `` && addr != `unknown` {
		return addr
	}
	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	return host
}

func (req *Request) Scheme() string {
	if proto := req.Header.Get("X-Forwarded-Proto"); proto != `` {
		return proto
	}
	return `http`
}

func (req *Request) Url() string {
	return req.Scheme() + `://` + req.Host + req.RequestURI
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
	if req.sess != nil && req.sess.Get(req.Request, p) && p != nil {
		v := reflect.ValueOf(p).Elem()
		if v.IsValid() {
			req.sessData = v
		}
	}
	req.sessParsed = true
}

func (req *Request) SetSession(i interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.IsValid() {
		req.sessData = v
		req.sessParsed = true
	}
}

func (req *Request) Context() context.Context {
	return tracer.Context(context.Background(), req.Span)
}

func (req *Request) GetBody() []byte {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("read http body error: %v", err)
	}
	return buf
}
