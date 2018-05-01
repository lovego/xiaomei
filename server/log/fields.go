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

	ReqBody int64  `json:"req_body"`
	ResBody int64  `json:"res_body"`
	Proto   string `json:"proto"`
	Ip      string `json:"ip"`
	Agent   string `json:"agent"`
	Refer   string `json:"refer"`

	Session interface{} `json:"session,omitempty"`
	Error   string      `json:"error,omitempty"`
	Stack   string      `json:"stack,omitempty"`
}

func getFields(req *xiaomei.Request, res *xiaomei.Response) *logFields {
	var sess interface{}
	req.Session(&sess)

	return &logFields{
		Span: req.Span,

		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Query:  req.URL.RawQuery,
		Status: res.Status(),

		ReqBody: req.ContentLength,
		ResBody: res.Size(),
		Proto:   req.Proto,
		Ip:      req.ClientAddr(),
		Agent:   req.UserAgent(),
		Refer:   req.Referer(),
		Session: sess,

		Error: req.Error,
		Stack: req.Stack,
	}
}

func formatFields(f *logFields, highlight bool) (result string) {
	if highlight {
		result += color.GreenString("at: %v", f.At)
	} else {
		result += fmt.Sprintf("at: %v", f.At)
	}

	result += fmt.Sprintf(`
duration: %v
host: %s
method: %s
path: %s
query: %s
status: %d
req_body: %d
res_body: %d
proto: %d
ip: %s
agent: %s
refer: %s
session: %v
children: %+v
tags: %v
`, f.Duration, f.Host, f.Method, f.Path, f.Query, f.Status,
		f.ReqBody, f.ResBody, f.Proto, f.Ip, f.Agent, f.Refer, f.Session,
		f.Children, f.Tags,
	)

	if f.Error != "" {
		result += f.Error + "\n"
	}
	if f.Stack != "" {
		result += f.Stack + "\n"
	}
	return
}
