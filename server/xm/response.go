package xm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"reflect"

	"github.com/lovego/xiaomei/server/xm/renderer"
	"github.com/lovego/xiaomei/server/xm/session"
	"github.com/lovego/xiaomei/utils/errs"
)

type LayoutDataFunc func(layout string, data interface{}, req *Request, res *Response) interface{}

type Response struct {
	http.ResponseWriter
	*Request
	sess           session.Session
	renderer       *renderer.Renderer
	layoutDataFunc LayoutDataFunc
}

func NewResponse(
	responseWriter http.ResponseWriter, request *Request, sess session.Session,
	rendrr *renderer.Renderer, layoutDataFunc LayoutDataFunc,
) *Response {
	return &Response{
		ResponseWriter: responseWriter,
		Request:        request,
		sess:           sess,
		renderer:       rendrr,
		layoutDataFunc: layoutDataFunc,
	}
}

func (res *Response) Session(data interface{}) {
	res.sess.Set(res.ResponseWriter, data)
}

func (res *Response) GetLayoutData(layout string, data interface{}) interface{} {
	if res.layoutDataFunc == nil {
		return data
	}
	return res.layoutDataFunc(layout, data, res.Request, res)
}

func (res *Response) Render(path string, data interface{}, options ...renderer.O) {
	var option renderer.O
	if len(options) > 0 {
		option = options[0]
	}
	option.LayoutDataGetter = res
	var buf bytes.Buffer
	res.renderer.Render(&buf, path, data, option)
	res.Write(buf.Bytes())
}

func (res *Response) RenderToBuffer(
	path string, data interface{}, options ...renderer.O,
) bytes.Buffer {
	var option renderer.O
	if len(options) > 0 {
		option = options[0]
	}
	option.LayoutDataGetter = res
	var buf bytes.Buffer
	res.renderer.Render(&buf, path, data, option)
	return buf
}

func (res *Response) Json(data interface{}) {
	bytes, err := json.Marshal(data)
	if err == nil {
		res.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		res.Write(bytes)
	} else {
		panic(err)
	}
}

func (res Response) Result(data interface{}, err error) {
	type result struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Result  interface{} `json:"result,omitempty"`
	}
	if err == nil {
		res.Json(result{Code: `ok`, Message: `success`, Result: data})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(result{Code: e.Code(), Message: e.Message(), Result: data})
		} else {
			res.Json(result{Code: `error`, Message: e.Error(), Result: data})
		}
	}
}

func (res *Response) Redirect(path string) {
	res.Header().Set("Location", path)
	res.WriteHeader(302)
}

func (res *Response) Status() int64 {
	s := reflect.ValueOf(res.ResponseWriter).Elem().FieldByName(`status`)
	if s.IsValid() {
		return s.Int()
	} else {
		return 0
	}
}

func (res *Response) Size() int64 {
	s := reflect.ValueOf(res.ResponseWriter).Elem().FieldByName(`written`)
	if s.IsValid() {
		return s.Int()
	} else {
		return 0
	}
}

func (res *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := res.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
}

func (res *Response) Flush() {
	if flusher, ok := res.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
