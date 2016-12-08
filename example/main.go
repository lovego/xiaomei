package main

import (
	"os"
	"strings"

	"github.com/bughou-go/xiaomei/cli"
	"github.com/bughou-go/xiaomei/example/filter"
	"github.com/bughou-go/xiaomei/example/routes"
	"github.com/bughou-go/xiaomei/server"
)

func main() {
	if len(os.Args) > 1 {
		cli.Run()
	} else {
		svr := server.New(routes.Get())
		svr.FilterFunc = filter.Process
		svr.ListenAndServe()
	}
}

func layoutData(layout string, data interface{},
	req *server.Request, res *server.Response) interface{} {
	if strings.HasPrefix(layout, `layout/`) {
		return struct {
			Data interface{}
		}{data}
	}
	return data
}
