package main

import (
	_ "github.com/lovego/xiaomei/server/init" // this package must be the first.

	"strings"

	"{{ .ProPath }}/filter"
	"{{ .ProPath }}/routes"
	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/server"
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
	layout string, data interface{}, req *xiaomei.Request, res *xiaomei.Response,
) interface{} {
	if strings.HasPrefix(layout, `layout/`) {
		return struct {
			Data interface{}
		}{data}
	}
	return data
}
