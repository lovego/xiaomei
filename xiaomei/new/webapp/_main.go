package main

import (
	_ "github.com/lovego/xiaomei/server/init" // this package must be the first.

	"strings"

	"{{ .ProPath }}/filter"
	"{{ .ProPath }}/routes"
	"github.com/lovego/xiaomei/server"
	"github.com/lovego/xiaomei/server/xm"
)

func main() {
	svr := &server.Server{
		FilterFunc:     filter.Process,
		Router:         routes.Routes(),
		Session:        server.NewSession(),
		Renderer:       server.NewRenderer(),
		LayoutDataFunc: layoutData,
	}
	svr.ListenAndServe()
}

func layoutData(
	layout string, data interface{}, req *xm.Request, res *xm.Response,
) interface{} {
	if strings.HasPrefix(layout, `layout/`) {
		return struct {
			Data interface{}
		}{data}
	}
	return data
}
