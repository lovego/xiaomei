package main

import (
	_ "github.com/bughou-go/xiaomei/server/a-init" // this package must be the first.

	"strings"

	"github.com/bughou-go/xiaomei/example/filter"
	"github.com/bughou-go/xiaomei/example/routes"
	"github.com/bughou-go/xiaomei/server"
	"github.com/bughou-go/xiaomei/server/xm"
)

func main() {
	svr := &server.Server{
		FilterFunc:     filter.Process,
		Router:         routes.Get(),
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
