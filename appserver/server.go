package main

import (
	"fmt"
	"github.com/bughou-go/xiaomei/appserver/filter"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/routes"
	"github.com/bughou-go/xm"
	"net/http"
	"strings"
	"time"
)

func main() {
	var router = routes.Get()
	var renderer = getRenderer()

	addr := config.CurrentAppServer().AppAddr + `:` + config.Data.AppPort
	fmt.Printf("%s listen at %s\n", time.Now().Format(`2006-01-02 15:04:05 -0700`), addr)

	if err := http.ListenAndServe(addr, http.HandlerFunc(
		func(response_writer http.ResponseWriter, request *http.Request) {
			req := xm.NewRequest(request)
			res := xm.NewResponse(response_writer, req, renderer, layoutData)

			defer errorHandler(time.Now(), req, res)

			// 如果返回true，继续交给路由处理
			if filter.Process(req, res) {
				router.Handle(req, res, notFound)
			}
		})); err != nil {
		panic(err)
	}
}

func layoutData(layout string, data interface{}, req *xm.Request, res *xm.Response) interface{} {
	if strings.HasPrefix(layout, `layout/`) {
		return struct {
			Data interface{}
		}{data}
	}
	return data
}
