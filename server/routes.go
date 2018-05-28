package server

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/router"
)

const pprofPath = `/_pprof`

var logBody = true

func sysRoutes(router *router.Router) {
	router.Root().
		// 存活检测
		Get(`/_alive`, func(req *xiaomei.Request, res *xiaomei.Response) {
			res.Write([]byte(`ok`))
		}).
		// 当前正在处理的请求列表
		Get(`/_ps`, func(req *xiaomei.Request, res *xiaomei.Response) {
			res.Write(psData.ToJson())
		}).
		// 返回当前机器是否记录的标记
		GetX(`/_log_body`,
			func(req *xiaomei.Request, res *xiaomei.Response, params []string) {
				hostname, err := os.Hostname()
				res.Data(struct {
					Host    string
					LogBody bool
				}{hostname + `:` + os.Getenv(`GOPORT`), logBody}, err)
			}).
		// 修改是否记录request body和response body的标记
		GetX(`/_log_body/(true|false)`,
			func(req *xiaomei.Request, res *xiaomei.Response, params []string) {
				logBody = params[1] == `true`
				res.Write([]byte(`ok`))
			}).
		// 性能分析
		Group(pprofPath).Get(`/`, routePprofIndex).GetX(`/(.+)`, routePprofGet)
}

var pprofIndexHtml []byte

func routePprofIndex(req *xiaomei.Request, res *xiaomei.Response) {
	if pprofIndexHtml == nil {
		var tmpl = template.Must(template.New(``).Parse(`<html>
<head>
<title>pprof/</title>
</head>
<body>
pprof<br>
<br>
profiles:<br>
<table>
{{range .}}
<tr><td align=right>{{.Count}}<td><a href="` + pprofPath + `/{{.Name}}?debug=1">{{.Name}}</a>
{{end}}
</table>
<br>
<a href="` + pprofPath + `/goroutine?debug=2">full goroutine stack dump</a><br>
</body>
</html>
`))
		profiles := pprof.Profiles()
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, profiles); err != nil {
			panic(err)
		}
		pprofIndexHtml = buf.Bytes()
	}
	res.Write(pprofIndexHtml)
}

func routePprofGet(req *xiaomei.Request, res *xiaomei.Response, params []string) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	p := pprof.Lookup(params[1])
	if p == nil {
		res.WriteHeader(404)
		fmt.Fprintf(res, "Unknown profile: %s\n", params[1])
		return
	}
	if params[1] == "heap" && req.FormValue("gc") != `` {
		runtime.GC()
	}
	debug, _ := strconv.Atoi(req.FormValue("debug"))
	p.WriteTo(res, debug)
}
