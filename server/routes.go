package server

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/lovego/xiaomei/server/xm"
)

const alivePath = `/_alive`
const psPath = `/_ps`
const pprofPath = `/_pprof`

var sysPaths = map[string]bool{alivePath: true, psPath: true, pprofPath: true}

func sysRoutes(router *xm.Router) {
	router.Root().
		// 存活检测
		Get(alivePath, func(req *xm.Request, res *xm.Response) {
			res.Write([]byte(`ok`))
		}).
		// 当前正在处理的请求列表
		Get(psPath, func(req *xm.Request, res *xm.Response) {
			res.Write(psData.ToJson())
		}).
		// 性能分析
		Group(pprofPath).Get(`/`, routePprofIndex).GetX(`/(.+)`, routePprofGet)
}

var pprofIndexHtml []byte

func routePprofIndex(req *xm.Request, res *xm.Response) {
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

func routePprofGet(req *xm.Request, res *xm.Response, params []string) {
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
