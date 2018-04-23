package log

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei"
)

func getFields(
	req *xiaomei.Request, res *xiaomei.Response, t time.Time,
) map[string]interface{} {
	m := map[string]interface{}{
		`at`:       t.Format(time.RFC3339),
		`duration`: fmt.Sprintf(`%.2f`, float64(time.Since(t))/1e6),
		`host`:     req.Host,
		`method`:   req.Method,
		`path`:     req.URL.Path,
		`query`:    req.URL.RawQuery,
		`status`:   res.Status(),

		`req_body`: req.ContentLength,
		`res_body`: res.Size(),
		`proto`:    req.Proto,
		`ip`:       req.ClientAddr(),
		`agent`:    req.UserAgent(),
		`refer`:    req.Referer(),
	}

	var sess interface{}
	req.Session(&sess)
	if sess != nil {
		m[`session`] = sess
	}

	// custom log fields
	for k, v := range req.GetLog() {
		m[k] = v
	}
	return m
}

var fieldsAry = []string{
	`duration`, `host`, `method`, `path`, `query`, `status`,
	`req_body`, `res_body`, `proto`, `ip`, `agent`, `refer`,
}
var fieldsMap = makeFieldsMap(fieldsAry)

func formatFields(fields map[string]interface{}, highlight bool) (result string) {
	if highlight {
		result += color.GreenString("at: %v\n", fields[`at`])
	} else {
		result += fmt.Sprintf("at: %v\n", fields[`at`])
	}

	for _, k := range fieldsAry {
		result += fmt.Sprintf("%s: %v\n", k, fields[k])
	}
	for k, v := range fields {
		if !fieldsMap[k] {
			result += fmt.Sprintf("%s: %v\n", k, v)
		}
	}
	if err := fields[`err`]; err != nil {
		if errStr, ok := err.(string); ok {
			result += errStr + "\n"
		} else {
			result += fmt.Sprint(err) + "\n"
		}
	}
	if stack := fields[`stack`]; stack != nil {
		if stackStr, ok := stack.(string); ok {
			result += stackStr + "\n"
		} else {
			result += fmt.Sprint(stack) + "\n"
		}
	}
	return
}

func makeFieldsMap(ary []string) map[string]bool {
	m := map[string]bool{`at`: true}
	for _, k := range ary {
		m[k] = true
	}
	m[`err`] = true
	m[`stack`] = true
	return m
}
