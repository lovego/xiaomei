package xiaomei

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
	"reflect"

	"github.com/lovego/xiaomei/renderer"
	"github.com/lovego/xiaomei/session"
)

type LayoutDataFunc func(layout string, data interface{}, req *Request, res *Response) interface{}

type Response struct {
	http.ResponseWriter
	*Request
	sess           session.Session
	renderer       *renderer.Renderer
	layoutDataFunc LayoutDataFunc
	body           []byte
	err            error
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

func (res *Response) Write(content []byte) (int, error) {
	res.body = append(res.body, content...)
	return res.ResponseWriter.Write(content)
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

func (res *Response) GetBody() []byte {
	return res.body
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
