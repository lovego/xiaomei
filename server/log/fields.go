package log

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei"
)

type logFields struct {
	*tracer.Span

	Host   string `json:"host"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Query  string `json:"query"`
	Status int64  `json:"status"`

	ReqBody     string `json:"req_body"`
	ReqBodySize int64  `json:"req_body_size"`
	ResBody     string `json:"res_body"`
	ResBodySize int64  `json:"res_body_size"`
	Proto       string `json:"proto"`
	Ip          string `json:"ip"`
	Agent       string `json:"agent"`
	Refer       string `json:"refer"`

	Session interface{} `json:"session,omitempty"`
	Error   string      `json:"error,omitempty"`
	Stack   string      `json:"stack,omitempty"`
}

func getFields(req *xiaomei.Request, res *xiaomei.Response, logBody bool) *logFields {
	var sess interface{}
	req.Session(&sess)

	fields := &logFields{
		Span: req.Span,

		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Query:  req.URL.RawQuery,
		Status: res.Status(),

		ReqBodySize: req.ContentLength,
		ResBodySize: res.Size(),
		Proto:       req.Proto,
		Ip:          req.ClientAddr(),
		Agent:       req.UserAgent(),
		Refer:       req.Referer(),
		Session:     sess,

		Error: req.Error,
		Stack: req.Stack,
	}
	if logBody {
		fields.ReqBody = string(req.GetBody())
		fields.ResBody = string(res.GetBody())
	}
	return fields
}

func formatFields(f *logFields, highlight, logBody bool) (result string) {
	if highlight {
		result += color.GreenString("at: %v", f.At)
	} else {
		result += fmt.Sprintf("at: %v", f.At)
	}

	line := fmt.Sprintf(`
duration: %v
host: %s
method: %s
path: %s
query: %s
status: %d
req_body_size: %d
res_body_size: %d
proto: %s
ip: %s
agent: %s
refer: %s
session: %+v
children: %+v
tags: %v
`, f.Duration, f.Host, f.Method, f.Path, f.Query, f.Status,
		f.ReqBodySize, f.ResBodySize, f.Proto, f.Ip, f.Agent, f.Refer, f.Session,
		f.Children, f.Tags,
	)
	if logBody {
		line += fmt.Sprintf("req_body: %s\nres_body: %s\n", f.ReqBody, f.ResBody)
	}
	result += line

	if f.Error != "" {
		result += f.Error + "\n"
	}
	if f.Stack != "" {
		result += f.Stack + "\n"
	}
	return
}
