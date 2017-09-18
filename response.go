package xiaomei

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"reflect"

	"github.com/lovego/xiaomei/renderer"
	"github.com/lovego/xiaomei/session"
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
		log.Panic(err)
	}
}

func (res *Response) Json2(data interface{}, err error) {
	if err != nil {
		res.LogError(err)
	}
	bytes, err := json.Marshal(data)
	if err == nil {
		res.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		res.Write(bytes)
	} else {
		panic(err)
	}
}

func (res Response) Data(data interface{}, err error) {
	type dataT struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	if err == nil {
		res.Json(dataT{Code: `ok`, Message: `success`, Data: data})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(dataT{Code: e.Code(), Message: e.Message(), Data: data})
		} else {
			res.LogError(err)
			res.Json(dataT{Code: `error`, Message: err.Error(), Data: data})
		}
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
			res.LogError(err)
			res.Json(result{Code: `error`, Message: err.Error(), Result: data})
		}
	}
}

func (res Response) Message(err error) {
	type result struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err == nil {
		res.Json(result{Code: `ok`, Message: `success`})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(result{Code: e.Code(), Message: e.Message()})
		} else {
			res.LogError(err)
			res.Json(result{Code: `error`, Message: err.Error()})
		}
	}
}

func (res *Response) LogError(err error) {
	if err == nil {
		return
	}
	log := map[string]interface{}{`err`: err}
	if stack, ok := err.(interface {
		Stack() string
	}); ok {
		log[`stack`] = stack.Stack()
	}
	res.Log(log)
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
