package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"reflect"

	"github.com/bughou-go/xiaomei/server/renderer"
)

type Response struct {
	http.ResponseWriter
	*Request
	*renderer.Renderer
	LayoutDataFunc
	rwValue reflect.Value
}

type LayoutDataFunc func(layout string, data interface{}, req *Request, res *Response) interface{}

func NewResponse(
	resWriter http.ResponseWriter, req *Request, render *renderer.Renderer, fun LayoutDataFunc,
) *Response {
	return &Response{
		ResponseWriter: resWriter,
		Request:        req,
		Renderer:       render,
		LayoutDataFunc: fun,
		rwValue:        reflect.Indirect(reflect.ValueOf(resWriter)),
	}
}

func (res *Response) Render(path string, data interface{}) {
	if res.LayoutDataFunc != nil && res.Layout != `` {
		data = res.LayoutDataFunc(res.Layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.Render(&buf, path, data)
	res.Write(buf.Bytes())
}

func (res *Response) RenderL(path, layout string, data interface{}) {
	if res.LayoutDataFunc != nil && layout != `` {
		data = res.LayoutDataFunc(layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.RenderL(&buf, path, layout, data)
	res.Write(buf.Bytes())
}

func (res *Response) RenderF(path string, funcs map[string]interface{}, data interface{}) {
	if res.LayoutDataFunc != nil && res.Layout != `` {
		data = res.LayoutDataFunc(res.Layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.RenderF(&buf, path, funcs, data)
	res.Write(buf.Bytes())
}

func (res *Response) RenderLF(path, layout string, funcs map[string]interface{}, data interface{}) {
	if res.LayoutDataFunc != nil && layout != `` {
		data = res.LayoutDataFunc(layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.RenderLF(&buf, path, layout, funcs, data)
	res.Write(buf.Bytes())
}

func (res *Response) RenderS(path string, data interface{}) string {
	if res.LayoutDataFunc != nil && res.Layout != `` {
		data = res.LayoutDataFunc(res.Layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.Render(&buf, path, data)
	return buf.String()
}

func (res *Response) RenderLS(path, layout string, data interface{}) string {
	if res.LayoutDataFunc != nil && layout != `` {
		data = res.LayoutDataFunc(layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.RenderL(&buf, path, layout, data)
	return buf.String()
}

func (res *Response) RenderFS(path string, funcs map[string]interface{}, data interface{}) string {
	if res.LayoutDataFunc != nil && res.Layout != `` {
		data = res.LayoutDataFunc(res.Layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.RenderF(&buf, path, funcs, data)
	return buf.String()
}

func (res *Response) RenderLFS(
	path, layout string, funcs map[string]interface{}, data interface{},
) string {
	if res.LayoutDataFunc != nil && layout != `` {
		data = res.LayoutDataFunc(layout, data, res.Request, res)
	}
	var buf bytes.Buffer
	res.Renderer.RenderLF(&buf, path, layout, funcs, data)
	return buf.String()
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

func (res Response) Json2(data interface{}) {
	bytes, err := json.Marshal(data)
	if err == nil {
		res.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		res.Write(bytes)
	} else {
		res.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		res.Write([]byte(`null`))
	}
}

func (res *Response) Redirect(path string) {
	res.Header().Set("Location", path)
	res.WriteHeader(302)
}

func (res *Response) Status() int64 {
	return res.rwValue.FieldByName(`status`).Int()
}

func (res *Response) Size() int64 {
	return res.rwValue.FieldByName(`written`).Int()
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
