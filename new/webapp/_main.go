package main

import (
	_ "github.com/lovego/config/init" // this package must be the first.

	"runtime"

	"github.com/lovego/goa"
	"github.com/lovego/goa/server"
	"github.com/lovego/goa/utilroutes"
	"{{ .ProPath }}/middlewares"
	"{{ .ProPath }}/routes"
	// "{{ .ProPath }}/tasks"
)

func main() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}

	router := goa.New()
	router.Use(middlewares.Logger.Record)
	router.Use(middlewares.SessionParse)
	router.Use(middlewares.CORS.Check)
	utilroutes.Setup(router)
	router.Use(middlewares.Filter)

	routes.Setup(router)

	// tasks.Start()
	server.ListenAndServe(router)
}
